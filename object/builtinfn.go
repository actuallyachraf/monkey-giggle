package object

import "fmt"

// Builtins indexes built-in functions
var Builtins = []struct {
	Name string
	Fn   *BuiltIn
}{
	{
		Name: "len",
		Fn: &BuiltIn{
			func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments, expected %d got %d", 1, len(args))
				}
				switch arg := args[0].(type) {
				case *String:
					return &Integer{Value: int64(len(arg.Value))}
				case *Array:
					return &Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
				}

			},
		},
	}, {
		Name: "head",

		Fn: &BuiltIn{func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, expected %d got %d", 1, len(args))
			}
			if args[0].Type() != ARRAY {
				return newError("argument to `head` not supported got %s", args[0].Type())
			}
			arr := args[0].(*Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return nil
		},
		},
	}, {
		Name: "tail",
		Fn: &BuiltIn{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments, expected %d got %d", 1, len(args))
				}
				if args[0].Type() != ARRAY {
					return newError("argument to `push` not supported got %s", args[0].Type())
				}
				arr := args[0].(*Array)
				length := len(arr.Elements)
				if length > 0 {
					newElems := make([]Object, length-1, length-1)
					copy(newElems, arr.Elements[1:length])
					return &Array{Elements: newElems}
				}
				return nil
			},
		},
	}, {
		Name: "last",
		Fn: &BuiltIn{
			Fn: func(args ...Object) Object {
				if len(args) != 1 {
					return newError("wrong number of arguments, expected %d got %d", 1, len(args))
				}
				if args[0].Type() != ARRAY {
					return newError("argument to `push` not supported got %s", args[0].Type())
				}
				arr := args[0].(*Array)
				length := len(arr.Elements)
				if length > 0 {
					return arr.Elements[length-1]
				}
				return nil
			},
		},
	}, {
		Name: "append",
		Fn: &BuiltIn{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments, expected %d got %d", 2, len(args))
				}
				if args[0].Type() != ARRAY {
					return newError("argument to `append` must be ARRAY got %s", args[0].Type())
				}
				arr := args[0].(*Array)
				length := len(arr.Elements)
				newElems := make([]Object, length+1, length+1)
				copy(newElems, arr.Elements)
				newElems[length] = args[1]
				return &Array{Elements: newElems}

			},
		},
	}, {
		Name: "concat",
		Fn: &BuiltIn{
			Fn: func(args ...Object) Object {
				if len(args) != 2 {
					return newError("wrong number of arguments, expected %d got %d", 2, len(args))
				}
				if args[0].Type() != ARRAY || args[1].Type() != ARRAY {
					return newError("argument to `concat` must be ARRAY got %s", args[0].Type())
				}
				arr1 := args[0].(*Array)
				arr2 := args[1].(*Array)
				length := len(arr1.Elements) + len(arr2.Elements)
				newElemens := make([]Object, 0, length)
				newElemens = append(newElemens, arr1.Elements...)
				newElemens = append(newElemens, arr2.Elements...)

				return &Array{Elements: newElemens}
			},
		},
	},
}

// newError creates a new error message
func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

// GetBuiltInByName fetches a built-in func by its name
func GetBuiltInByName(name string) *BuiltIn {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Fn
		}
	}
	return nil
}
