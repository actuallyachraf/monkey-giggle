// Package lexer implements lexical analysis for the monkey source code
// the goal is to simply read a string representation of code
// and turn it into a token representation.
package lexer

import "github.com/actuallyachraf/monkey-giggle/token"

// Lexer represents a lexical analysis engine.
type Lexer struct {
	input   string // represents the input string (TODO:replace with io.Reader)
	pos     int    // current position in input (current char)
	readPos int    // next position in input
	ch      byte   // current char
}

// New creates a new instance of lexer.
func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}
	l.readChar()
	return l
}

// readChar reads a single byte from the string and update positions.a
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}

	l.pos = l.readPos
	l.readPos++
}

// NextToken parses the and returns the next token in the input.
func (l *Lexer) NextToken() token.Token {

	var tok token.Token
	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.NewLiteral(token.EQ, string(ch)+string(l.ch))
		} else {
			tok = token.New(token.ASSIGN, l.ch)
		}
	case '+':
		tok = token.New(token.ADD, l.ch)
	case '-':
		tok = token.New(token.SUB, l.ch)
	case '*':
		tok = token.New(token.MUL, l.ch)
	case '/':
		tok = token.New(token.DIV, l.ch)
	case '%':
		tok = token.New(token.MOD, l.ch)
	case ';':
		tok = token.New(token.SEMICOLON, l.ch)
	case ',':
		tok = token.New(token.COMMA, l.ch)
	case '(':
		tok = token.New(token.LPAREN, l.ch)
	case ')':
		tok = token.New(token.RPAREN, l.ch)
	case '{':
		tok = token.New(token.LBRACE, l.ch)
	case '}':
		tok = token.New(token.RBRACE, l.ch)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.NewLiteral(token.NEQ, string(ch)+string(l.ch))
		} else {
			tok = token.New(token.BANG, l.ch)
		}

	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.NewLiteral(token.LE, string(ch)+string(l.ch))
		} else {
			tok = token.New(token.LT, l.ch)
		}

	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.NewLiteral(token.GE, string(ch)+string(l.ch))
		} else {
			tok = token.New(token.GT, l.ch)
		}
	case 0:
		tok = token.Token{Type: token.EOF, Literal: token.Literal("")}

	default:
		if isLetter(l.ch) {
			tok.Literal = token.Literal(l.readIdentifier())
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = token.Literal(l.readNumber())
			return tok
		} else {
			tok = token.New(token.ILLEGAL, l.ch)
		}
	}

	// move to the next char
	l.readChar()

	return tok
}

// readIdentifier reads the next identifier
func (l *Lexer) readIdentifier() string {

	pos := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.pos]
}

// skipWhitespace is used to escape whitespace between keywords and literal identifiers.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readNumber is used to read a whole number (composed of multiplie digits).
func (l *Lexer) readNumber() string {
	pos := l.pos

	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.pos]
}

// peekChar is used to read past the current char to determine multi char tokens
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]

}

// isLetter checks whether the current char is valid ASCII letter
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks whether the current char is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}
