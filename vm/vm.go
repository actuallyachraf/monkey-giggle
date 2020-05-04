package vm

import (
	"fmt"

	"github.com/actuallyachraf/monkey-giggle/code"
	"github.com/actuallyachraf/monkey-giggle/compiler"
	"github.com/actuallyachraf/monkey-giggle/object"
)

const (
	// StackSize is the max number of items that can be on the stack.
	StackSize = 2048
	// GlobalsSize is the max number of global binds in a given program.
	GlobalsSize = 65536
	// MaxFrames is the max number of call frames that can be on the stack frame
	MaxFrames = 1024
)

var (
	// True marks truth value
	True = &object.Boolean{Value: true}
	// False marks false value
	False = &object.Boolean{Value: false}
	// Null marks the null value
	Null = &object.Null{}
)

// VM represents a virtual machine.
type VM struct {
	constants []object.Object

	stack []object.Object
	sp    int

	globals []object.Object

	frames      []*Frame
	framesIndex int
}

// New creates a new instance of VM using bytecode to execute.
func New(bytecode compiler.Bytecode) *VM {

	// the program bytecode is considered an entire function and is pushed
	// as part of it's own call frame
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants:   bytecode.Constants,
		globals:     make([]object.Object, GlobalsSize),
		stack:       make([]object.Object, StackSize),
		sp:          0,
		frames:      frames,
		framesIndex: 1,
	}
}

// NewWithGlobalState creates a new instace of VM with a global state
func NewWithGlobalState(bytecode compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s

	return vm
}

// CurrentFrame returns the current call frame
func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

// pushFrame pushes a new call frame to the frame stack
func (vm *VM) pushFrame(f *Frame) {
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

// popFrame removes the last used call frame from the frame stack
func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

// StackTop returns the top most element of the stack.
func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}

	return vm.stack[vm.sp-1]
}

