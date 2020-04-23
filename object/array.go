package object

import (
	"bytes"
	"strings"
)

// Array wraps native arrays
type Array struct {
	Elements []Object
}

// Type implements the object interface
func (a *Array) Type() Type {
	return ARRAY
}

// Inspect implements the object interface
func (a *Array) Inspect() string {

	var out bytes.Buffer

	elements := []string{}

	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()

}
