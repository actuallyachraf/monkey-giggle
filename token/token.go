package token

// Type encodes the type of the token
type Type string

// Literal encodes a literal value
type Literal string

// Token represents the actual token holds the type and it's literal representation.
type Token struct {
	Type
	Literal
}

// New creates a new token instance
func New(tok Type, ch byte) Token {
	return Token{
		Type:    tok,
		Literal: Literal(ch),
	}
}

// NewLiteral creates a new token with a literal
func NewLiteral(tok Type, lit string) Token {
	return Token{
		Type:    tok,
		Literal: Literal(lit),
	}
}
