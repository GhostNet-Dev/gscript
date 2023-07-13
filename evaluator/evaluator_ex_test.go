package evaluator

import (
	"testing"

	"github.com/GhostNet-Dev/glambda/object"
)

func TestNewBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`int("1")`, 1},
		{`string(1)`, "1"},
	}
	for i, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected), i)
		case string:
			errObj, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("object is Error. got=%T (%+v)", evaluated, evaluated)
				continue
			}
			if errObj.Value != expected {
				t.Errorf("wrong Result. expected=%q, got=%q", expected, errObj.Value)
			}
		}
	}
}

func TestHashVariable(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`let a = {"foo": 5};a["foo"];`, 5},
		{`{"foo": 5}["foo"]`, 5},
		{"let myArray = [1, 2, 3];myArray[2];", 3},
	}
	for i, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer), i)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestUseStruct(t *testing.T) {
	input := `type a struct {}
	a b;`
	evaluated := testEval(input)
	_, ok := evaluated.(*object.Struct)
	if !ok {
		t.Fatalf("Eval didn't return Identifier. got=%T (%+v)", evaluated, evaluated)
	}
}

func TestTypeStruct(t *testing.T) {
	input := `type a struct
	{
		let b = 0;
	}`
	evaluated := testEval(input)
	_, ok := evaluated.(*object.Struct)
	if !ok {
		t.Fatalf("Eval didn't return Struct. got=%T (%+v)", evaluated, evaluated)
	}
}

func TestAssignStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 4;a = 5;a;", 5},
	}

	for i, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, i)
	}
}

func TestLetIdentifierStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; let b = a; b;", 5},
		{"let a = 5; let b = a; let c = a + b + 5; c;", 15},
	}

	for i, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, i)
	}
}

func TestForExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"for (let a = 0;a < 5;a = a + 1) {};a;", 5},
		{"let a = 0;for (a = 0;a < 5;a = a + 1) {};a;", 5},
		{"let a = 0;let b = 0;for (a = 0;a < 5;a = a + 1) {b = b + 1;};b;", 5},
	}

	for i, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, i)
	}
}
func TestFunctionNameObject(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let x = 1;let func = fn() { let x = 0;x = x + 2; };func();x", 1},
		{"let x = 1;let func = fn() { x = x + 2; };func();x", 3},
		{"let x = 1;let func = fn() { let x = 0;x + 2; };func();", 2},
	}

	for i, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, i)
	}
}
