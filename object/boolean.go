package object

import "fmt"

// Boolean represents a boolean value
type Boolean struct {
	Value bool
}

// Inspect implements the object interface
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

// Type returns the integer type enum.
func (b *Boolean) Type() Type {
	return BOOLEAN
}
