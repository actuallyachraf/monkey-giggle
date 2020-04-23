package eval

import (
	"testing"

	"github.com/actuallyachraf/monkey-giggle/lexer"
	"github.com/actuallyachraf/monkey-giggle/object"
	"github.com/actuallyachraf/monkey-giggle/parser"
)

func TestEval(t *testing.T) {

	t.Run("TestEvalIntegerExpression", func(t *testing.T) {

		tests := []struct {
			input    string
			expected int64
		}{
			{"5", 5},
			{"10", 10},
		}

		for _, tt := range tests {
			evaled := testEval(tt.input)
			testIntegerObject(t, evaled, tt.expected)
		}
	})
	t.Run("TestEvalBooleanExpression", func(t *testing.T) {
		tests := []struct {
			input    string
			expected bool
		}{
			{"true", true},
			{"false", false},
			{"1 < 2", true},
			{"1 > 2", false},
			{"1 < 1", false},
			{"1 > 1", false},
			{"1 == 1", true},
			{"1 == 2", false},
			{"1 != 2", true},
			{"1 != 1", false},
			{"1 <= 1", true},
			{"1 >= 1", true},
			{"true == true", true},
			{"false == false", true},
			{"true == false", false},
			{"true != false", true},
			{"false != true", true},
			{"(1 < 2) == true", true},
			{"(1 < 2) == false", false},
			{"(1 > 2) == true", false},
			{"(1 > 2) == false", true},
		}

		for _, tt := range tests {
			evaled := testEval(tt.input)
			testBooleanBoject(t, evaled, tt.expected)
		}
	})
	t.Run("TestEvalBangOperator", func(t *testing.T) {
		tests := []struct {
			input    string
			expected bool
		}{
			{"!true", false},
			{"!false", true},
			{"!5", false},
			{"!!true", true},
			{"!!false", false},
			{"!!5", true},
		}

		for _, tt := range tests {
			evaled := testEval(tt.input)
			testBooleanBoject(t, evaled, tt.expected)
		}
	})
	t.Run("TestEvalMinusOperator", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{"-5", -5},
			{"-10", -10},
		}

		for _, tt := range tests {
			evaled := testEval(tt.input)
			testIntegerObject(t, evaled, tt.expected)
		}
	})
	t.Run("TestEvalIntegerExpression", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{"5", 5},
			{"10", 10},
			{"-5", -5},
			{"-10", -10},
			{"5 + 5 + 5 + 5 - 10", 10},
			{"2 * 2 * 2 * 2 * 2", 32},
			{"-50 + 100 + -50", 0},
			{"5 * 2 + 10", 20},
			{"5 + 2 * 10", 25},
			{"20 + 2 * -10", 0},
			{"50 / 2 * 2 + 10", 60},
			{"2 * (5 + 10)", 30},
			{"3 * 3 * 3 + 10", 37},
			{"3 * (3 * 3) + 10", 37},
			{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
		}

		for _, tt := range tests {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		}
	})
	t.Run("TestEvalIfElseExpression", func(t *testing.T) {
		tests := []struct {
			input    string
			expected interface{}
		}{
			{"if (true) { 10 }", 10},
			{"if (false) { 10 }", nil},
			{"if (1) { 10 }", 10},
			{"if (1 < 2) { 10 }", 10},
			{"if (1 > 2) { 10 }", nil},
			{"if (1 > 2) { 10 } else { 20 }", 20},
			{"if (1 < 2) { 10 } else { 20 }", 10},
		}

		for _, tt := range tests {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		}
	})
	t.Run("TestEvalReturnStatement", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{"return 10;", 10},
			{"return 10; 9;", 10},
			{"return 2 * 5; 9;", 10},
			{"9; return 2 * 5; 9;", 10},
			{"if (10 > 1) { return 10; }", 10},
			{
				`
	if (10 > 1) {
	  if (10 > 1) {
		return 10;
	  }

	  return 1;
	}
	`,
				10,
			},
		}

		for _, tt := range tests {
			evaluated := testEval(tt.input)
			testIntegerObject(t, evaluated, tt.expected)
		}
	})
	t.Run("TestEvalErrorHandling", func(t *testing.T) {
		tests := []struct {
			input           string
			expectedMessage string
		}{
			{
				"5 + true;",
				"type mismatch: INTEGER + BOOLEAN",
			},
			{
				"5 + true; 5;",
				"type mismatch: INTEGER + BOOLEAN",
			},
			{
				"-true",
				"unknown operator: -BOOLEAN",
			},
			{
				"true + false;",
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"true + false + true + false;",
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"5; true + false; 5",
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"if (10 > 1) { true + false; }",
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				`
	if (10 > 1) {
	  if (10 > 1) {
		return true + false;
	  }

	  return 1;
	}
	`,
				"unknown operator: BOOLEAN + BOOLEAN",
			},
			{
				"foobar",
				"identifier not found: foobar",
			},
		}

		for _, tt := range tests {
			evaluated := testEval(tt.input)

			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("no error object returned. got=%T(%+v)",
					evaluated, evaluated)
				continue
			}

			if errObj.Message != tt.expectedMessage {
				t.Errorf("wrong error message. expected=%q, got=%q",
					tt.expectedMessage, errObj.Message)
			}
		}
	})
	t.Run("TestEvalLetStatement", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{"let a = 5; a;", 5},
			{"let a = 5 * 5; a;", 25},
			{"let a = 5; let b = a; b;", 5},
			{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
		}

		for _, tt := range tests {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		}
	})
	t.Run("TestEvalFunctionStatement", func(t *testing.T) {
		input := "fn(x) { x + 2; };"

		evaluated := testEval(input)
		fn, ok := evaluated.(*object.Function)
		if !ok {
			t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
		}

		if len(fn.Parameters) != 1 {
			t.Fatalf("function has wrong parameters. Parameters=%+v",
				fn.Parameters)
		}

		if fn.Parameters[0].String() != "x" {
			t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
		}

		expectedBody := "(x + 2)"

		if fn.Body.String() != expectedBody {
			t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
		}
	})
	t.Run("TestEvalFunctionApplication", func(t *testing.T) {
		tests := []struct {
			input    string
			expected int64
		}{
			{"let identity = fn(x) { x; }; identity(5);", 5},
			{"let identity = fn(x) { return x; }; identity(5);", 5},
			{"let double = fn(x) { x * 2; }; double(5);", 10},
			{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
			{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
			{"fn(x) { x; }(5)", 5},
			{"let addMod = fn(x,y,m){return (x+y)%m;}; addMod(5,10,3)", 0},
		}

		for _, tt := range tests {
			testIntegerObject(t, testEval(tt.input), tt.expected)
		}
	})
	t.Run("TestEvalEnclosedEnv", func(t *testing.T) {
		input := `
		let first = 10;
		let second = 10;
		let third = 10;

		let ourFunction = fn(first) {
		  let second = 20;

		  first + second + third;
		};

		ourFunction(20) + first + second;`

		testIntegerObject(t, testEval(input), 70)
	})
	t.Run("TestEvalStringExpression", func(t *testing.T) {

		input := `"hello world!"`

		evaled := testEval(input)

		str, ok := evaled.(*object.String)
		if !ok {
			t.Fatalf("object is not String got %T , (%+v)", evaled, evaled)
		}

		if str.Value != "hello world!" {
			t.Fatalf("wrong object value expected %s got %s", input, str.Value)
		}
	})
	t.Run("TestEvalStringConcat", func(t *testing.T) {
		input := `"Hello" +" " + "World!"`
		evaled := testEval(input)
		str, ok := evaled.(*object.String)
		if !ok {
			t.Fatalf("object is not String got %T , (%+v)", evaled, evaled)
		}

		if str.Value != "Hello World!" {
			t.Fatalf("wrong object value expected %s got %s", input, str.Value)
		}
	})
	t.Run("TestEvalBuiltInFunction", func(t *testing.T) {
		tests := []struct {
			input    string
			expected interface{}
		}{
			{`len("")`, 0},
			{`len("four")`, 4},
			{`len("Hello World!")`, 12},
			{`len(1)`, "argument to `len` not supported, got INTEGER"},
			{`len("one","two")`, "wrong number of arguments, expected 1 got 2"},
			{`head([1,2,3])`, 1},
			{`tail([1,2,3])`, []int{2, 3}},
			{`append([],1)`, []int{1}},
		}

		for _, tt := range tests {
			evaled := testEval(tt.input)
			switch expected := tt.expected.(type) {
			case int:
				testIntegerObject(t, evaled, int64(expected))
			case string:
				errObj, ok := evaled.(*object.Error)
				if !ok {
					t.Errorf("object is not an error got %T (%+v)", evaled, evaled)
					continue
				}
				if errObj.Message != expected {
					t.Errorf("wrong error message ! expected %q got %q", expected, errObj.Message)
				}
			case []int:
				array, ok := evaled.(*object.Array)
				if !ok {
					t.Errorf("obj not Array. got=%T (%+v)", evaled, evaled)
					continue
				}

				if len(array.Elements) != len(expected) {
					t.Errorf("wrong num of elements. want=%d, got=%d",
						len(expected), len(array.Elements))
					continue
				}

				for i, expectedElem := range expected {
					testIntegerObject(t, array.Elements[i], int64(expectedElem))
				}
			}
		}
	})
	t.Run("TestEvalArrayLiteral", func(t *testing.T) {
		input := "[1, 2 * 2, 3 + 3]"

		evaluated := testEval(input)
		result, ok := evaluated.(*object.Array)
		if !ok {
			t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
		}

		if len(result.Elements) != 3 {
			t.Fatalf("array has wrong num of elements. got=%d",
				len(result.Elements))
		}

		testIntegerObject(t, result.Elements[0], 1)
		testIntegerObject(t, result.Elements[1], 4)
		testIntegerObject(t, result.Elements[2], 6)
	})
	t.Run("TestEvalIndexExpression", func(t *testing.T) {
		tests := []struct {
			input    string
			expected interface{}
		}{
			{
				"[1, 2, 3][0]",
				1,
			},
			{
				"[1, 2, 3][1]",
				2,
			},
			{
				"[1, 2, 3][2]",
				3,
			},
			{
				"let i = 0; [1][i];",
				1,
			},
			{
				"[1, 2, 3][1 + 1];",
				3,
			},
			{
				"let myArray = [1, 2, 3]; myArray[2];",
				3,
			},
			{
				"let myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
				6,
			},
			{
				"let myArray = [1, 2, 3]; let i = myArray[0]; myArray[i]",
				2,
			},
			{
				"[1, 2, 3][3]",
				nil,
			},
			{
				"[1, 2, 3][-1]",
				nil,
			},
		}

		for _, tt := range tests {
			evaluated := testEval(tt.input)
			integer, ok := tt.expected.(int)
			if ok {
				testIntegerObject(t, evaluated, int64(integer))
			} else {
				testNullObject(t, evaluated)
			}
		}
	})

}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.Parse()
	env := object.NewEnv()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	res, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not an integer got %T", obj)
		return false
	}

	if res.Value != expected {
		t.Errorf("object value mismatch expected %d got %d", expected, res.Value)
		return false
	}

	return true
}

func testBooleanBoject(t *testing.T, obj object.Object, expected bool) bool {
	res, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not a boolean got %T", obj)
		return false
	}

	if res.Value != expected {
		t.Errorf("object value mismatch expected %t got %t", expected, res.Value)
		return false
	}

	return true
}
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
