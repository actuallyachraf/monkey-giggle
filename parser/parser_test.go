package parser

import (
	"errors"
	"fmt"
	"testing"

	"github.com/actuallyachraf/monkey-giggle/ast"
	"github.com/actuallyachraf/monkey-giggle/lexer"
	"github.com/actuallyachraf/monkey-giggle/token"
)

func TestParser(t *testing.T) {

	t.Run("TestParseLetStatement", func(t *testing.T) {
		input := `

		let x = 5;
		let y = 10;
		let foobar = 9876543210;
		`
		badInput := `

		let x 5;
		let  = 10;
		let  9876543210;
		`

		l := lexer.New(input)
		p := New(l)

		program := p.Parse()

		checkParserError(t, p)

		if program == nil {
			t.Fatal("Parse() returned nil")
		}
		if len(program.Statements) != 3 {
			t.Fatal("Parse() error : expected 3 statments got ", len(program.Statements), " instead")
		}
		tests := []struct {
			expectedIdentifier token.Literal
		}{
			{"x"},
			{"y"},
			{"foobar"},
		}

		for i, tt := range tests {
			stmt := program.Statements[i]
			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}
		}

		l = lexer.New(badInput)
		p = New(l)

		program = p.Parse()

		if checkParserError(t, p) == nil {
			t.Fatal("Parser should fail on bad input")
		}
	})
	t.Run("TestParseLetStatementWithExpressions", func(t *testing.T) {
		tests := []struct {
			input              string
			expectedIdentifier token.Literal
			expectedValue      interface{}
		}{
			{"let x = 5;", "x", 5},
			{"let y = true;", "y", true},
			{"let foobar = y;", "foobar", "y"},
			{"let ourFunction = 5;", "ourFunction", 5},
		}

		for _, tt := range tests {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.Parse()
			checkParserError(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program.Statements does not contain 1 statements. got=%d",
					len(program.Statements))
			}

			stmt := program.Statements[0]
			if !testLetStatement(t, stmt, tt.expectedIdentifier) {
				return
			}

			val := stmt.(*ast.LetStatement).Value
			if !testLiteralExpression(t, val, tt.expectedValue) {
				return
			}
		}
	})
	t.Run("TestParseReturnStatement", func(t *testing.T) {

		input := `
		return 5;
		return 10;
		return 935834;
		`
		l := lexer.New(input)
		p := New(l)

		program := p.Parse()

		checkParserError(t, p)

		if program == nil {
			t.Fatal("Parse() returned nil")
		}
		if len(program.Statements) != 3 {
			t.Fatal("Parse() error : expected 3 statments got ", len(program.Statements), " instead")
		}
		for _, stmt := range program.Statements {
			returnStmt, ok := stmt.(*ast.ReturnStatement)
			if !ok {
				t.Errorf("statement not *ast.ReturnStatement got %T", returnStmt)
				continue
			}
			if returnStmt.TokenLiteral() != "return" {
				t.Errorf("returnStmt.TokenLiteral not return got %s", returnStmt.TokenLiteral())
			}
		}
	})
	t.Run("TestParseIdentifierExpression", func(t *testing.T) {

		input := "foobar"

		l := lexer.New(input)
		p := New(l)
		program := p.Parse()
		checkParserError(t, p)
		if len(program.Statements) != 1 {
			t.Fatal("program has not enough statement expected 1 got ", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program has wrong statement expected ExpressionStatement got %T", program.Statements[0])
		}
		ident, ok := stmt.Expression.(*ast.Identifier)
		if !testIdentifier(t, ident, input) {
			t.Fatal("program failed to parse identifier")
		}
	})
	t.Run("TestParseIntegerLiteralExpression", func(t *testing.T) {
		input := "5"

		l := lexer.New(input)
		p := New(l)
		program := p.Parse()
		checkParserError(t, p)
		if len(program.Statements) != 1 {
			t.Fatal("program has not enough statement expected 1 got ", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program has wrong statement expected ExpressionStatement got %T", program.Statements[0])
		}
		literal, ok := stmt.Expression.(*ast.IntegerLiteral)
		if !ok {
			t.Fatalf("program has wrong expression expected *ast.Identifier got %T", stmt.Expression)
		}
		if literal.Value != 5 {
			t.Errorf("program has wrong identifier value expected %d got %d", 5, literal.Value)
		}
		if literal.TokenLiteral() != "5" {
			t.Errorf("program has wrong token literal expected %s got %s", "5", literal.TokenLiteral())
		}
	})
	t.Run("TestParsePrefixExpression", func(t *testing.T) {
		prefixTests := []struct {
			input    string
			operator string
			value    interface{}
		}{
			{"!5", "!", 5},
			{"-15", "-", 15},
			{"!true", "!", true},
			{"!false", "!", false},
		}

		for _, tt := range prefixTests {

			l := lexer.New(tt.input)
			p := New(l)
			program := p.Parse()
			checkParserError(t, p)

			if len(program.Statements) != 1 {
				t.Fatal("program has not enough statement expected 1 got ", len(program.Statements))
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program has wrong statement expected ExpressionStatement got %T", program.Statements[0])
			}
			exp, ok := stmt.Expression.(*ast.PrefixExpression)
			if !ok {
				t.Fatalf("program has wrong expression expected PrefixExpression got %T", program.Statements[0])
			}

			if string(exp.Operator) != tt.operator {
				t.Fatalf("program has wrong expression operator expected %s got %s", exp.Operator, tt.operator)
			}
			if !testLiteralExpression(t, exp.Right, tt.value) {
				return
			}
		}
	})
	t.Run("TestParseInfixExpression", func(t *testing.T) {
		infixTests := []struct {
			input    string
			leftVal  interface{}
			operator string
			rightVal interface{}
		}{
			{"5 + 5;", 5, "+", 5},
			{"5 - 5;", 5, "-", 5},
			{"5 * 5;", 5, "*", 5},
			{"5 / 5;", 5, "/", 5},
			{"5 % 5;", 5, "%", 5},
			{"5 < 5;", 5, "<", 5},
			{"5 <= 5;", 5, "<=", 5},
			{"5 >= 5;", 5, ">=", 5},
			{"5 == 5;", 5, "==", 5},
			{"5 != 5;", 5, "!=", 5},
			{"true == true", true, "==", true},
			{"true != false", true, "!=", false},
			{"false == false", false, "==", false},
		}

		for _, tt := range infixTests {

			l := lexer.New(tt.input)
			p := New(l)
			program := p.Parse()
			checkParserError(t, p)

			if len(program.Statements) != 1 {
				t.Fatal("program has not enough statement expected 1 got ", len(program.Statements))
			}
			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program has wrong statement expected ExpressionStatement got %T", program.Statements[0])
			}
			exp, ok := stmt.Expression.(*ast.InfixExpression)

			if !testInfixExpression(t, exp, tt.leftVal, tt.operator, tt.rightVal) {
				t.Fatalf("program has failed to parse infix expression")
			}

		}
	})
	t.Run("TestParseOperatorPrecedence", func(t *testing.T) {
		tests := []struct {
			input    string
			expected string
		}{
			{
				"-a * b",
				"((-a) * b)",
			},
			{
				"!-a",
				"(!(-a))",
			},
			{
				"a + b + c",
				"((a + b) + c)",
			},
			{
				"a + b - c",
				"((a + b) - c)",
			},
			{
				"a * b * c",
				"((a * b) * c)",
			},
			{
				"a * b / c",
				"((a * b) / c)",
			},
			{
				"a + b / c",
				"(a + (b / c))",
			},
			{
				"a + b * c + d / e - f",
				"(((a + (b * c)) + (d / e)) - f)",
			},
			{
				"3 + 4; -5 * 5",
				"(3 + 4)((-5) * 5)",
			},
			{
				"5 > 4 == 3 < 4",
				"((5 > 4) == (3 < 4))",
			},
			{
				"5 < 4 != 3 > 4",
				"((5 < 4) != (3 > 4))",
			},
			{
				"3 + 4 * 5 == 3 * 1 + 4 * 5",
				"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
			},
			{
				"true",
				"true",
			},
			{
				"false",
				"false",
			},
			{
				"3 > 5 == false",
				"((3 > 5) == false)",
			},
			{
				"3 < 5 == true",
				"((3 < 5) == true)",
			},
			{
				"1 + (2 + 3) + 4",
				"((1 + (2 + 3)) + 4)",
			},
			{
				"(5 + 5) * 2",
				"((5 + 5) * 2)",
			},
			{
				"2 / (5 + 5)",
				"(2 / (5 + 5))",
			},
			{
				"(5 + 5) * 2 * (5 + 5)",
				"(((5 + 5) * 2) * (5 + 5))",
			},
			{
				"-(5 + 5)",
				"(-(5 + 5))",
			},
			{
				"!(true == true)",
				"(!(true == true))",
			},
			{
				"a + add(b * c) + d",
				"((a + add((b * c))) + d)",
			},
			{
				"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
				"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
			},
			{
				"add(a + b + c * d / f + g)",
				"add((((a + b) + ((c * d) / f)) + g))",
			},
		}

		for _, tt := range tests {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.Parse()
			checkParserError(t, p)

			actual := program.String()
			if actual != tt.expected {
				t.Errorf("expected=%q, got=%q", tt.expected, actual)
			}
		}
	})
	t.Run("TestParseBooleanLiteralExpression", func(t *testing.T) {
		tests := []struct {
			input           string
			expectedBoolean bool
		}{
			{"true;", true},
			{"false;", false},
		}

		for _, tt := range tests {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.Parse()
			checkParserError(t, p)

			if len(program.Statements) != 1 {
				t.Fatalf("program has not enough statements. got=%d",
					len(program.Statements))
			}

			stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
					program.Statements[0])
			}

			boolean, ok := stmt.Expression.(*ast.BooleanLiteral)
			if !ok {
				t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
			}
			if boolean.Value != tt.expectedBoolean {
				t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
					boolean.Value)
			}
		}
	})
	t.Run("TestParseIfExpression", func(t *testing.T) {
		input := `if (x < y) { x }`

		l := lexer.New(input)
		p := New(l)
		program := p.Parse()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has wrong number of statements expected 1 got %d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program has wrong statement type expected *ast.ExpressionStatement got %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("program has wrong expression type exepcted *ast.IfExpression got %T", stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}

		if len(exp.Consequence.Statements) != 1 {
			t.Fatalf("program has wrong number of statements expected 1 got %d", len(exp.Consequence.Statements))

		}

		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program has wrong statement type expected *ast.ExpressionStatement got %T", exp.Consequence.Statements[0])
		}

		if !testIdentifier(t, consequence.Expression, "x") {
			return
		}
		if exp.Alternative != nil {
			t.Fatalf("program has wrong alternative expected nil got %+v", exp.Alternative)
		}
	})
	t.Run("TestParseIfElseExpression", func(t *testing.T) {
		input := `if (x < y) { x } else { y }`

		l := lexer.New(input)
		p := New(l)
		program := p.Parse()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Body does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.IfExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
		}

		if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
			return
		}

		if len(exp.Consequence.Statements) != 1 {
			t.Errorf("consequence is not 1 statements. got=%d\n",
				len(exp.Consequence.Statements))
		}

		consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
				exp.Consequence.Statements[0])
		}

		if !testIdentifier(t, consequence.Expression, "x") {
			return
		}

		if len(exp.Alternative.Statements) != 1 {
			t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
				len(exp.Alternative.Statements))
		}

		alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
				exp.Alternative.Statements[0])
		}

		if !testIdentifier(t, alternative.Expression, "y") {
			return
		}
	})
	t.Run("TestParseFunctionLiterals", func(t *testing.T) {
		input := `fn(x, y) { x + y; }`

		l := lexer.New(input)
		p := New(l)
		program := p.Parse()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Body does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		function, ok := stmt.Expression.(*ast.FunctionLiteral)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
				stmt.Expression)
		}

		if len(function.Parameters) != 2 {
			t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
				len(function.Parameters))
		}

		testLiteralExpression(t, function.Parameters[0], "x")
		testLiteralExpression(t, function.Parameters[1], "y")

		if len(function.Body.Statements) != 1 {
			t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
				len(function.Body.Statements))
		}

		bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
				function.Body.Statements[0])
		}

		testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
	})
	t.Run("TestParseFunctionParameters", func(t *testing.T) {
		tests := []struct {
			input          string
			expectedParams []string
		}{
			{input: "fn() {};", expectedParams: []string{}},
			{input: "fn(x) {};", expectedParams: []string{"x"}},
			{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
		}

		for _, tt := range tests {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.Parse()
			checkParserError(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			function := stmt.Expression.(*ast.FunctionLiteral)

			if len(function.Parameters) != len(tt.expectedParams) {
				t.Errorf("length parameters wrong. want %d, got=%d\n",
					len(tt.expectedParams), len(function.Parameters))
			}

			for i, ident := range tt.expectedParams {
				testLiteralExpression(t, function.Parameters[i], ident)
			}
		}
	})
	t.Run("TestParseCallExpression", func(t *testing.T) {
		input := "add(1, 2 * 3, 4 + 5);"

		l := lexer.New(input)
		p := New(l)
		program := p.Parse()
		checkParserError(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, "add") {
			return
		}

		if len(exp.Arguments) != 3 {
			t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
		}

		testLiteralExpression(t, exp.Arguments[0], 1)
		testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
		testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
	})
	t.Run("TestParseCallArguments", func(t *testing.T) {
		tests := []struct {
			input         string
			expectedIdent string
			expectedArgs  []string
		}{
			{
				input:         "add();",
				expectedIdent: "add",
				expectedArgs:  []string{},
			},
			{
				input:         "add(1);",
				expectedIdent: "add",
				expectedArgs:  []string{"1"},
			},
			{
				input:         "add(1, 2 * 3, 4 + 5);",
				expectedIdent: "add",
				expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
			},
		}

		for _, tt := range tests {
			l := lexer.New(tt.input)
			p := New(l)
			program := p.Parse()
			checkParserError(t, p)

			stmt := program.Statements[0].(*ast.ExpressionStatement)
			exp, ok := stmt.Expression.(*ast.CallExpression)
			if !ok {
				t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
					stmt.Expression)
			}

			if !testIdentifier(t, exp.Function, tt.expectedIdent) {
				return
			}

			if len(exp.Arguments) != len(tt.expectedArgs) {
				t.Fatalf("wrong number of arguments. want=%d, got=%d",
					len(tt.expectedArgs), len(exp.Arguments))
			}

			for i, arg := range tt.expectedArgs {
				if exp.Arguments[i].String() != arg {
					t.Errorf("argument %d wrong. want=%q, got=%q", i,
						arg, exp.Arguments[i].String())
				}
			}
		}
	})
}

