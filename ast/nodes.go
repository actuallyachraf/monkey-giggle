package ast

import (
	"bytes"
	"strings"

	"github.com/actuallyachraf/monkey-giggle/token"
)

// nodes.go implements Node for various declarations.

// LetStatement implements the Node interface for let statements
type LetStatement struct {
	Token token.Token // Let token
	Name  *Identifier // Name of the identifier used to hold the left-value expression
	Value Expression  // The value held by this identifier
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral implements the interface and returns the token literal.
func (ls *LetStatement) TokenLiteral() token.Literal {
	return ls.Token.Literal
}

// Stringer implements the stringer interface.
func (ls *LetStatement) String() string {

	var out bytes.Buffer

	out.WriteString(string(ls.TokenLiteral()) + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")
	return out.String()
}

// Identifier implements the Node interface for identifier declarations.
type Identifier struct {
	Token token.Token
	Value token.Literal
}

func (i *Identifier) expressionNode() {}

// TokenLiteral implements the interface and returns the identifier literal.
func (i *Identifier) TokenLiteral() token.Literal {
	return i.Token.Literal
}

// String implements the stringer interface.
func (i *Identifier) String() string {
	return string(i.Value)
}

// ReturnStatement implements the Node interface for identifier declarations.
type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}

// TokenLiteral implements the interface and returns the identifier literal.
func (rs *ReturnStatement) TokenLiteral() token.Literal {
	return rs.Token.Literal
}

// String implements the stringer interface.
func (rs *ReturnStatement) String() string {

	var out bytes.Buffer

	out.WriteString(string(rs.TokenLiteral()) + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// ExpressionStatement implements the interface for parsing expressions.
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

// TokenLiteral implements the interface and returns the Expression literal.
func (es *ExpressionStatement) TokenLiteral() token.Literal {
	return es.Token.Literal
}

// String implements the stringer interface.
func (es *ExpressionStatement) String() string {

	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// IntegerLiteral represents a literal integer expression
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) expressionNode() {}

// TokenLiteral implements the interface and returns the Expression literal.
func (il *IntegerLiteral) TokenLiteral() token.Literal {
	return il.Token.Literal
}

// String implements the stringer interface.
func (il *IntegerLiteral) String() string {
	return string(il.Token.Literal)
}

// StringLiteral represents a literal string
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode() {}

// TokenLiteral implements the node interface
func (sl *StringLiteral) TokenLiteral() token.Literal {
	return sl.Token.Literal
}

// ArrayLiteral represents arrays
type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

// TokenLiteral implements the node interface
func (al *ArrayLiteral) TokenLiteral() token.Literal {
	return al.Token.Literal
}

// String implements the stringer interface
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}

	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// String implements the stringer interface
func (sl *StringLiteral) String() string {
	return string(sl.Token.Literal)
}

// PrefixExpression represents prefixed expressions.
type PrefixExpression struct {
	Token    token.Token
	Operator token.Literal
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}

// TokenLiteral implements the interface of node and returns the literal prefix expression.
func (pe *PrefixExpression) TokenLiteral() token.Literal {
	return pe.Token.Literal
}

// String implements the stringer interface
func (pe *PrefixExpression) String() string {

	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(string(pe.Token.Literal))
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()

}

// InfixExpression represents prefixed expressions.
type InfixExpression struct {
	Token    token.Token
	Operator token.Literal
	Right    Expression
	Left     Expression
}

func (ie *InfixExpression) expressionNode() {}

// TokenLiteral implements the interface of node and returns the literal prefix expression.
func (ie *InfixExpression) TokenLiteral() token.Literal {
	return ie.Token.Literal
}

// String implements the stringer interface
func (ie *InfixExpression) String() string {

	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(string(ie.Left.String()) + " ")
	out.WriteString(string(ie.Operator) + " ")
	out.WriteString(string(ie.Right.String()))

	out.WriteString(")")

	return out.String()

}

// BooleanLiteral represents a literal boolean expression.
type BooleanLiteral struct {
	Token token.Token
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}

// TokenLiteral implements the interface and returns the Expression literal.
func (bl *BooleanLiteral) TokenLiteral() token.Literal {
	return bl.Token.Literal
}

// String implements the stringer interface.
func (bl *BooleanLiteral) String() string {
	return string(bl.Token.Literal)
}

// IfExpression represents a conditional if expression.
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode() {}

// TokenLiteral implements the expression interface and returns the literal.
func (ie *IfExpression) TokenLiteral() token.Literal {
	return ie.Token.Literal
}

// String implements the stringer interface
func (ie *IfExpression) String() string {

	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(ie.Condition.String())
	out.WriteString(" ")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}

	return out.String()
}

// BlockStatement represents a sequence of statements that execute within
// a conditional branch.
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

// TokenLiteral implements the interface and returns the token literal
func (bs *BlockStatement) TokenLiteral() token.Literal {
	return bs.Token.Literal
}

// String implements the stringer interface.
func (bs *BlockStatement) String() string {

	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

// FunctionLiteral represents nodes for expressions of the type fn <params> <block>
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

// TokenLiteral implements the interface and returns the token literal fn.
func (fl *FunctionLiteral) TokenLiteral() token.Literal {
	return fl.Token.Literal
}

// String implements the stringer interface.
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}

	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(string(fl.TokenLiteral()))
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

// CallExpression represents function calls.
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

// TokenLiteral implements the expressionNode interface returns the token literal.
func (ce *CallExpression) TokenLiteral() token.Literal {
	return ce.Token.Literal
}

// String implements the stringer interface.
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}

	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()

}

// IndexExpression represents indexing expressions for index accessible ds.
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode() {}

// TokenLiteral implements the node interface
func (ie *IndexExpression) TokenLiteral() token.Literal {
	return ie.Token.Literal
}

// String implements the stringer interface
func (ie *IndexExpression) String() string {

	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}

// HashmapLiteral represents a hashmap
type HashmapLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashmapLiteral) expressionNode() {}

// TokenLiteral implements the node interface
func (hl *HashmapLiteral) TokenLiteral() token.Literal {
	return hl.Token.Literal
}

// String implements the stringer interface
func (hl *HashmapLiteral) String() string {
	var out bytes.Buffer

	pairs := []string{}

	for k, v := range hl.Pairs {
		pairs = append(pairs, k.String()+":"+v.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()

}
