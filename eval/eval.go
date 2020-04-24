package eval

import (
	"fmt"

	"github.com/actuallyachraf/monkey-giggle/ast"
	"github.com/actuallyachraf/monkey-giggle/object"
	"github.com/actuallyachraf/monkey-giggle/token"
)

// Native Object References are used to reference predefined values
// instead of creating new ones.
var (
	// TRUE denotes the boolean true value
	TRUE = &object.Boolean{Value: true}
	// FALSE denotes the boolean false value
	FALSE = &object.Boolean{Value: false}
	// NULL dentoes the null value
	NULL = &object.Null{}
)

// Eval works by tree walking given an ast.Node it evaluates and returns a host
// type value.
func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		if node == nil {
			fmt.Println("Nil node value !")
		}
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(string(node.Name.Value), val)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBoolean(node.Value)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)
	case *ast.HashmapLiteral:
		return evalHashmapLiteral(node, env)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}
	return nil
}

// evalProgram runs down the eval function on program/block statements.
func evalProgram(program *ast.Program, env *object.Environment) object.Object {

	var res object.Object

	for _, s := range program.Statements {
		res = Eval(s, env)

		switch res := res.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}

	return res
}

// evalBlockStatement evaluates block statements
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {

	var res object.Object

	for _, statement := range block.Statements {
		res = Eval(statement, env)

		if res != nil {
			rt := res.Type()
			if rt == object.RETURN || rt == object.ERROR {
				return res
			}
		}
	}

	return res
}

// evalPrefixExpression is called when we need to evaluate a prefixed expression.
func evalPrefixExpression(operator token.Literal, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("unknwon operator :%s%s", operator, right.Type())
	}
}

// evalInfixExpression is called when we need to evaluate infixed expression
// of the type val op val.
func evalInfixExpression(operator token.Literal, left object.Object, right object.Object) object.Object {
	// order matters here if the operator is == we need to evaluate whether it's
	// operating on integer operands before assuming left and right are both boolean.
	// edge case (1 == 2) == true => false
	// this bug is artefact of the fact that we also compare pointer values since the TRUE/FALSE booleans
	// have long-life pointers during repl runtime.
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.INTEGER && right.Type() == object.INTEGER:
		return evalIntegerExpression(operator, left, right)
	case left.Type() == object.STRING && right.Type() == object.STRING:
		return evalStringExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToBoolean(left == right)
	case operator == "!=":
		return nativeBoolToBoolean(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())

	}
}

// evalIntegerExpression evaluates infix expression where both operands are integers.
func evalIntegerExpression(operator token.Literal, left object.Object, right object.Object) object.Object {

	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "<":
		return nativeBoolToBoolean(leftVal < rightVal)
	case ">":
		return nativeBoolToBoolean(leftVal > rightVal)
	case "<=":
		return nativeBoolToBoolean(leftVal <= rightVal)
	case ">=":
		return nativeBoolToBoolean(leftVal >= rightVal)
	case "==":
		return nativeBoolToBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringExpression evaluates infix expression where both operands are strings.
func evalStringExpression(operator token.Literal, left object.Object, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError("unknown operator for string type %s %s %s", operator, leftVal, rightVal)
	}
}

// evalBangOperatorExpression is used to evaluate expressions prefixed by the
// bang operator.
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

// evalMinusOperatorExpression is used to evaluate expressions prefixed by the
// minus operator
func evalMinusOperatorExpression(right object.Object) object.Object {

	if right.Type() != object.INTEGER {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value

	return &object.Integer{Value: -value}
}

// evalExpressions evaluations expressions that appear as arguments.
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {

	var res []object.Object

	for _, e := range exps {
		evaled := Eval(e, env)
		if isError(evaled) {
			return []object.Object{evaled}
		}

		res = append(res, evaled)
	}

	return res
}

// evalIfExpression is used to evaluate conditionnal branches.
func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)

	if isTruth(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL

}

// evalIndexExpression evaluates expressions in the indexing op
func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH:
		return evalHashmapIndexExpression(left, index)
	default:
		return NULL
	}
}

// evalArrayIndexExpression evaluates expression in array indexes
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value
	max := int64(len(arrayObj.Elements) - 1)

	if idx < 0 || idx > max {
		return NULL
	}

	return arrayObj.Elements[idx]
}

// evalHashmapIndexExpression evaluates expression in hashmap keys
func evalHashmapIndexExpression(hash, index object.Object) object.Object {

	hashObj := hash.(*object.HashMap)

	key, ok := index.(object.Hashable)
	if !ok {
		return newError("not valid key for hashmap literal : %s", index.Type())
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}

// evalHashmapLiteral evaluates hashmap literals
func evalHashmapLiteral(node *ast.HashmapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for k, v := range node.Pairs {
		key := Eval(k, env)
		if isError(key) {
			return key
		}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return newError("not valid key for hashmap literal : %s", key.Type())
		}
		val := Eval(v, env)
		if isError(val) {
			return val
		}
		hashed := hashKey.HashKey()
		pairs[hashed] = object.HashPair{Key: key, Value: val}
	}

	return &object.HashMap{Pairs: pairs}
}

// applyFunction fn to a list of arguments
func applyFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaled := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaled)
	case *object.BuiltIn:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

// extendFunctionEnv from current environment
func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(string(param.Value), args[paramIdx])
	}

	return env
}

// unwrapReturnValue unwraps the return value to a return statement
func unwrapReturnValue(obj object.Object) object.Object {
	if retVal, ok := obj.(*object.ReturnValue); ok {
		return retVal.Value
	}

	return obj
}

// nativeBoolToBoolean returns a referenced type instead of creating a new one
// for the Boolean type.
func nativeBoolToBoolean(in bool) *object.Boolean {
	if in {
		return TRUE
	}

	return FALSE
}

// isTruth returns whether a given object is a true expression or not.
func isTruth(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case FALSE:
		return false
	default:
		return true
	}
}

// newError creates a new error message
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// isError checks if a given object is an error
func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR
	}

	return false
}

// evalIdentifier ...
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(string(node.Value)); ok {
		return val
	}
	if builtin, ok := builtins[string(node.Value)]; ok {
		return builtin
	}

	return newError("identifier not found: " + string(node.Value))
}
