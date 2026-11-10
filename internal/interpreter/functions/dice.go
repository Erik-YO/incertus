package functions

import (
	"incertus/internal/interpreter"
	"incertus/internal/interpreter/operations"
)

// ====================================

func DiceSum1(r []DataRecord) (DataRecord, error) {
	if err := ensureArgumentCount(1, r); err != nil {
		return DataRecord{}, err
	}
	return DiceSum2([]DataRecord{{1: 1}, r[0]})
}

// This is O(N^3)
func DiceSum2(r []DataRecord) (DataRecord, error) {
	if err := ensureArgumentCount(2, r); err != nil {
		return DataRecord{}, err
	}

	a := r[0]
	b := r[1]

	result := DataRecord{}

	for nV, nF := range a {

		tmp := DataRecord{}
		for sidesV, sidesF := range b {
			r := dice(nV, sidesV)

			for iV, iF := range r {
				tmp[iV] += iF * sidesF
			}
		}

		for kV, kF := range tmp {
			result[kV] += kF * nF
		}
	}

	return result, nil
}

func DiceBest(r []DataRecord) (DataRecord, error) {
	if err := ensureArgumentCount(2, r); err != nil {
		return DataRecord{}, err
	}

	return operations.ApplyBinary(maxValue, r[0], r[1])
}

func DiceWorst(r []DataRecord) (DataRecord, error) {
	if err := ensureArgumentCount(2, r); err != nil {
		return DataRecord{}, err
	}

	return operations.ApplyBinary(minValue, r[0], r[1])
}

// ====================================

func maxValue(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func minValue(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func dice(n int64, sides int64) DataRecord {
	if n <= 0 {
		return DataRecord{}
	}

	var current, operand DataRecord
	operand = interpreter.DataRecord{}

	for i := range sides {
		operand[i+1] = 1
	}

	current = operand

	for range n - 1 {
		current = add(current, operand)
	}
	return current
}

func add(a, b DataRecord) DataRecord {

	result := DataRecord{}

	for aV, aF := range a {
		for bV, bF := range b {
			result[aV+bV] += aF * bF
		}
	}

	return result
}
