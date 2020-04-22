package object

// ReturnValue wraps values that are return expressions.
type ReturnValue struct {
	Value Object
}

// Type implements the object interface
func (rv *ReturnValue) Type() Type {
	return RETURN
}

// Inspect implements the ojbect interface
func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}
