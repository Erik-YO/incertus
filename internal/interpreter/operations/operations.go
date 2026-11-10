package operations

import "incertus/internal/common"

type DataRecord = common.DataRecord

// ====================================

func Multiply(a, b int64) int64 {
	return a * b
}

func Add(a, b int64) int64 {
	return a + b
}

func Subtract(a, b int64) int64 {
	return a - b
}

func BitwiseLeftShift(a, b int64) int64 {
	return a << b
}

func BitwiseRightShift(a, b int64) int64 {
	return a >> b
}

func Equal(a, b int64) int64 {
	if a == b {
		return 1
	}
	return 0
}

func NotEqual(a, b int64) int64 {
	if a != b {
		return 1
	}
	return 0
}

func Greater(a, b int64) int64 {
	if a > b {
		return 1
	}
	return 0
}

func GreaterEqual(a, b int64) int64 {
	if a >= b {
		return 1
	}
	return 0
}

func Lesser(a, b int64) int64 {
	if a < b {
		return 1
	}
	return 0
}

func LesserEqual(a, b int64) int64 {
	if a <= b {
		return 1
	}
	return 0
}

func BooleanAnd(a, b int64) int64 {
	if a != 0 && b != 0 {
		return 1
	}
	return 0
}

func BooleanOr(a, b int64) int64 {
	if a != 0 || b != 0 {
		return 1
	}
	return 0
}

// ====================================

func ApplyBinary(function func(int64, int64) int64, a DataRecord, b DataRecord) (DataRecord, error) {

	result := DataRecord{}

	for aV, aF := range a {
		for bV, bF := range b {
			r := function(aV, bV)
			result[r] += aF * bF
		}
	}

	return result, nil
}
