package token

// types.go holds the defined tokens of the language.defined

const (
	// ILLEGAL denotes an unknown (illegal) token type
	ILLEGAL Type = "ILLEGAL"
	// EOF marks end of file attributes
	EOF = "EOF"

	// IDENT denotes the identifier which is the "name" used for a token foo, bar x,y,z and such.
	IDENT = "IDENT"
	// INT denotes the integer type
	INT = "INT"
	// STRING denotes the string type
	STRING = "STRING"

	// Operators

	// ASSIGN denotes the assignment operator token.
	ASSIGN = "ASSIGN"
	// ADD denotes addition token
	ADD = "ADD"
	// SUB denotes substraction token
	SUB = "SUB"
	// MUL denotes multiplication token
	MUL = "MUL"
	// DIV denotes integer division token
	DIV = "DIV"
	// MOD denotes modulo operation token
	MOD = "MOD"
	// LT denotes lesser than
	LT = "<"
	// GT denotes greater than
	GT = ">"
	// EQ denotes equality test
	EQ = "=="
	// NEQ denotes the not equal test
	NEQ = "!="
	// LE denotes the lesser than or equal
	LE = "<="
	// GE denotes the greater than or equal
	GE = ">="

	// BANG denotes the bang token
	BANG = "!"

	// Delimiters are used to separate text representations

	// SEMICOLON represents the semicolon delimiter for scopes
	SEMICOLON = ";"
	// COMMA represents the comma delimiter for values
	COMMA = ","
	// LPAREN represents a left parenthesis
	LPAREN = "("
	// RPAREN represents a right parenthesis
	RPAREN = ")"
	// LBRACE represents a left curly brace
	LBRACE = "{"
	// RBRACE represents a right curly brace
	RBRACE = "}"
	// LBRACKET represents a left bracket
	LBRACKET = "["
	// RBRACKET represents a right bracked
	RBRACKET = "]"

	// Keywords are special words that can't be used as identifiers
	// such as the "let" declaration for values, "fn" for functions.

	// FUNCTION represents a new function declaration
	FUNCTION = "FUNCTION"
	// LET represents the value declaration token
	LET = "LET"
	// TRUE represents the boolean value true
	TRUE = "TRUE"
	// FALSE represents the boolean value false
	FALSE = "FALSE"
	// IF represents the conditional if branch
	IF = "IF"
	// ELSE represents the conditional else branch
	ELSE = "ELSE"
	// RETURN represents the return instruction
	RETURN = "RETURN"
)
