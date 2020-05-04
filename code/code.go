// Package code implements Opcode semantics for the virtual machine bytecode.
package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// OpCode is a single byte constant specifying an instruction.
type OpCode byte

// Instructions are sequences of opcodes.
type Instructions []byte

const (
	// OpConstant is a fetch and push opcode for constants
	OpConstant OpCode = iota
	// OpAdd pops two integers from the stack and pushes their sum
	OpAdd
	// OpSub pops two integers from the stack and pushes their difference
	OpSub
	// OpMul pops two integers from the stack and pushes their product
	OpMul
	// OpDiv pops two integers from the stack and pushes their quotient
	OpDiv
	// OpMod pops two integers from the stack and pushes their remainder
	OpMod
	// OpNot pops a boolean from the stack and pushes it's opposite
	OpNot
	// OpNeg pops an integer from the stack and pushes it's negation
	OpNeg
	// OpPop pops an element from the stack
	OpPop
	// OpTrue pushes a truth value
	OpTrue
	// OpFalse pushes a false value
	OpFalse
	// OpEqual pops two elements from the stack pushes true if equal false otherwise.
	OpEqual
	// OpNotEqual pops two elements from the stack pushes true if not equal false otherwise.
	OpNotEqual
	// OpGreaterThan pops two elements from stack push true if a > b
	OpGreaterThan
	// OpGreaterOrEqual pops two elements from the stack pushes true if a >= b
	OpGreaterOrEqual
	// OpJNE implements Jump if not equal
	OpJNE
	// OpJump implements jump to address
	OpJump
	// OpNull pushes the null value to the stack
	OpNull
	// OpGetGlobal fetches binding values from the globals state
	OpGetGlobal
	// OpSetGlobal binds a value to an identifier
	OpSetGlobal
	// OpArray constructs an array object by popping N objects from the stack
	OpArray
	// OpHashTable constructs an associative array of objects
	OpHashTable
	// OpIndex reads the object and index and pushes the indexed element to the stack
	OpIndex
	// OpCall executes function calls
	OpCall
	// OpReturnValue explicitly pushes function return values to the stack
	OpReturnValue
	// OpReturn returns from the function call to the caller
	OpReturn
	// OpGetLocal is used to get local bindings
	OpGetLocal
	// OpSetLocal is used to set local bindings
	OpSetLocal
	// OpGetBuiltin is used to fetch built-in function from their scope
	OpGetBuiltin
)

// Definition represents information about opcodes.
type Definition struct {
	Name          string
	OperandWidths []int
}

var lookupTable = map[OpCode]Definition{
	OpConstant:       {"OpConstant", []int{2}},
	OpAdd:            {"OpAdd", []int{}},
	OpSub:            {"OpSub", []int{}},
	OpMul:            {"OpMul", []int{}},
	OpDiv:            {"OpDiv", []int{}},
	OpMod:            {"OpMod", []int{}},
	OpPop:            {"OpPop", []int{}},
	OpTrue:           {"OpTrue", []int{}},
	OpFalse:          {"OpFalse", []int{}},
	OpEqual:          {"OpEqual", []int{}},
	OpNotEqual:       {"OpNotEqual", []int{}},
	OpGreaterThan:    {"OpGreaterThan", []int{}},
	OpGreaterOrEqual: {"OpGreaterThanOrEqual", []int{}},
	OpNeg:            {"OpNeg", []int{}},
	OpNot:            {"OpNot", []int{}},
	OpJNE:            {"OpJumpIfNotEqual", []int{2}},
	OpJump:           {"OpJump", []int{2}},
	OpNull:           {"OpNull", []int{}},
	OpGetGlobal:      {"OpGetGlobal", []int{2}},
	OpSetGlobal:      {"OpSetGlobal", []int{2}},
	OpArray:          {"OpArray", []int{2}},
	OpHashTable:      {"OpHashTable", []int{2}},
	OpIndex:          {"OpIndex", []int{}},
	OpCall:           {"OpCall", []int{1}},
	OpReturnValue:    {"OpReturnValue", []int{}},
	OpReturn:         {"OpReturn", []int{}},
	OpGetLocal:       {"OpGetLocal", []int{1}},
	OpSetLocal:       {"OpSetLocal", []int{1}},
	OpGetBuiltin:     {"OpGetBuiltin", []int{1}},
}

// Lookup fetches the opcode definition.
func Lookup(op OpCode) (Definition, error) {
	def, ok := lookupTable[op]
	if !ok {
		return Definition{}, fmt.Errorf("Opcode %d is undefined", op)
	}

	return def, nil
}

// Make creates an instruction sequence given an opcode and operands.
func Make(op OpCode, operands ...int) []byte {
	def, ok := lookupTable[op]
	if !ok {
		return []byte{}
	}

	instLen := 1

	for _, w := range def.OperandWidths {
		instLen += w
	}

	inst := make([]byte, instLen)
	inst[0] = byte(op)

	offset := 1

	for i, operand := range operands {
		width := def.OperandWidths[i]

		switch width {
		case 2:
			binary.BigEndian.PutUint16(inst[offset:], uint16(operand))
		case 1:
			inst[offset] = byte(operand)
		}
		offset += width
	}

	return inst

}

// FormatInstruction returns a pretty printed instruction
func (inst Instructions) FormatInstruction(def Definition, operands []int) string {
	operandCount := len(def.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR : operand length %d does not match defined %d", len(operands), operandCount)
	}
	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operand count for %s", def.Name)
}

// String implements the stringer interface
func (inst Instructions) String() string {

	var out bytes.Buffer

	i := 0
	for i < len(inst) {
		def, err := Lookup(OpCode(inst[i]))
		if err != nil {
			fmt.Fprintf(&out, "ERROR : %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, inst[i+1:])

		fmt.Fprintf(&out, "%04d %s\n", i, inst.FormatInstruction(def, operands))

		i += 1 + read

	}

	return out.String()
}

// ReadUint16 reads a big-endian encoded uint16 from an instruction slice
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

// ReadOperands parses definition operand width
func ReadOperands(def Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))

	offset := 0

	for i, width := range def.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUint8(ins[offset:]))
		}

		offset += width
	}

	return operands, offset
}

// ReadUint8 reads a single byte integer
func ReadUint8(ins Instructions) uint8 {
	return uint8(ins[0])
}
