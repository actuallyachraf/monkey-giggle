package compiler

import (
	"fmt"
	"testing"

	"github.com/actuallyachraf/monkey-giggle/ast"
	"github.com/actuallyachraf/monkey-giggle/code"
	"github.com/actuallyachraf/monkey-giggle/lexer"
	"github.com/actuallyachraf/monkey-giggle/object"
	"github.com/actuallyachraf/monkey-giggle/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []interface{}
	expectedInstructions []code.Instructions
}

func TestCompiler(t *testing.T) {
	t.Run("TestIntegerArithmetic", func(t *testing.T) {

		tests := []compilerTestCase{
			{
				input:             "1 + 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 - 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSub),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 * 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpMul),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 / 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpDiv),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 % 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpMod),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestBooleanExpressions", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "true",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "false",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpFalse),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 > 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThan),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 >= 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterOrEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 < 2",
				expectedConstants: []interface{}{2, 1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterThan),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 <= 2",
				expectedConstants: []interface{}{2, 1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpGreaterOrEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 == 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "1 != 2",
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpNotEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "true != false",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpFalse),
					code.Make(code.OpNotEqual),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "true == true",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpTrue),
					code.Make(code.OpEqual),
					code.Make(code.OpPop),
				},
			},
		}

		runCompilerTests(t, tests)
	})
	t.Run("TestPrefixExpression", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "-1",
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpNeg),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "!true",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpNot),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestConditionalExpression", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             `if (true) {10}; 3333;`,
				expectedConstants: []interface{}{10, 3333},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpJNE, 10),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpJump, 11),
					code.Make(code.OpNull),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPop),
				},
			}, {
				input:             `if (true) {10} else {20}; 3333;`,
				expectedConstants: []interface{}{10, 20, 3333},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpTrue),
					code.Make(code.OpJNE, 10),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpJump, 13),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpPop),
				},
			}, {
				input:             `if (false) {10} else {20}; 3333;`,
				expectedConstants: []interface{}{10, 20, 3333},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpFalse),
					code.Make(code.OpJNE, 10),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpJump, 13),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpPop),
				},
			}, {
				input:             `if (false) {10} ; 3333;`,
				expectedConstants: []interface{}{10, 3333},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpFalse),
					code.Make(code.OpJNE, 10),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpJump, 11),
					code.Make(code.OpNull),
					code.Make(code.OpPop),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestGlobalLetStatement", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: `
				let one = 1;
				let two = 2;
				`,
				expectedConstants: []interface{}{1, 2},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpSetGlobal, 1),
				},
			}, {
				input: `
				let one = 1;
				one;
				`,
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpPop),
				},
			}, {
				input: `
				let one = 1;
				let two = one;
				two;`,
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpSetGlobal, 1),
					code.Make(code.OpGetGlobal, 1),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestStringExpression", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             `"monkey"`,
				expectedConstants: []interface{}{"monkey"},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpPop),
				},
			}, {
				input:             `"mon" + "key"`,
				expectedConstants: []interface{}{"mon", "key"},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestArrayLiterals", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "[]",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpArray, 0),
					code.Make(code.OpPop),
				},
			}, {
				input:             "[1,2,3]",
				expectedConstants: []interface{}{1, 2, 3},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpArray, 3),
					code.Make(code.OpPop),
				},
			}, {
				input:             "[1+2,3-4,5*6]",
				expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpAdd),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpSub),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpConstant, 5),
					code.Make(code.OpMul),
					code.Make(code.OpArray, 3),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestHashmapLiteral", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "{}",
				expectedConstants: []interface{}{},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpHashTable, 0),
					code.Make(code.OpPop),
				},
			}, {
				input:             "{1:2,3:4,5:6}",
				expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpConstant, 5),
					code.Make(code.OpHashTable, 6),
					code.Make(code.OpPop),
				},
			}, {
				input:             "{1:2+3,4:5*6}",
				expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpAdd),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpConstant, 5),
					code.Make(code.OpMul),
					code.Make(code.OpHashTable, 4),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestIndexExpression", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input:             "[1,2,3][1+1]",
				expectedConstants: []interface{}{1, 2, 3, 1, 1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpArray, 3),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpAdd),
					code.Make(code.OpIndex),
					code.Make(code.OpPop),
				},
			},
			{
				input:             "{1:2,3:4}[1+1]",
				expectedConstants: []interface{}{1, 2, 3, 4, 1, 1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpHashTable, 4),
					code.Make(code.OpConstant, 4),
					code.Make(code.OpConstant, 5),
					code.Make(code.OpAdd),
					code.Make(code.OpIndex),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestFunctionLiteral", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: "fn() { return 5 + 10; }",
				expectedConstants: []interface{}{
					5,
					10,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpPop),
				},
			}, {
				input: `fn() { 5 + 10 }`,
				expectedConstants: []interface{}{
					5,
					10,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpPop),
				},
			}, {
				input: `fn() {  1;2 }`,
				expectedConstants: []interface{}{
					1,
					2,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpPop),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpPop),
				},
			}, {
				input: `fn() {}`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpReturn),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestFunctionCall", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: `fn(){24}();`,
				expectedConstants: []interface{}{
					24,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpCall, 0),
					code.Make(code.OpPop),
				},
			}, {
				input: `
				let noArg = fn(){24};
				noArg();
				`,
				expectedConstants: []interface{}{
					24,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpCall, 0),
					code.Make(code.OpPop),
				},
			}, {
				input: `
				let oneArg = fn(a) {};
				oneArg(24);
				`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpReturn),
					},
					24,
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpCall, 1),
					code.Make(code.OpPop),
				},
			}, {
				input: `
				let manyArg = fn(a,b,c){};
				manyArg(24,25,26);
				`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpReturn),
					},
					24,
					25,
					26,
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpCall, 3),
					code.Make(code.OpPop),
				},
			}, {
				input: `
				let oneArg = fn(a) {a};
				oneArg(24);
				`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpReturnValue),
					},
					24,
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpCall, 1),
					code.Make(code.OpPop),
				},
			}, {
				input: `
				let manyArg = fn(a,b,c){a;b;c};
				manyArg(24,25,26);
				`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpPop),
						code.Make(code.OpGetLocal, 1),
						code.Make(code.OpPop),
						code.Make(code.OpGetLocal, 2),
						code.Make(code.OpReturnValue),
					},
					24,
					25,
					26,
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpGetGlobal, 0),
					code.Make(code.OpConstant, 1),
					code.Make(code.OpConstant, 2),
					code.Make(code.OpConstant, 3),
					code.Make(code.OpCall, 3),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestLetStatementWithScope", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: `
				let num = 55;
				fn() { num }
				`,
				expectedConstants: []interface{}{
					55,
					[]code.Instructions{
						code.Make(code.OpGetGlobal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpConstant, 0),
					code.Make(code.OpSetGlobal, 0),
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpPop),
				},
			}, {
				input: `
				fn(){
					let num = 55;
					num
				}

				`,
				expectedConstants: []interface{}{
					55,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpPop),
				},
			}, {
				input: `
				fn(){
					let a = 55;
					let b = 75;
					a+b
				}

				`,
				expectedConstants: []interface{}{
					55,
					75,
					[]code.Instructions{
						code.Make(code.OpConstant, 0),
						code.Make(code.OpSetLocal, 0),
						code.Make(code.OpConstant, 1),
						code.Make(code.OpSetLocal, 1),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpGetLocal, 1),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 2, 0),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestBuiltinFunc", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: `
				len([]);
				append([],1);
				`,
				expectedConstants: []interface{}{1},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpGetBuiltin, 0),
					code.Make(code.OpArray, 0),
					code.Make(code.OpCall, 1),
					code.Make(code.OpPop),
					code.Make(code.OpGetBuiltin, 4),
					code.Make(code.OpArray, 0),
					code.Make(code.OpConstant, 0),
					code.Make(code.OpCall, 2),
					code.Make(code.OpPop),
				},
			}, {
				input: `
				fn(){len([])}
				`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpGetBuiltin, 0),
						code.Make(code.OpArray, 0),
						code.Make(code.OpCall, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 0, 0),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
	t.Run("TestClosures", func(t *testing.T) {
		tests := []compilerTestCase{
			{
				input: `
			fn(a){
				fn(b){
					a+b
				}
			}
			`,
				expectedConstants: []interface{}{
					[]code.Instructions{
						code.Make(code.OpGetFree, 0),
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpAdd),
						code.Make(code.OpReturnValue),
					},
					[]code.Instructions{
						code.Make(code.OpGetLocal, 0),
						code.Make(code.OpClosure, 0, 1),
						code.Make(code.OpReturnValue),
					},
				},
				expectedInstructions: []code.Instructions{
					code.Make(code.OpClosure, 1, 0),
					code.Make(code.OpPop),
				},
			},
		}
		runCompilerTests(t, tests)
	})
}
func TestSymbolTable(t *testing.T) {
	t.Run("TestDefine", func(t *testing.T) {
		expected := map[string]Symbol{
			"a": {Name: "a", Scope: GlobalScope, Index: 0},
			"b": {Name: "b", Scope: GlobalScope, Index: 1},
			"c": {Name: "c", Scope: LocalScope, Index: 0},
			"d": {Name: "d", Scope: LocalScope, Index: 1},
			"e": {Name: "e", Scope: LocalScope, Index: 0},
			"f": {Name: "f", Scope: LocalScope, Index: 1},
		}

		global := NewSymbolTable()

		a := global.Define("a")
		if a != expected["a"] {
			t.Errorf("expected a=%+v, got=%+v", expected["a"], a)
		}

		b := global.Define("b")
		if b != expected["b"] {
			t.Errorf("expected b=%+v, got=%+v", expected["b"], b)
		}

		firstLocal := NewEnclosedSymbolTable(global)

		c := firstLocal.Define("c")
		if c != expected["c"] {
			t.Errorf("expected c=%+v, got=%+v", expected["c"], c)
		}

		d := firstLocal.Define("d")
		if d != expected["d"] {
			t.Errorf("expected d=%+v, got=%+v", expected["d"], d)
		}

		secondLocal := NewEnclosedSymbolTable(firstLocal)

		e := secondLocal.Define("e")
		if e != expected["e"] {
			t.Errorf("expected e=%+v, got=%+v", expected["e"], e)
		}

		f := secondLocal.Define("f")
		if f != expected["f"] {
			t.Errorf("expected f=%+v, got=%+v", expected["f"], f)
		}
	})
	t.Run("TestGlobalResolve", func(t *testing.T) {

		globals := NewSymbolTable()
		globals.Define("a")
		globals.Define("b")

		expected := map[string]Symbol{
			"a": {Name: "a", Scope: GlobalScope, Index: 0},
			"b": {Name: "b", Scope: GlobalScope, Index: 1},
		}
		for _, sym := range expected {
			result, ok := globals.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v instead got %+v", sym.Name, sym, result)
			}
		}
	})
	t.Run("TestLocalResolve", func(t *testing.T) {
		globals := NewSymbolTable()
		globals.Define("a")
		globals.Define("b")

		locals := NewEnclosedSymbolTable(globals)
		locals.Define("c")
		locals.Define("d")

		expected := map[string]Symbol{
			"a": {Name: "a", Scope: GlobalScope, Index: 0},
			"b": {Name: "b", Scope: GlobalScope, Index: 1},
			"c": {Name: "c", Scope: LocalScope, Index: 0},
			"d": {Name: "d", Scope: LocalScope, Index: 1},
		}
		for _, sym := range expected {
			result, ok := locals.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v instead got %+v", sym.Name, sym, result)
			}
		}
	})
	t.Run("TestResolveNested", func(t *testing.T) {
		global := NewSymbolTable()
		global.Define("a")
		global.Define("b")

		firstLocal := NewEnclosedSymbolTable(global)
		firstLocal.Define("c")
		firstLocal.Define("d")

		secondLocal := NewEnclosedSymbolTable(firstLocal)
		secondLocal.Define("e")
		secondLocal.Define("f")

		tests := []struct {
			table           *SymbolTable
			expectedSymbols []Symbol
		}{
			{
				firstLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0},
					{Name: "b", Scope: GlobalScope, Index: 1},
					{Name: "c", Scope: LocalScope, Index: 0},
					{Name: "d", Scope: LocalScope, Index: 1},
				},
			},
			{
				secondLocal,
				[]Symbol{
					{Name: "a", Scope: GlobalScope, Index: 0},
					{Name: "b", Scope: GlobalScope, Index: 1},
					{Name: "e", Scope: LocalScope, Index: 0},
					{Name: "f", Scope: LocalScope, Index: 1},
				},
			},
		}

		for _, tt := range tests {
			for _, sym := range tt.expectedSymbols {
				result, ok := tt.table.Resolve(sym.Name)
				if !ok {
					t.Errorf("name %s not resolvable", sym.Name)
					continue
				}
				if result != sym {
					t.Errorf("expected %s to resolve to %+v, got=%+v",
						sym.Name, sym, result)
				}
			}
		}
	})
}
func TestCompilerScope(t *testing.T) {
	compiler := New()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 0)
	}
	globalSymbolTable := compiler.symbolTable

	compiler.emit(code.OpMul)

	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d", compiler.scopeIndex, 1)
	}

	compiler.emit(code.OpSub)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 1 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last := compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpSub {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
			last.Opcode, code.OpSub)
	}

	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Errorf("compiler did not enclose symbolTable")
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex wrong. got=%d, want=%d",
			compiler.scopeIndex, 0)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not restore global symbol table")
	}
	if compiler.symbolTable.Outer != nil {
		t.Errorf("compiler modified global symbol table incorrectly")
	}

	compiler.emit(code.OpAdd)

	if len(compiler.scopes[compiler.scopeIndex].instructions) != 2 {
		t.Errorf("instructions length wrong. got=%d",
			len(compiler.scopes[compiler.scopeIndex].instructions))
	}

	last = compiler.scopes[compiler.scopeIndex].lastInstruction
	if last.Opcode != code.OpAdd {
		t.Errorf("lastInstruction.Opcode wrong. got=%d, want=%d",
			last.Opcode, code.OpAdd)
	}

	previous := compiler.scopes[compiler.scopeIndex].previousInstruction
	if previous.Opcode != code.OpMul {
		t.Errorf("previousInstruction.Opcode wrong. got=%d, want=%d",
			previous.Opcode, code.OpMul)
	}

}
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)

	program := p.Parse()

	return program
}
func runCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		compiler := New()
		program := parse(tt.input)
		err := compiler.Compile(program)

		if err != nil {
			t.Fatalf("compiler error : %s", err)
		}

		bytecode := compiler.Bytecode()

		err = testInstructions(tt.expectedInstructions, bytecode.Instructions)

		if err != nil {
			t.Fatalf("testInstructions failed : %s -- expected %q got %q", err, tt.expectedInstructions, bytecode.Instructions)
		}

		err = testConstants(t, tt.expectedConstants, bytecode.Constants)

		if err != nil {
			t.Fatalf("testConstants failed : %s", err)
		}
	}
}

func testInstructions(expected []code.Instructions, actual code.Instructions) error {
	concatted := concatInstructions(expected)
	if len(actual) != len(concatted) {
		return fmt.Errorf("wrong instructions length expected %d %q got %d %q", len(concatted), concatted, len(actual), actual)
	}

	for i, inst := range concatted {
		if actual[i] != inst {
			return fmt.Errorf("wrong instruction at offset %d expected %q got %q", i, concatted, actual)
		}
	}
	return nil
}

func concatInstructions(s []code.Instructions) code.Instructions {
	out := code.Instructions{}

	for _, inst := range s {
		out = append(out, inst...)
	}

	return out
}

func testConstants(t *testing.T, expected []interface{}, actual []object.Object) error {

	if len(expected) != len(actual) {
		return fmt.Errorf("wrong instructions length expected %d got %d", len(expected), len(actual))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			err := testIntegerObject(int64(constant), actual[i])
			if err != nil {
				return fmt.Errorf("constant[%d] = %d - testIntegerObject failed with error : %s", i, constant, err)
			}
		case string:
			err := testStringObject(constant, actual[i])
			if err != nil {
				return fmt.Errorf("constant[%d] = %s - testStringObject failed with error : %s", i, constant, err)
			}
		case []code.Instructions:
			fn, ok := actual[i].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("constant[%d] not a function: %T", i, actual[i])
			}
			err := testInstructions(constant, fn.Instructions)
			if err != nil {
				return fmt.Errorf("constant[%d] - testInstruction failed with error %s", i, err)
			}
		}
	}
	return nil
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer got %T (%+v) ", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value expected %d got %d", expected, result.Value)
	}
	return nil
}
func testStringObject(expected string, actual object.Object) error {
	result, ok := actual.(*object.String)
	if !ok {
		return fmt.Errorf("object is not Integer got %T (%+v) ", actual, actual)
	}
	if result.Value != expected {
		return fmt.Errorf("object has wrong value expected %s got %s", expected, result.Value)
	}
	return nil
}
