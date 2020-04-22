package lexer

import (
	"testing"

	"github.com/actuallyachraf/monkey-giggle/token"
)

func TestLexer(t *testing.T) {

	t.Run("TestBasicToken", func(t *testing.T) {

		input := `=+-*/%(){},;`

		tests := []struct {
			expectedType    token.Type
			expectedLiteral token.Literal
		}{
			{token.ASSIGN, "="},
			{token.ADD, "+"},
			{token.SUB, "-"},
			{token.MUL, "*"},
			{token.DIV, "/"},
			{token.MOD, "%"},
			{token.LPAREN, "("},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RBRACE, "}"},
			{token.COMMA, ","},
			{token.SEMICOLON, ";"},
		}
		l := New(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - wrong token type : expected %q, got %q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - wrong token literal : expected %q, got %q", i, tt.expectedLiteral, tok.Literal)

			}
		}
	})
	t.Run("TestTokenChain", func(t *testing.T) {

		input := `let five = 5;
				  let ten = 10;
				  let add = fn(x,y){
					  x + y;
				  };
				  let Result = add(five,ten);
				  !-/*5;
				  5 < 10 > 5;
				  if (5 < 10 ){
					  return true;
				  } else {
					  return false;
				  }
				  10 == 10;
				  10 != 9;
				  5 <= 5;
				  6 >= 6;
				  `
		tests := []struct {
			expectedType    token.Type
			expectedLiteral token.Literal
		}{
			{token.LET, "let"},
			{token.IDENT, "five"},
			{token.ASSIGN, "="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "ten"},
			{token.ASSIGN, "="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "add"},
			{token.ASSIGN, "="},
			{token.FUNCTION, "fn"},
			{token.LPAREN, "("},
			{token.IDENT, "x"},
			{token.COMMA, ","},
			{token.IDENT, "y"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.IDENT, "x"},
			{token.ADD, "+"},
			{token.IDENT, "y"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.SEMICOLON, ";"},
			{token.LET, "let"},
			{token.IDENT, "Result"},
			{token.ASSIGN, "="},
			{token.IDENT, "add"},
			{token.LPAREN, "("},
			{token.IDENT, "five"},
			{token.COMMA, ","},
			{token.IDENT, "ten"},
			{token.RPAREN, ")"},
			{token.SEMICOLON, ";"},
			{token.BANG, "!"},
			{token.SUB, "-"},
			{token.DIV, "/"},
			{token.MUL, "*"},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.INT, "5"},
			{token.LT, "<"},
			{token.INT, "10"},
			{token.GT, ">"},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.IF, "if"},
			{token.LPAREN, "("},
			{token.INT, "5"},
			{token.LT, "<"},
			{token.INT, "10"},
			{token.RPAREN, ")"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.TRUE, "true"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.ELSE, "else"},
			{token.LBRACE, "{"},
			{token.RETURN, "return"},
			{token.FALSE, "false"},
			{token.SEMICOLON, ";"},
			{token.RBRACE, "}"},
			{token.INT, "10"},
			{token.EQ, "=="},
			{token.INT, "10"},
			{token.SEMICOLON, ";"},
			{token.INT, "10"},
			{token.NEQ, "!="},
			{token.INT, "9"},
			{token.SEMICOLON, ";"},
			{token.INT, "5"},
			{token.LE, "<="},
			{token.INT, "5"},
			{token.SEMICOLON, ";"},
			{token.INT, "6"},
			{token.GE, ">="},
			{token.INT, "6"},
			{token.SEMICOLON, ";"},
			{token.EOF, ""},
		}
		l := New(input)

		for i, tt := range tests {
			tok := l.NextToken()
			if tok.Type != tt.expectedType {
				t.Fatalf("tests[%d] - wrong token type : expected %q, got %q", i, tt.expectedType, tok.Type)
			}

			if tok.Literal != tt.expectedLiteral {
				t.Fatalf("tests[%d] - wrong token literal : expected %q, got %q", i, tt.expectedLiteral, tok.Literal)

			}
		}
	})
}
