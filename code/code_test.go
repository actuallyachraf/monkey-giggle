package code

import (
	"testing"
)

func TestOpcode(t *testing.T) {
	t.Run("TestMake", func(t *testing.T) {

		tests := []struct {
			op       OpCode
			operands []int
			expected []byte
		}{
			{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
			{OpAdd, []int{}, []byte{byte(OpAdd)}},
			{OpClosure, []int{65534, 255}, []byte{byte(OpClosure), 255, 254, 255}},
		}

		for _, tt := range tests {
			inst := Make(tt.op, tt.operands...)
			if len(inst) != len(tt.expected) {
				t.Errorf("instruction has wrong length expected %d got %d", len(tt.expected), len(inst))
			}
			for i, b := range tt.expected {
				if inst[i] != b {
					t.Errorf("instruction has wrong byte at pos[%d] expected %d got %d", i, b, inst[i])
				}
			}
		}
	})
	t.Run("TestDisassemble", func(t *testing.T) {
		instructions := []Instructions{
			Make(OpAdd),
			Make(OpConstant, 2),
			Make(OpConstant, 65535),
		}

		expected := `0000 OpAdd
0001 OpConstant 2
0004 OpConstant 65535
`

		concatted := Instructions{}
		for _, inst := range instructions {
			concatted = append(concatted, inst...)
		}
		if expected != concatted.String() {
			t.Errorf("instructions wrongly formatted want %s got %s", expected, concatted.String())
		}

	})
	t.Run("TestReadOperand", func(t *testing.T) {
		tests := []struct {
			op        OpCode
			operands  []int
			bytesRead int
		}{
			{
				OpConstant,
				[]int{65535},
				2,
			},
		}

		for _, tt := range tests {
			inst := Make(tt.op, tt.operands...)

			def, err := Lookup(tt.op)
			if err != nil {
				t.Errorf("Lookup failed with error :%s", err)
			}
			operandsRead, n := ReadOperands(def, inst[1:])
			if n != tt.bytesRead {
				t.Errorf("wrong count of read bytes expected %d got %d", tt.bytesRead, n)
			}

			for i, operand := range tt.operands {
				if operandsRead[i] != operand {
					t.Errorf("wrong operand expected %d got %d", operand, operandsRead[i])
				}
			}

		}
	})
}
