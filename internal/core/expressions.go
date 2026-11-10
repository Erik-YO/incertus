package core

import (
	"fmt"
)

type Expression[T any] interface {
	String() string
	Process(Engine[T]) (T, error)
}

// ====================================

type AstExpressionError[T any] struct{}

func (e AstExpressionError[T]) String() string { panic("") }

func (e AstExpressionError[T]) Process(Engine[T]) (T, error) { panic("") }

// ====================================

type AstBinaryExpression[T any] struct {
	left     Expression[T]
	operator Token
	right    Expression[T]
}

func (e AstBinaryExpression[T]) String() string {
	return fmt.Sprintf("(%s %s %s)", e.left.String(), e.operator.String(), e.right.String())
}

func (e AstBinaryExpression[T]) Process(engine Engine[T]) (T, error) {

	left, err := e.left.Process(engine)
	if err != nil {
		return left, err
	}
	right, err := e.right.Process(engine)
	if err != nil {
		return right, err
	}

	switch e.operator.Kind {
	case TokenOpAnd:
		return engine.BooleanAnd(left, right)
	case TokenOpOr:
		return engine.BooleanOr(left, right)
	case TokenOpEqual:
		return engine.Equal(left, right)
	case TokenOpNotEqual:
		return engine.NotEqual(left, right)
	case TokenOpGreater:
		return engine.Greater(left, right)
	case TokenOpLesser:
		return engine.Lesser(left, right)
	case TokenOpGreaterEqual:
		return engine.GreaterEqual(left, right)
	case TokenOpLesserEqual:
		return engine.LesserEqual(left, right)
	case TokenOpBitLeft:
		return engine.BitwiseLeftShift(left, right)
	case TokenOpBitRight:
		return engine.BitwiseRightShift(left, right)
	case TokenOpPlus:
		return engine.Add(left, right)
	case TokenOpMinus:
		return engine.Subtract(left, right)
	case TokenOpMultiplication:
		return engine.Multiply(left, right)
	case TokenOpDivision:
		return engine.Divide(left, right)
	case TokenOpModulo:
		return engine.Modulo(left, right)
	default:
		// Probably should panic
		return engine.Empty(), fmt.Errorf("Internal parsing error: Established a non binary operator to a binary operation (%v)", e.operator.Kind)
	}
}

// ====================================

type AstUnaryExpression[T any] struct {
	operator Token
	right    Expression[T]
}

func (e AstUnaryExpression[T]) String() string {
	return fmt.Sprintf("(%s %s)", e.operator.String(), e.right.String())
}

func (e AstUnaryExpression[T]) Process(engine Engine[T]) (T, error) {
	right, err := e.right.Process(engine)
	if err != nil {
		return right, err
	}

	switch e.operator.Kind {
	case TokenOpMinus:
		return engine.Negative(right)
	default:
		// Probably should panic
		return engine.Empty(), fmt.Errorf("Internal parsing error: Established a non unary operator to a unary operation (%v)", e.operator.Kind)
	}
}

// ====================================

type AstLiteralExpression[T any] struct {
	value Token
}

func (e AstLiteralExpression[T]) String() string {
	return e.value.String()
}
func (e AstLiteralExpression[T]) Process(engine Engine[T]) (T, error) {
	if e.value.Kind != TokenNumber {
		return engine.Empty(), fmt.Errorf("Internal parsing error: Got Literal from not a number (%v)", e.value.Kind)
	}
	return engine.Literal(e.value.Number), nil
}

// ====================================

type AstFunctionCallExpression[T any] struct {
	function  Token
	arguments []Expression[T]
}

func (e AstFunctionCallExpression[T]) String() string {
	return fmt.Sprintf("%s%v", e.function.String(), e.arguments)
}
func (e AstFunctionCallExpression[T]) Process(engine Engine[T]) (T, error) {
	if e.function.Kind != TokenIdentifier {
		return engine.Empty(), fmt.Errorf("Internal parsing error: Got function name from not an identifier (%v)", e.function.Kind)
	}

	var arguments []T
	for _, expr := range e.arguments {
		argument, err := expr.Process(engine)
		if err != nil {
			return argument, err
		}
		arguments = append(arguments, argument)
	}

	return engine.CallFunction(e.function.Identifier, arguments)
}
