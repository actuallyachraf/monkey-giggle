package vm

import (
	"fmt"
	"testing"

	"github.com/actuallyachraf/monkey-giggle/ast"
	"github.com/actuallyachraf/monkey-giggle/compiler"
	"github.com/actuallyachraf/monkey-giggle/lexer"
	"github.com/actuallyachraf/monkey-giggle/object"
	"github.com/actuallyachraf/monkey-giggle/parser"
)

func TestVM(t *testing.T) {
	t.Run("TestIntegerArithmetic", func(t *testing.T) {
		tests := []vmTestCase{
			{
				"1", 1,
			},
			{
				"2", 2,
			},
			{
				"1 + 2", 3,
			},
			{
				"50 / 2 * 2 + 10 -5", 55,
			},
			{
				"5 * (2 + 10)", 60,
			},
			{
				"350 % 3", 2,
			},
		}
		runVMTests(t, tests)
	})
	t.Run("TestBooleanExpression", func(t *testing.T) {
		tests := []vmTestCase{
			{"true", true},
			{"false", false},
			{"1 < 2", true},
			{"1 > 2", false},
			{"1 >= 1", true},
			{"2 <= 2", true},
			{"1 == 1", true},
			{"1 == 2", false},
			{"1 != 2", true},
			{"1 != 1", false},
			{"true == true", true},
			{"false == false", true},
			{"true == false", false},
			{"true != false", true},
			{"false != true", true},
			{"false != false", false},
			{"(1 < 2) == true", true},
			{"(1 > 2) == false", true},
		}
		runVMTests(t, tests)
	})
	t.Run("TestPrefixExpression", func(t *testing.T) {
		tests := []vmTestCase{
			{"!true", false},
			{"!false", true},
			{"-5", -5},
			{"-6", -6},
		}
		runVMTests(t, tests)
	})
	t.Run("TestConditionalExpression", func(t *testing.T) {
		tests := []vmTestCase{
			{"if (true) {10}", 10},
			{"if (true) {10} else {20}", 10},
			{"if (false) {10} else {20}", 20},
			{"if (1) {10}", 10},
			{"if (1 < 2) {10} else {20}", 10},
			{"if (1 > 2) {10} else {20}", 20},
			{"if (false) {10}", Null},
		}
		runVMTests(t, tests)
	})
	t.Run("TestGlobalLetStatement", func(t *testing.T) {
		tests := []vmTestCase{
			{"let one = 1;let two = 2;one + two;", 3},
			{"let one = 1;let two = 2; let three = one + two; three;", 3},
			{"let one = 1;one;", 1},
			{"let one = 1;let two = one + one; one + two;", 3},
		}
		runVMTests(t, tests)
	})
	t.Run("TestStringExpression", func(t *testing.T) {
		tests := []vmTestCase{
			{`"foobar"`, "foobar"},
			{`"foo"+"bar"`, "foobar"},
			{`"foo"+"bar"+"banana"`, "foobarbanana"},
			{`"Hello " + " World !"`, "Hello  World !"},
		}
		runVMTests(t, tests)
	})
	t.Run("TestArrayLiterals", func(t *testing.T) {
		tests := []vmTestCase{
			{"[]", []int{}},
			{"[1,2,3]", []int{1, 2, 3}},
			{"[1+2,3-4,5*6]", []int{3, -1, 30}},
			{`["foobar",3*4,"Hello"+"World"]`, []interface{}{"foobar", 12, "HelloWorld"}},
		}
		runVMTests(t, tests)
	})
	t.Run("TestHashmapLiterals", func(t *testing.T) {
		tests := []vmTestCase{
			{
				"{}", map[object.HashKey]int64{},
			}, {
				"{1:2,3:4,5:6}",
				map[object.HashKey]int64{
					(&object.Integer{Value: 1}).HashKey(): 2,
					(&object.Integer{Value: 3}).HashKey(): 4,
					(&object.Integer{Value: 5}).HashKey(): 6,
				},
			},
			{
				"{1 + 1:2*2,3+3:4*4}",
				map[object.HashKey]int64{
					(&object.Integer{Value: 2}).HashKey(): 4,
					(&object.Integer{Value: 6}).HashKey(): 16,
				},
			},
		}
		runVMTests(t, tests)
	})
	t.Run("TestIndexExpression", func(t *testing.T) {
		tests := []vmTestCase{
			{"[1,2,3][1]", 2},
			{"[1,2,3][0+2]", 3},
			{"[[1,2,3],[4,5,6]][1][1]", 5},
			{"[[1,2,3]][0][0]", 1},
			{"[1,2][-1]", Null},
			{"[][0]", Null},
			{"[1,2,3][99]", Null},
			{"{1:1,2:2}[1]", 1},
			{"{}[0]", Null},
			{"{1:1}[0]", Null},
			{"{1:1,2:2}[2]", 2},
		}
		runVMTests(t, tests)
	})
	t.Run("TestFunctionCallNoArgs", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let fivePlusTen = fn(){ 5 + 10;};
			fivePlusTen()
			`,
				expected: 15,
			},
			{
				input: `

				let one = fn(){1 ;};
				let two = fn(){2;};
				one() + two()
				`,
				expected: 3,
			}, {
				input: `
				let returnBefore = fn(){ return 99;100;}
				returnBefore()
				`,
				expected: 99,
			}, {
				input: `
				let earlyRet = fn(){ return 99; return 100;}
				earlyRet()
				`,
				expected: 99,
			}, {
				input: `
				let noRet = fn(){};
				noRet();
				`,
				expected: Null,
			}, {
				input: `
				let noRet = fn(){ };
				let noRetBis = fn(){ noRet(); };
				noRet();
				noRetBis();
				`,
				expected: Null,
			}, {
				input: `
				let returnsOne = fn(){1;};
				let returnsOnRet = fn() { returnsOne;};
				returnsOnRet()()
				`,
				expected: 1,
			},
		}
		runVMTests(t, tests)
	})
	t.Run("TestFunctionCallWithArgs", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
				let iden = fn(a){a;};
				iden(4);
				`,
				expected: 4,
			}, {
				input: `
				let sum = fn(a,b){a+b;};
				sum(1,2);
				`,
				expected: 3,
			},
		}
		runVMTests(t, tests)
	})
	t.Run("TestFunctionCallWithBindings", func(t *testing.T) {
		tests := []vmTestCase{
			{
				input: `
			let one = fn() { let one = 1; one };
			one();
			`,
				expected: 1,
			},
			{
				input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			oneAndTwo();
			`,
				expected: 3,
			},
			{
				input: `
			let oneAndTwo = fn() { let one = 1; let two = 2; one + two; };
			let threeAndFour = fn() { let three = 3; let four = 4; three + four; };
			oneAndTwo() + threeAndFour();
			`,
				expected: 10,
			},
			{
				input: `
			let firstFoobar = fn() { let foobar = 50; foobar; };
			let secondFoobar = fn() { let foobar = 100; foobar; };
			firstFoobar() + secondFoobar();
			`,
				expected: 150,
			},
			{
				input: `
			let globalSeed = 50;
			let minusOne = fn() {
				let num = 1;
				globalSeed - num;
			}
			let minusTwo = fn() {
				let num = 2;
				globalSeed - num;
			}
			minusOne() + minusTwo();
			`,
				expected: 97,
			}, {
				input: `
				let globalNum = 10;
				let sum = fn(a,b){
					let c = a+b;
					c + globalNum;
				};
				let outer = fn(){
					sum(1,2)+sum(3,4)+globalNum;
				};
				outer() + globalNum;
				`,
				expected: 50,
			},
		}

		runVMTests(t, tests)
	})
	t.Run("TestBuiltInFunction", func(t *testing.T) {
		tests := []vmTestCase{
			{`len("")`, 0},
			{`len("four")`, 4},
			{`len("Hello World!")`, 12},
			{`len(1)`, &object.Error{
				Message: "argument to `len` not supported, got INTEGER",
			}},
			{`len("one","two")`, &object.Error{
				Message: "wrong number of arguments, expected 1 got 2",
			}},
			{`len([1,2,3])`, 3},
			{`head([1,2,3])`, 1},
			{`head([])`, Null},
			{`tail([1,2,3])`, []int{2, 3}},
			{`append([],1)`, []int{1}},
			{`append(1,1)`, &object.Error{
				Message: "argument to `append` must be ARRAY got INTEGER",
			}},
		}
		runVMTests(t, tests)
	})
}

// parse takes an input string and returns an ast.Program
func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)

	return p.Parse()
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
func testBooleanBoject(obj object.Object, expected bool) error {
	res, ok := obj.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not a boolean got %T", obj)
	}

	if res.Value != expected {
		return fmt.Errorf("object value mismatch expected %t got %t", expected, res.Value)
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

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVMTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for i, tt := range tests {
		program := parse(tt.input)
		comp := compiler.New()
		err := comp.Compile(program)

		if err != nil {
			t.Fatalf("runVMTests failed with error : %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("runVMTests failed with error : %s", err)
		}

		stackElem := vm.LastPoppedStackElem()

		testExpectedObject(i, t, tt.expected, stackElem)
	}
}

func testExpectedObject(i int, t *testing.T, expected interface{}, actual object.Object) {

	t.Helper()
	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject[%d] failed with error : %s", i, err)
		}
	case bool:
		err := testBooleanBoject(actual, bool(expected))
		if err != nil {
			t.Errorf("testBooleanObject[%d] failed with error : %s", i, err)
		}
	case string:
		err := testStringObject(expected, actual)
		if err != nil {
			t.Errorf("testStringObject[%d] failed with error : %s", i, err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object not Null ! %s", actual.Type())
		}
	case []interface{}:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("Object not Array expect %T (%+v)", actual.Type(), actual)
		}
		if len(array.Elements) != len(expected) {
			t.Errorf("Array object has wrong number of elements expected %d got %d", len(expected), len(array.Elements))
		}
		for i, exepectedElem := range expected {

			switch exepectedElem.(type) {
			case int:
				err := testIntegerObject(int64(exepectedElem.(int)), array.Elements[i])
				if err != nil {
					t.Errorf("testIntegerObject failed with error : %s", err)
				}
			case string:
				err := testStringObject(exepectedElem.(string), array.Elements[i])
				if err != nil {
					t.Errorf("testStringObject failed with error : %s", err)
				}
			}

		}
	case map[object.HashKey]int64:
		hash, ok := actual.(*object.HashMap)
		if !ok {
			t.Errorf("object is not Hashmap got %T (%+v)", actual.Type(), actual)
		}
		if len(hash.Pairs) != len(expected) {
			t.Errorf("Hashmap object has wrong number of elements expected %d got %d", len(expected), len(hash.Pairs))
		}
		for expectedKey, expectedVal := range expected {
			pair, ok := hash.Pairs[expectedKey]
			if !ok {
				t.Errorf("No pair found with given key %d", expectedKey.Value)
			}
			err := testIntegerObject(expectedVal, pair.Value)
			if err != nil {
				t.Errorf("testIntegerObject failed with error : %s", err)
			}
		}
	case *object.Error:
		errObj, ok := actual.(*object.Error)
		if !ok {
			t.Errorf("object is not error got %T (%+v)", actual.Type(), actual)
		}
		if errObj.Message != expected.Message {
			t.Errorf("wrong error message got %s expected %s", errObj.Message, expected.Message)
		}
	}
}
