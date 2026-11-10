package core

import (
	"testing"
)

func TestParse(t *testing.T) {

	tests := []struct {
		expression  string
		shouldError bool
	}{
		{"", true},
		{"1", false},
		{"123", false},
		{" 1", false},
		{"    12    ", false},

		{"func", true},
		{" func ", true},
		{" a", true},
		{"a ", true},

		{"+", true},
		{" + ", true},
		{" +", true},
		{"+ ", true},

		{"==", true},

		{" > ", true},
		{" >=", true},
		{"<", true},
		{" <= ", true},
		{"//", true},

		{"1 + 1", false},
		{"2 * 3", false},
		{"4 + 5 * 6", false},
		{"1 + -2", false},
		{"4 * (5 + 6)", false},
		{"(4 - (5 + 6)) * 7", false},
		{"(4-(5+6))*7", false},

		{"dice.sum(6)", false},
		{"dice.sum(2,6)", false},
		{"dice.sum(2, 6)", false},

		// (1 << 63) - 1 -> 9223372036854775807
		{"9223372036854775807", false},
		// 1 << 62 -> 4611686018427387904
		{"4611686018427387904", false},
	}

	// Expression
	for _, test := range tests {
		evaluate(t, test.expression, test.shouldError)
		// name := fmt.Sprintf("Expression-%03d:%s", i, test.expression)
		// t.Run(name, func(t *testing.T) {
		// })
	}
}

func evaluate(t *testing.T, expression string, shouldError bool) {

	tokens, err := Lex(expression)
	if err != nil {
		t.Fatalf("Lexing should not error for \"%s\"; got %s", expression, err.Error())
	}

	ast, err := Parse[map[int64]uint64](tokens)

	if err != nil {
		if shouldError {
			return
		}
		t.Fatalf("Parsing should not error for \"%s\"; got %s", expression, err.Error())
	}

	if shouldError {
		t.Fatalf("Parsing should error for \"%s\"; got %v", expression, ast)
	}

	t.Logf("\n\n%s\nAST: %v\n\n", expression, ast)
}
