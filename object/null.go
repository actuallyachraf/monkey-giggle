package object

// Null represents the null type (absence of value) a better idea would
// be to use Option instead I think.
type Null struct{}

// Inspect implements the object interface
func (n *Null) Inspect() string {
	return "null"
}

// Type returns the integer type enum.
func (n *Null) Type() Type {
	return NULL
}