func testLetStatement(t *testing.T, s ast.Statement, name token.Literal) bool {

	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral got %s instead of let", s.TokenLiteral())
		return false
	}
	letStatement, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("interface breakage expected *ast.LetStatement got %T", s)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("letStatement.Name.Value not %s got %s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("letStatement.Name not %s got %s", name, letStatement.Name)
		return false
	}

	return true
}
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("wrong expression type expected *ast.IntegerLiteral got %T", il)
		return false
	}
	if integ.Value != value {
		t.Errorf("wrong integer literal expected %d got %d", value, integ.Value)
	}

	valueLiteral := fmt.Sprintf("%d", value)

	if integ.String() != valueLiteral {
		t.Errorf("wrong token literal expected %s got %s", valueLiteral, integ.TokenLiteral())
		return false
	}

	return true
}
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier got %T", exp)
		return false
	}
	if string(ident.Value) != value {
		t.Errorf("wrong identifier value expected %s got %s", value, ident.Value)
		return false
	}
	if string(ident.TokenLiteral()) != value {
		t.Errorf("wrong token literal value expected %s got %s", value, ident.TokenLiteral())
		return false
	}

	return true
}
func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of expression not handled got %T", exp)
	return false
}
func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {

	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not Operator expression got %T", opExp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	if string(opExp.Operator) != operator {
		t.Errorf("exp.Operator is not %s got %q", operator, opExp.Operator)
		return false
	}
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if string(bo.TokenLiteral()) != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}
func checkParserError(t *testing.T, p *Parser) error {
	parserErrors := p.Errors()

	if len(parserErrors) == 0 {
		return nil
	}

	t.Logf("parse failed with %d errors", len(parserErrors))

	for _, msg := range parserErrors {
		t.Logf("parser error :%q", msg)
	}

	return errors.New("Parser Failed")
}