// Run is the main loop that runs a fetch-decode-execute cycle.
func (vm *VM) Run() error {

	var ip int
	var inst code.Instructions
	var op code.OpCode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {

		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		inst = vm.currentFrame().Instructions()
		op = code.OpCode(inst[ip])

		switch op {
		case code.OpPop:
			vm.pop()
		case code.OpConstant:
			constIndex := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpMod, code.OpDiv:
			err := vm.executeBinOp(op)
			if err != nil {
				return err
			}
		case code.OpEqual, code.OpNotEqual, code.OpGreaterOrEqual, code.OpGreaterThan:
			err := vm.executeCompare(op)
			if err != nil {
				return err
			}
		case code.OpNot:
			err := vm.executeNotOp()
			if err != nil {
				return err
			}
		case code.OpNeg:
			err := vm.executeNegOp()
			if err != nil {
				return err
			}
		case code.OpTrue:
			err := vm.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := vm.push(False)
			if err != nil {
				return err
			}
		case code.OpJump:
			pos := int(code.ReadUint16(inst[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case code.OpJNE:
			pos := int(code.ReadUint16(inst[ip+1:]))
			vm.currentFrame().ip += 2

			condition := vm.pop()
			if !isTrue(condition) {
				vm.currentFrame().ip = pos - 1
			}
		case code.OpNull:
			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip += 2

			vm.globals[globalIndex] = vm.pop()
		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(inst[ip+1:])
			vm.currentFrame().ip += 2

			err := vm.push(vm.globals[globalIndex])
			if err != nil {
				return err
			}
		case code.OpSetLocal:
			localIndex := code.ReadUint8(inst[ip+1:])
			vm.currentFrame().ip++

			frame := vm.currentFrame()
			vm.stack[frame.basePointer+int(localIndex)] = vm.pop()
		case code.OpGetLocal:
			localIndex := code.ReadUint8(inst[ip+1:])
			vm.currentFrame().ip++

			frame := vm.currentFrame()
			err := vm.push(vm.stack[frame.basePointer+int(localIndex)])
			if err != nil {
				return err
			}

		case code.OpArray:
			numElements := int(code.ReadUint16(inst[ip+1:]))
			vm.currentFrame().ip += 2

			array := vm.buildArrayObject(vm.sp-numElements, vm.sp)
			vm.sp = vm.sp - numElements

			err := vm.push(array)
			if err != nil {
				return err
			}
		case code.OpHashTable:
			numElements := int(code.ReadUint16(inst[ip+1:]))
			vm.currentFrame().ip += 2

			hashmap, err := vm.buildHashmapObject(vm.sp-numElements, vm.sp)
			if err != nil {
				return err
			}
			vm.sp = vm.sp - numElements

			err = vm.push(hashmap)
			if err != nil {
				return err
			}
		case code.OpIndex:
			index := vm.pop()
			left := vm.pop()

			err := vm.executeIndexExpression(left, index)
			if err != nil {
				return err
			}
		case code.OpCall:

			numArgs := code.ReadUint8(inst[ip+1:])
			vm.currentFrame().ip++
			err := vm.executeFunctionCall(int(numArgs))
			if err != nil {
				return err
			}
		case code.OpReturnValue:
			// pop the return value from stack
			returnVal := vm.pop()
			// pop the stack frame
			frame := vm.popFrame()
			// restore the stack pointer
			// we substract 1 to throw away the last value on the stack (function call)
			// without explicitly poping the value from the stack
			vm.sp = frame.basePointer - 1
			// push the actual return value
			err := vm.push(returnVal)
			if err != nil {
				return err
			}
		case code.OpReturn:
			frame := vm.popFrame()
			vm.sp = frame.basePointer - 1

			err := vm.push(Null)
			if err != nil {
				return err
			}
		case code.OpGetBuiltin:
			builtinIndex := code.ReadUint8(inst[ip+1:])
			vm.currentFrame().ip++

			def := object.Builtins[builtinIndex]

			err := vm.push(def.Fn)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// push an element to the stack and increment the stack pointer.
func (vm *VM) push(obj object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = obj
	vm.sp++

	return nil
}

// pop the top element of the stack
func (vm *VM) pop() object.Object {
	obj := vm.stack[vm.sp-1]
	vm.sp--

	return obj
}

// isTrue checks if the given object evaluates to boolean true
func isTrue(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

// LastPoppedStackElem returns the last pop'd item.
func (vm *VM) LastPoppedStackElem() object.Object {
	return vm.stack[vm.sp]
}

// executeFunctionCall executes a function using arguments
func (vm *VM) executeFunctionCall(numArgs int) error {

	callee := vm.stack[vm.sp-1-int(numArgs)]

	switch callee := callee.(type) {
	case *object.CompiledFunction:
		return vm.callFunc(callee, numArgs)
	case *object.BuiltIn:
		return vm.callBuiltIn(callee, numArgs)
	default:
		return fmt.Errorf("calling non-function")
	}
}

// callFunc executes a function call on user defined functions
func (vm *VM) callFunc(fn *object.CompiledFunction, numArgs int) error {
	if numArgs != fn.NumParams {
		return fmt.Errorf("wrong number of parameters : want %d, got %d", fn.NumParams, numArgs)
	}
	// substract numArgs to correctly set bp
	frame := NewFrame(fn, vm.sp-numArgs)
	vm.pushFrame(frame)
	// increment the stack pointer to make place for local variables
	vm.sp = frame.basePointer + fn.NumLocals
	return nil
}

// callBuiltIn executes a builtin function call
func (vm *VM) callBuiltIn(builtin *object.BuiltIn, numArgs int) error {

	args := vm.stack[vm.sp-numArgs : vm.sp]

	res := builtin.Fn(args...)
	vm.sp = vm.sp - numArgs - 1

	if res != nil {
		vm.push(res)
	} else {
		vm.push(Null)
	}

	return nil
}

// executeBinOp executes a binary opeartion
func (vm *VM) executeBinOp(op code.OpCode) error {
	right := vm.pop()
	left := vm.pop()

	rightType := right.Type()
	leftType := right.Type()

	if rightType == object.INTEGER && leftType == object.INTEGER {
		return vm.executeIntegerOp(op, left, right)
	} else if rightType == object.STRING && leftType == object.STRING {
		return vm.executeStringOp(op, left, right)
	}

	return fmt.Errorf("unsupported types for binary operation for %d %s %s", op, rightType, leftType)
}

// executeIntegerOp executes a binary operation on integers
func (vm *VM) executeIntegerOp(op code.OpCode, left, right object.Object) error {

	rightVal := right.(*object.Integer).Value
	leftVal := left.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftVal + rightVal
	case code.OpSub:
		result = leftVal - rightVal
	case code.OpMul:
		result = leftVal * rightVal
	case code.OpDiv:
		result = leftVal / rightVal
	case code.OpMod:
		result = leftVal % rightVal
	}
	return vm.push(&object.Integer{Value: result})
}

// executeStringOp executes a binary operation on integers
func (vm *VM) executeStringOp(op code.OpCode, left, right object.Object) error {

	rightVal := right.(*object.String).Value
	leftVal := left.(*object.String).Value

	if op != code.OpAdd {
		return fmt.Errorf("Unknown operator %d for type %s %s", op, left.Type(), right.Type())
	}
	result := leftVal + rightVal

	return vm.push(&object.String{Value: result})
}

// executeCompare executes comparison opcodes pushing the result to the stack
func (vm *VM) executeCompare(op code.OpCode) error {

	right := vm.pop()
	left := vm.pop()

	leftType := left.Type()
	rightType := right.Type()

	if rightType == object.INTEGER && leftType == object.INTEGER {
		return vm.executeIntegerCompare(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(left == right))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(left != right))
	default:
		return fmt.Errorf("unknown operator : %d for type %s %s", op, leftType, rightType)
	}
}

// executeIntegerCompare compares two integers using a comparison operator and pushes the result to the stack
func (vm *VM) executeIntegerCompare(op code.OpCode, left, right object.Object) error {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch op {

	case code.OpEqual:
		return vm.push(nativeBoolToBooleanObject(leftVal == rightVal))
	case code.OpNotEqual:
		return vm.push(nativeBoolToBooleanObject(leftVal != rightVal))
	case code.OpGreaterOrEqual:
		return vm.push(nativeBoolToBooleanObject(leftVal >= rightVal))
	case code.OpGreaterThan:
		return vm.push(nativeBoolToBooleanObject(leftVal > rightVal))
	default:
		return fmt.Errorf("unknown operator %d for INTEGER Type", op)
	}
}

// executeNotOp on the top item of the stack
func (vm *VM) executeNotOp() error {

	operand := vm.pop()

	switch operand {
	case True:
		return vm.push(False)
	case False:
		return vm.push(True)
	case Null:
		return vm.push(True)
	default:
		return vm.push(False)
	}
}

// executeNegOp on the top item of the stack (must be an INTEGER)
func (vm *VM) executeNegOp() error {
	operand := vm.pop()

	switch operand.Type() {
	case object.INTEGER:
		val := operand.(*object.Integer).Value
		return vm.push(&object.Integer{Value: -val})
	default:
		return fmt.Errorf("Unsupported operand type %s for NEG operator", operand.Type())
	}
}

// executeIndexExpression pops the index and the object to be indexed from the stack
// and pushes the value
func (vm *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY && index.Type() == object.INTEGER:
		return vm.executeArrayIndex(left, index)
	case left.Type() == object.HASH:
		return vm.executeHashMapIndex(left, index)
	default:
		return fmt.Errorf("index operator %T not supported for type %T", index.Type(), left.Type())
	}
}

// executeArrayIndex executes the indexing operations for array cases
func (vm *VM) executeArrayIndex(array, index object.Object) error {
	arrayObj := array.(*object.Array)
	idx := index.(*object.Integer).Value

	max := int64(len(arrayObj.Elements) - 1)

	if idx < 0 || idx > max {
		return vm.push(Null)
	}
	return vm.push(arrayObj.Elements[idx])
}

// executeHashmapIndex executes the indexing operations for hashmaps
func (vm *VM) executeHashMapIndex(hash, index object.Object) error {
	hashObj := hash.(*object.HashMap)

	key, ok := index.(object.Hashable)
	if !ok {
		return fmt.Errorf("invalid key for hashmap type :%T", index.Type())
	}
	pair, ok := hashObj.Pairs[key.HashKey()]
	if !ok {
		return vm.push(Null)
	}

	return vm.push(pair.Value)
}

// nativeBoolToBooleanObject returns a boolean object equivalent to input
func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}

	return False
}

// buildArrayObject pops n elements from the stack starting at index k
// and ending at index j and builds an array object of those.
func (vm *VM) buildArrayObject(startIndex, endIndex int) object.Object {
	elements := make([]object.Object, endIndex-startIndex)

	for i := startIndex; i < endIndex; i++ {
		elements[i-startIndex] = vm.stack[i]
	}

	return &object.Array{Elements: elements}
}

// buildHashmapObject creates a new object.Hashmap from stack elements
func (vm *VM) buildHashmapObject(startIndex, endIndex int) (object.Object, error) {

	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := startIndex; i < endIndex; i += 2 {
		key := vm.stack[i]
		val := vm.stack[i+1]

		pair := object.HashPair{Key: key, Value: val}
		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("Type unusable for hashmap key : %s", key.Type())
		}
		hashedPairs[hashKey.HashKey()] = pair
	}

	return &object.HashMap{Pairs: hashedPairs}, nil
}
