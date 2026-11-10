package core

import (
	"fmt"
	"slices"
)

type parserState[T any] struct {
	tokens  []Token
	current int
}

func newParser[T any](tokens []Token) parserState[T] {
	return parserState[T]{tokens: tokens}
}

func Parse[T any](tokens []Token) (Expression[T], error) {
	var parser parserState[T] = newParser[T](tokens)
	expr, err := parser.parseExpression()
	if err != nil {
		return expr, err
	}
	err = parser.validatePending()
	if err != nil {
		return expr, err
	}
	return expr, err
}

// ====================================

var booleanTokens []TokenKind = []TokenKind{TokenOpAnd, TokenOpOr}
var comparisonTokens []TokenKind = []TokenKind{TokenOpEqual, TokenOpNotEqual, TokenOpGreater, TokenOpLesser, TokenOpGreaterEqual, TokenOpLesserEqual}
var bitopTokens []TokenKind = []TokenKind{TokenOpBitLeft, TokenOpBitRight}
var termTokens []TokenKind = []TokenKind{TokenOpPlus, TokenOpMinus}
var factorTokens []TokenKind = []TokenKind{TokenOpMultiplication, TokenOpDivision, TokenOpModulo}
var unaryTokens []TokenKind = []TokenKind{TokenOpMinus}

// ====================================

func (state parserState[T]) validatePending() error {

	if state.current >= len(state.tokens) {
		return nil
	}

	current := state.tokens[state.current]

	return fmt.Errorf("Unexpected token '%v'", current)
}

func (state *parserState[T]) parseExpression() (Expression[T], error) {
	return state.parseBoolean()
}

func (state *parserState[T]) parseBoolean() (Expression[T], error) {
	var expr Expression[T]
	expr, err := state.parseComparison()
	if err != nil {
		return expr, err
	}

	for state.match(booleanTokens...) {
		operator := state.tokens[state.current]
		state.current++

		right, err := state.parseComparison()
		if err != nil {
			// Equality check requires a second operand
			return expr, err
		}

		expr = AstBinaryExpression[T]{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (state *parserState[T]) parseComparison() (Expression[T], error) {
	var expr Expression[T]
	expr, err := state.parseBitOp()
	if err != nil {
		return expr, err
	}

	for state.match(comparisonTokens...) {
		operator := state.tokens[state.current]
		state.current++

		right, err := state.parseBitOp()
		if err != nil {
			// Comparison check requires a second operand
			return expr, err
		}

		expr = AstBinaryExpression[T]{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (state *parserState[T]) parseBitOp() (Expression[T], error) {
	var expr Expression[T]
	expr, err := state.parseTerm()
	if err != nil {
		return expr, err
	}

	for state.match(bitopTokens...) {
		operator := state.tokens[state.current]
		state.current++

		right, err := state.parseTerm()
		if err != nil {
			// Bit operation requires a second operand
			return expr, err
		}

		expr = AstBinaryExpression[T]{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (state *parserState[T]) parseTerm() (Expression[T], error) {
	var expr Expression[T]
	expr, err := state.parseFactor()
	if err != nil {
		return expr, err
	}

	for state.match(termTokens...) {
		operator := state.tokens[state.current]
		state.current++

		right, err := state.parseFactor()
		if err != nil {
			// Operation requires a second operand
			return expr, err
		}

		expr = AstBinaryExpression[T]{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (state *parserState[T]) parseFactor() (Expression[T], error) {
	var expr Expression[T]
	expr, err := state.parseUnary()
	if err != nil {
		return expr, err
	}

	for state.match(factorTokens...) {
		operator := state.tokens[state.current]
		state.current++

		right, err := state.parseUnary()
		if err != nil {
			// Operation requires a second operand
			return expr, err
		}

		expr = AstBinaryExpression[T]{left: expr, operator: operator, right: right}
	}

	return expr, nil
}

func (state *parserState[T]) parseUnary() (Expression[T], error) {
	// No left part for unary

	if state.match(unaryTokens...) {
		operator := state.tokens[state.current]
		state.current++

		right, err := state.parseUnary()
		if err != nil {
			// Operation requires an expression
			return right, err
		}

		return AstUnaryExpression[T]{operator: operator, right: right}, nil
	}

	return state.parsePrimary()
}

func (state *parserState[T]) parsePrimary() (Expression[T], error) {

	if state.match(TokenNumber) {
		value := state.tokens[state.current]
		state.current++
		return AstLiteralExpression[T]{value: value}, nil
	}

	if state.match(TokenParenLeft) {
		state.current++
		expr, err := state.parseExpression()
		if err != nil {
			return expr, err
		}
		if !state.match(TokenParenRight) {
			return AstExpressionError[T]{}, fmt.Errorf("Unmatched parenthesis")
		}
		state.current++
		return expr, nil
	}

	if state.match(TokenIdentifier) {
		function := state.tokens[state.current]
		state.current++

		if !state.match(TokenParenLeft) {
			return AstExpressionError[T]{}, fmt.Errorf("Bare function name calling it with 'function(...)'")
		}
		state.current++

		// No arguments
		if state.match(TokenParenRight) {
			state.current++
			return AstFunctionCallExpression[T]{function: function, arguments: nil}, nil
		}

		arguments, err := state.parseCallArguments()
		if err != nil {
			return AstExpressionError[T]{}, err
		}

		if !state.match(TokenParenRight) {
			return AstExpressionError[T]{}, fmt.Errorf("Unclosed parenthesis on function call")
		}
		state.current++

		return AstFunctionCallExpression[T]{function: function, arguments: arguments}, nil
	}

	return AstExpressionError[T]{}, fmt.Errorf("You should not be seeing this")
}

// What's inside the parenthesis, without including them
func (state *parserState[T]) parseCallArguments() ([]Expression[T], error) {
	var expressions []Expression[T]
	var expr Expression[T]

	expr, err := state.parseExpression()
	if err != nil {
		return nil, err
	}
	expressions = append(expressions, expr)

	for state.match(TokenColon) {
		state.current++

		expr, err := state.parseExpression()
		if err != nil {
			return nil, err
		}

		expressions = append(expressions, expr)
	}

	return expressions, nil
}

// ====================================

// Does not advance state's cursor
// that's why we accept state as value, not pointer
func (state parserState[T]) match(tokens ...TokenKind) bool {
	if state.current >= len(state.tokens) {
		return false
	}
	return slices.Contains(tokens, state.tokens[state.current].Kind)
}
