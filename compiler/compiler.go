// Package compiler defines the core compiler code.
package compiler

import (
	"fmt"
	"sort"

	"github.com/actuallyachraf/monkey-giggle/ast"
	"github.com/actuallyachraf/monkey-giggle/code"
	"github.com/actuallyachraf/monkey-giggle/object"
)

// Compiler represents structs and defined constant objects
type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable

	scopes     []CompilationScope
	scopeIndex int
}

// CompilationScope represents scopes for functions
type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

// EmittedInstruction represents an emitted compiler instruction
type EmittedInstruction struct {
	Opcode   code.OpCode
	Position int
}

// New creates a new compiler instance
func New() *Compiler {

	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	symTable := NewSymbolTable()

	for i, v := range object.Builtins {
		symTable.DefineBuiltIn(i, v.Name)
	}
	return &Compiler{
		constants:   []object.Object{},
		symbolTable: symTable,
		scopes:      []CompilationScope{mainScope},
		scopeIndex:  0,
	}
}

// currentInstructions returns the instructions within the current scope
func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

// NewWithState creates a new compiler instance with a predefined symbol table
func NewWithState(symTable *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = symTable
	compiler.constants = constants

	return compiler
}

// Compile takes an AST Node and returns an equivalent compiler.
func (c *Compiler) Compile(node ast.Node) error {

	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)
	case *ast.InfixExpression:
		if node.Operator == "<" || node.Operator == "<=" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}
			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			if node.Operator == "<" {
				c.emit(code.OpGreaterThan)
			} else if node.Operator == "<=" {
				c.emit(code.OpGreaterOrEqual)
			}
			return nil
		}
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case "%":
			c.emit(code.OpMod)
		case ">":
			c.emit(code.OpGreaterThan)
		case ">=":
			c.emit(code.OpGreaterOrEqual)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("Unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	case *ast.BooleanLiteral:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpNot)
		case "-":
			c.emit(code.OpNeg)
		default:
			return fmt.Errorf("Unknown operator %s", node.Operator)
		}
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}
		JNEPos := c.emit(code.OpJNE, 9999)
		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}
		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}
		JMPPos := c.emit(code.OpJump, 9999)
		afterConsqPos := len(c.currentInstructions())
		c.changeOperand(JNEPos, afterConsqPos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err := c.Compile(node.Alternative)
			if err != nil {
				return err
			}
			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}

		}
		afterAltPos := len(c.currentInstructions())
		c.changeOperand(JMPPos, afterAltPos)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.LetStatement:
		sym := c.symbolTable.Define(string(node.Name.Value))
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}
		if sym.Scope == GlobalScope {
			c.emit(code.OpSetGlobal, sym.Index)
		} else {
			c.emit(code.OpSetLocal, sym.Index)
		}
	case *ast.Identifier:
		sym, ok := c.symbolTable.Resolve(string(node.Value))
		if !ok {
			return fmt.Errorf("Undefined variable %s", node.Value)
		}
		c.loadSymbol(sym)
	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(str))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := c.Compile(el)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *ast.HashmapLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := c.Compile(k)
			if err != nil {
				return err
			}
			err = c.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}
		c.emit(code.OpHashTable, len(node.Pairs)*2)
	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Index)
		if err != nil {
			return err
		}
		c.emit(code.OpIndex)
	case *ast.FunctionLiteral:
		c.enterScope()

		for _, p := range node.Parameters {
			c.symbolTable.Define(string(p.Value))
		}

		err := c.Compile(node.Body)
		if err != nil {
			return err
		}
		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithRet()
		}

		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		numLocals := c.symbolTable.numDefinitions
		inst := c.leaveScope()

		compiledFn := &object.CompiledFunction{Instructions: inst, NumLocals: numLocals, NumParams: len(node.Parameters)}
		c.emit(code.OpConstant, c.addConstant(compiledFn))
	case *ast.ReturnStatement:
		err := c.Compile(node.ReturnValue)
		if err != nil {
			return err
		}
		c.emit(code.OpReturnValue)
	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}
		for _, arg := range node.Arguments {
			err := c.Compile(arg)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpCall, len(node.Arguments))
	}
	return nil
}

// Bytecode represents a sequence of instructions and object table.
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

// Bytecode returns the generated bytecode.
func (c *Compiler) Bytecode() Bytecode {

	return Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

// addConstant adds a constant to the constant pool
func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

// emit generates an instruction and add it to the bytecode
func (c *Compiler) emit(op code.OpCode, operands ...int) int {
	inst := code.Make(op, operands...)
	pos := c.addInstruction(inst)

	c.setLastEmittedInstruction(op, pos)

	return pos
}

// lastInstructionIsPop checks if the last emitted instruction is OpPop
func (c *Compiler) lastInstructionIs(op code.OpCode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}
	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

// removeLastPop clears the last emitted pop instruction
func (c *Compiler) removeLastPop() {

	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	newInst := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = newInst
	c.scopes[c.scopeIndex].lastInstruction = previous
}

// setLastEmittedInstruction populates the last instruction from the current
// one before we emit a new one.
func (c *Compiler) setLastEmittedInstruction(op code.OpCode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

// replaceInstruction at position with a new one
func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {

	inst := c.currentInstructions()
	for i := 0; i < len(newInstruction); i++ {
		inst[pos+i] = newInstruction[i]
	}
}

// changeOperand replaces an opcode operand's
func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.OpCode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

// addInstruction appends a new instruction to the generated bytecode
func (c *Compiler) addInstruction(inst []byte) int {
	posNewInst := len(c.currentInstructions())
	updatedInst := append(c.currentInstructions(), inst...)

	c.scopes[c.scopeIndex].instructions = updatedInst

	return posNewInst
}

// enterScope creates a new compiler scope and makes it the current working scope
func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

// leaveScope dismisses the current working scope and decrement scope index
func (c *Compiler) leaveScope() code.Instructions {
	inst := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer
	return inst
}

// replaceLastPopWithRet replaces the last opPop with opRet
func (c *Compiler) replaceLastPopWithRet() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

// loadSymbol emits the proper symbol opcode
func (c *Compiler) loadSymbol(sym Symbol) {
	switch sym.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, sym.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, sym.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, sym.Index)
	}
}
