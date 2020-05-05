package vm

import (
	"github.com/actuallyachraf/monkey-giggle/code"
	"github.com/actuallyachraf/monkey-giggle/object"
)

// Frame represents a stack frame used to execute function calls
type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

// NewFrame creates a new frame for a given compiled function
func NewFrame(cl *object.Closure, bp int) *Frame {
	return &Frame{
		cl:          cl,
		ip:          -1,
		basePointer: bp,
	}
}

// Instructions returns the set of instructions inside the call frame
func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
