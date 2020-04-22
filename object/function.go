package object

import (
	"bytes"
	"strings"

	"github.com/actuallyachraf/monkey-giggle/ast"
)

// Function represents a function literal each function gets a local environment.
// to enforce scope rules and prevent variable shadowing.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type implements the Object interface.
func (f *Function) Type() Type {
	return FUNCTION
}

// Inspect implements the Object interface.
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range params {
		params = append(params, p)
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
