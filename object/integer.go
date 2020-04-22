package object

import "fmt"

// Integer represents the language integer type.
type Integer struct {
	Value int64
}

// Inspect implements the object interface
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

// Type returns the integer type enum.
func (i *Integer) Type() Type {
	return INTEGER
}
