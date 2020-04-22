// Package parser implement a recursive descent parser.
package parser

import (
	"fmt"
	"strconv"

	"github.com/actuallyachraf/monkey-giggle/ast"
	"github.com/actuallyachraf/monkey-giggle/lexer"
	"github.com/actuallyachraf/monkey-giggle/token"
)

const (
	// Enumerate operations by their precedence order.
	_ int = iota
	// LOWEST marks lowest precedence order
	LOWEST
	// EQUALS marks equality
	EQUALS
	// LESSGREATER marks lesser or greater than operations
	LESSGREATER
	// SUM marks sum operation
	SUM
	// PRODUCT marks product operation
	PRODUCT
	// PREFIX marks prefix operators
	PREFIX
	// CALL marks function calls
	CALL
)

var precedenceTable = map[token.Type]int{
	token.EQ:     EQUALS,
	token.NEQ:    EQUALS,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.GE:     LESSGREATER,
	token.LE:     LESSGREATER,
	token.ADD:    SUM,
	token.SUB:    SUM,
	token.MUL:    PRODUCT,
	token.DIV:    PRODUCT,
	token.MOD:    PRODUCT,
	token.LPAREN: CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser implements the main parsing structure.
type Parser struct {
	l         *lexer.Lexer
	currToken token.Token
	peekToken token.Token

	errors []string

	prefixParseFuncs map[token.Type]prefixParseFn
	infixParseFuncs  map[token.Type]infixParseFn
}

// New creates a new parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:                l,
		errors:           []string{},
		prefixParseFuncs: make(map[token.Type]prefixParseFn),
		infixParseFuncs:  make(map[token.Type]infixParseFn),
	}

	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.SUB, p.parsePrefixExpression)

	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	p.registerInfix(token.ADD, p.parseInfixExpression)
	p.registerInfix(token.SUB, p.parseInfixExpression)
	p.registerInfix(token.MUL, p.parseInfixExpression)
	p.registerInfix(token.DIV, p.parseInfixExpression)
	p.registerInfix(token.MOD, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NEQ, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GE, p.parseInfixExpression)
	p.registerInfix(token.LE, p.parseInfixExpression)

	p.nextToken()
	p.nextToken()
	return p
}

// registerPrefix parsing function for a token type
func (p *Parser) registerPrefix(t token.Type, fn prefixParseFn) {
	p.prefixParseFuncs[t] = fn
}

// registerInfix parsing function for a token type
func (p *Parser) registerInfix(t token.Type, fn infixParseFn) {
	p.infixParseFuncs[t] = fn
}

// Errors returns the list of errors that occured during parsing
func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t token.Type) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// nextToken reads the next token and updates the fields.
func (p *Parser) nextToken() {
	p.currToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// parseIdentifier expression.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
}

// Parse is the main function call given a lexer instance it will parse
// and construct an abstract syntax tree for the given input.
func (p *Parser) Parse() *ast.Program {

	program := &ast.Program{Statements: []ast.Statement{}}

	for p.currToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// parseStatement parse a statement given it's defining token.
func (p *Parser) parseStatement() ast.Statement {

	switch p.currToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseIntegerLiteral parse a literal int from string to int64
func (p *Parser) parseIntegerLiteral() ast.Expression {

	value, err := strconv.ParseInt(string(p.currToken.Literal), 0, 64)
	if err != nil {
		msg := fmt.Sprintf("failed to parse %q as int64", p.currToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	return &ast.IntegerLiteral{
		Token: p.currToken,
		Value: value,
	}
}

// parsePrefixExpression constructs an AST expression node for prefix expression.
func (p *Parser) parsePrefixExpression() ast.Expression {

	expression := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}
	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression constructs an AST expression node for infix expressions.
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {

	expression := &ast.InfixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal,
		Left:     left,
	}
	precedence := p.currPrecendence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// noPrefixParseFnError writes an error message when a prefix parsing func
// isn't found for a given token type.
func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse func found for token type %s", t)
	p.errors = append(p.errors, msg)
}

// peekPrecedence returns the precedence level of the peek token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedenceTable[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

// currPrecendence returns the precedence level of the current token
func (p *Parser) currPrecendence() int {
	if p, ok := precedenceTable[p.currToken.Type]; ok {
		return p
	}

	return LOWEST
}

// noInfixParseFnError writes an error message when an infix parsing func
// isn't found for a given token type.
func (p *Parser) noInfixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no infix parse func found for token type %s", t)
	p.errors = append(p.errors, msg)
}

// parseExpression parses an expression given a precedence enum
func (p *Parser) parseExpression(precedence int) ast.Expression {

	prefix := p.prefixParseFuncs[p.currToken.Type]

	if prefix == nil {
		p.noPrefixParseFnError(p.currToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFuncs[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

// parseLetStatement parses and construct an ast node for the let statement.
func (p *Parser) parseLetStatement() *ast.LetStatement {

	stmt := &ast.LetStatement{
		Token: p.currToken,
	}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses and construct an ast node for return statements.
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {

	stmt := &ast.ReturnStatement{Token: p.currToken}
	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)

	for !p.currTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpressionStatement parses and construct an ast node for expressions.
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.currToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseBooleanLiteral parses and construct an ast node for boolean literals.
func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{
		Token: p.currToken,
		Value: p.currTokenIs(token.TRUE),
	}
}

// parseGroupedExpression parses and construct an ast node for expressions
// of the type ((a*b)+c).
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseIfExpression parses and construct an ast branch for conditionals
func (p *Parser) parseIfExpression() ast.Expression {

	exp := &ast.IfExpression{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	exp.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	exp.Consequence = p.parseBlockStatement()

	if p.peekTokenIs(token.ELSE) {
		p.nextToken()

		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		exp.Alternative = p.parseBlockStatement()
	}
	return exp
}

// parseBlockStatement parses and construct ast nodes for the block statements
// that follow a conditional branch.
func (p *Parser) parseBlockStatement() *ast.BlockStatement {

	block := &ast.BlockStatement{Token: p.currToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()

		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}

		p.nextToken()
	}

	return block
}

// parseFunctionParameters is used to construct a list of identifiers for
// function literal parameters
func (p *Parser) parseFunctionParameters() []*ast.Identifier {

	identifiers := []*ast.Identifier{}

	// if the next token is the right parenthesis return (they are no params)
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		ident := &ast.Identifier{Token: p.currToken, Value: p.currToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseFunctionLiteral constructs an ast branch for function literal expressions.
func (p *Parser) parseFunctionLiteral() ast.Expression {

	lit := &ast.FunctionLiteral{Token: p.currToken}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

// parseCallArguments is used to parse function arguments which are expressions.
func (p *Parser) parseCallArguments() []ast.Expression {

	args := []ast.Expression{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()

	args = append(args, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()

		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// parseCallExpression parses function calls.
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {

	return &ast.CallExpression{
		Token:     p.currToken,
		Function:  function,
		Arguments: p.parseCallArguments(),
	}
}
func (p *Parser) currTokenIs(t token.Type) bool {
	return p.currToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t token.Type) bool {

	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}
