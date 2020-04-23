package eval

import "github.com/actuallyachraf/monkey-giggle/object"

var builtins = map[string]*object.BuiltIn{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, expected %d got %d", 1, len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"head": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, expected %d got %d", 1, len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newError("argument to `push` not supported got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			if len(arr.Elements) > 0 {
				return arr.Elements[0]
			}
			return NULL
		},
	},
	"tail": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, expected %d got %d", 1, len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newError("argument to `push` not supported got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				newElems := make([]object.Object, length-1, length-1)
				copy(newElems, arr.Elements[1:length])
				return &object.Array{Elements: newElems}
			}
			return NULL
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments, expected %d got %d", 1, len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newError("argument to `push` not supported got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			if length > 0 {
				return arr.Elements[length-1]
			}
			return NULL
		},
	},
	"append": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments, expected %d got %d", 2, len(args))
			}
			if args[0].Type() != object.ARRAY {
				return newError("argument to `push` not supported got %s", args[0].Type())
			}
			arr := args[0].(*object.Array)
			length := len(arr.Elements)
			newElems := make([]object.Object, length+1, length+1)
			copy(newElems, arr.Elements)
			newElems[length] = args[1]
			return &object.Array{Elements: newElems}

		},
	},
}
