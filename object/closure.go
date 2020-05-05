package object

import "fmt"

// Closure represents compiled closures
type Closure struct {
	Fn            *CompiledFunction
	FreeVariables []Object
}

// Type implements the object interface
func (c *Closure) Type() Type {
	return CLOSURE
}

// Inspect implements the object interface
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", c)
}
