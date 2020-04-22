package object

// Error represents an error (a string with the error message).
type Error struct {
	Message string
}

// Type implements the Object interface.
func (e *Error) Type() Type {
	return ERROR
}

// Inspect implements the Object interface.
func (e *Error) Inspect() string {
	return "ERROR :" + e.Message
}
