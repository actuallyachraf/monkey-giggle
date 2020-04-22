// Package ast implement an abstract syntax tree api, the tree is built
// using recursive descent parsing.
package ast

import (
	"bytes"

	"github.com/actuallyachraf/monkey-giggle/token"
)

// Node describes a node in the ast.
type Node interface {
	TokenLiteral() token.Literal
	String() string
}

// Statement describes a statement node in the ast, statements are declarations
// and don't produce values.
type Statement interface {
	Node
	statementNode()
}

// Expression describes an expression node in the ast, expressions are value
// producing declarations.
type Expression interface {
	Node
	expressionNode()
}

// Program describes the root node of every ast produced by the parser,
// essentially a program is a sequence of statements which are represented
// by statement nodes in the ast.
type Program struct {
	Statements []Statement
}

// TokenLiteral returns the token literal at the current node.
func (p *Program) TokenLiteral() token.Literal {

	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

// String implements the stringer interface
func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())

	}

	return out.String()
}
