package evaluator

import "testing"

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
	}

	for i, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected, i)
	}
}
