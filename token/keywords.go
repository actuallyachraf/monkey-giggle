package token

// file keywords.go define the language proper keywords.

var keywords = map[Literal]Type{
	"fn":     FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent checks whether an identifier string is a keyword or not.
func LookupIdent(ident Literal) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
