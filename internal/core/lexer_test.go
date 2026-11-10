package core

import (
	"fmt"
	"testing"
)

func TestLexTokenType(t *testing.T) {

	tests := []struct {
		expression  string
		result      []TokenKind
		shouldError bool
	}{
		{"", []TokenKind{}, false},
		{"1", []TokenKind{TokenNumber}, false},
		{"123", []TokenKind{TokenNumber}, false},
		{" 1", []TokenKind{TokenNumber}, false},
		{"01", []TokenKind{}, true},
		{"    12    ", []TokenKind{TokenNumber}, false},

		{"func", []TokenKind{TokenIdentifier}, false},
		{" func ", []TokenKind{TokenIdentifier}, false},
		{" a", []TokenKind{TokenIdentifier}, false},
		{"a ", []TokenKind{TokenIdentifier}, false},
		{"_a ", []TokenKind{}, true},
		{"_", []TokenKind{}, true},
		{".a", []TokenKind{}, true},
		{"_1a ", []TokenKind{}, true},

		{"+", []TokenKind{TokenOpPlus}, false},
		{" + ", []TokenKind{TokenOpPlus}, false},
		{" +", []TokenKind{TokenOpPlus}, false},
		{"+ ", []TokenKind{TokenOpPlus}, false},
		{"_+ ", []TokenKind{}, true},
		{"_1+ ", []TokenKind{}, true},

		{"==", []TokenKind{TokenOpEqual}, false},
		{"=", []TokenKind{}, true},

		{" > ", []TokenKind{TokenOpGreater}, false},
		{" >=", []TokenKind{TokenOpGreaterEqual}, false},
		{"<", []TokenKind{TokenOpLesser}, false},
		{" <= ", []TokenKind{TokenOpLesserEqual}, false},
		{"//", []TokenKind{TokenOpDivision}, false},
		{"/", []TokenKind{}, true},

		{"1 + 1", []TokenKind{TokenNumber, TokenOpPlus, TokenNumber}, false},
		{"2 * 3", []TokenKind{TokenNumber, TokenOpMultiplication, TokenNumber}, false},
		{"4 + 5 * 6", []TokenKind{TokenNumber, TokenOpPlus, TokenNumber, TokenOpMultiplication, TokenNumber}, false},
		{"1 + -2", []TokenKind{TokenNumber, TokenOpPlus, TokenOpMinus, TokenNumber}, false},
		{"4 * (5 + 6)", []TokenKind{TokenNumber, TokenOpMultiplication, TokenParenLeft, TokenNumber, TokenOpPlus, TokenNumber, TokenParenRight}, false},
		{"(4 - (5 + 6)) * 7", []TokenKind{TokenParenLeft, TokenNumber, TokenOpMinus, TokenParenLeft, TokenNumber, TokenOpPlus, TokenNumber, TokenParenRight, TokenParenRight, TokenOpMultiplication, TokenNumber}, false},
		{"(4-(5+6))*7", []TokenKind{TokenParenLeft, TokenNumber, TokenOpMinus, TokenParenLeft, TokenNumber, TokenOpPlus, TokenNumber, TokenParenRight, TokenParenRight, TokenOpMultiplication, TokenNumber}, false},

		{"dice.sum(6)", []TokenKind{TokenIdentifier, TokenParenLeft, TokenNumber, TokenParenRight}, false},
		{"dice.sum(2,6)", []TokenKind{TokenIdentifier, TokenParenLeft, TokenNumber, TokenColon, TokenNumber, TokenParenRight}, false},
		{"dice.sum(2, 6)", []TokenKind{TokenIdentifier, TokenParenLeft, TokenNumber, TokenColon, TokenNumber, TokenParenRight}, false},

		// 1 << 64 -> 18446744073709551616
		{"18446744073709551616", []TokenKind{}, true},
		// 1 << 63 -> 9223372036854775808
		{"9223372036854775808", []TokenKind{}, true},
		// (1 << 63) - 1 -> 9223372036854775807
		{"9223372036854775807", []TokenKind{TokenNumber}, false},
		// 1 << 62 -> 4611686018427387904
		{"4611686018427387904", []TokenKind{TokenNumber}, false},
	}

	// Expression
	for i, test := range tests {
		name := fmt.Sprintf("Expression-%03d:%s", i, test.expression)
		t.Run(name, func(t *testing.T) {
			testExpectedTokenKinds(t, test.expression, test.result, test.shouldError)
		})
	}

	comments := []string{
		"",
		"#",
		"-------\n",
		"abcde",
		"1+1+#1+1",
		" abcd",
		" 1728 3993",
	}

	// Expression with comments
	for i, test := range tests {
		for j, comment := range comments {
			expr := fmt.Sprintf("%s#%s", test.expression, comment)
			name := fmt.Sprintf("Commented-%03d-%02d:%s", i, j, expr)
			t.Run(name, func(t *testing.T) {
				testExpectedTokenKinds(t, expr, test.result, test.shouldError)
			})
		}
	}
}

func testExpectedTokenKinds(t *testing.T, expression string, expected []TokenKind, shouldError bool) {

	result, err := Lex(expression)

	if shouldError {
		if err == nil {
			t.Fatalf("Lexing should return error for \"%s\"; got %v", expression, result)
		}
	} else {
		if err != nil {
			t.Fatalf("Lexing should not error for \"%s\"; got %s", expression, err.Error())
		}
		resultKinds := []TokenKind{}
		for _, token := range result {
			resultKinds = append(resultKinds, token.Kind)
		}
		if len(expected) != len(resultKinds) {
			t.Fatalf(
				"Lexing should return %d tokens for \"%s\"; got %d:\nExpected: %v\nResult:   %v",
				len(expected), expression, len(resultKinds), expected, resultKinds)
		}
		for i := range expected {
			if expected[i] != resultKinds[i] {
				t.Fatalf(
					"Lexing token kind mismatch on index %d:\nExpected: %v\nResult:   %v",
					i, expected, resultKinds)
			}
		}
	}

}
