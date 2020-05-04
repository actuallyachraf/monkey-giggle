package vm

import (
	"github.com/actuallyachraf/monkey-giggle/code"
	"github.com/actuallyachraf/monkey-giggle/object"
)

// Frame represents a stack frame used to execute function calls
type Frame struct {
	fn          *object.CompiledFunction
	ip          int
	basePointer int
}

// NewFrame creates a new frame for a given compiled function
func NewFrame(fn *object.CompiledFunction, bp int) *Frame {
	return &Frame{
		fn:          fn,
		ip:          -1,
		basePointer: bp,
	}
}

// Instructions returns the set of instructions inside the call frame
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
