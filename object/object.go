package object

const (
	// INTEGER are wrapped 64 integer values.
	INTEGER = "INTEGER"
	// BOOLEAN represents a wrapped bool value.
	BOOLEAN = "BOOLEAN"
	// NULL represents a null object which is just an empty struct
	NULL = "NULL"
	// RETURN represents a value that's to be "returned"
	RETURN = "RETURN"
	// FUNCTION represents a function.
	FUNCTION = "FUNCTION"
	// ERROR represents errors that happen during runtime
	ERROR = "ERROR"
)

// Type represents the type of a given object.
type Type string

// Object is an interface that the host (Go) wrapped types for the language
// implement.
type Object interface {
	Type() Type
	Inspect() string
}
