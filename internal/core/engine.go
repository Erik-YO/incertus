package core

type Engine[T any] interface {
	Empty() T
	Literal(int64) T

	Negative(T) (T, error)

	Multiply(T, T) (T, error)
	Divide(T, T) (T, error)
	Modulo(T, T) (T, error)
	Add(T, T) (T, error)
	Subtract(T, T) (T, error)
	BitwiseLeftShift(T, T) (T, error)
	BitwiseRightShift(T, T) (T, error)
	Equal(T, T) (T, error)
	NotEqual(T, T) (T, error)
	Greater(T, T) (T, error)
	GreaterEqual(T, T) (T, error)
	Lesser(T, T) (T, error)
	LesserEqual(T, T) (T, error)
	BooleanAnd(T, T) (T, error)
	BooleanOr(T, T) (T, error)

	CallFunction(string, []T) (T, error)
}
