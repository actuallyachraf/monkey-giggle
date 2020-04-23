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

// BuiltInFunc defines functions that are part of the language and operate
// on native objects.
type BuiltInFunc func(args ...Object) Object

// BuiltIn describes builtin function objects
type BuiltIn struct {
	Fn BuiltInFunc
}

// Type implements the object interface
func (b *BuiltIn) Type() Type {
	return BUILTIN
}

// Inspect implements the object interface
func (b *BuiltIn) Inspect() string {
	return "built-in function"
}
