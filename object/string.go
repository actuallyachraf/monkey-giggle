package object

// String represents the language String type.
type String struct {
	Value string
}

// Inspect implements the object interface
func (i *String) Inspect() string {
	return i.Value
}

// Type returns the String type enum.
func (i *String) Type() Type {
	return STRING
}
