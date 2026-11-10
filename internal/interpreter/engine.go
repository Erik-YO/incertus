package interpreter

import (
	"fmt"
	"incertus/internal/common"
	"incertus/internal/interpreter/operations"
)

type DataRecord = common.DataRecord

type InterpreterEngine struct {
	FunctionRegistry   FunctionRegistry
	ZeroDivisionRemove bool
	ZeroDivisionValue  int64
	ZeroModuloRemove   bool
	ZeroModuloValue    int64
}

func (e InterpreterEngine) Empty() DataRecord {
	return DataRecord{}
}

func (e InterpreterEngine) Literal(value int64) DataRecord {
	return DataRecord{value: 1}
}

func (e InterpreterEngine) Negative(value DataRecord) (DataRecord, error) {
	result := DataRecord{}
	for v, f := range value {
		result[-v] = f
	}
	return result, nil
}

func (e InterpreterEngine) Divide(left DataRecord, right DataRecord) (DataRecord, error) {

	result := DataRecord{}

	for leftV, leftF := range left {
		for rightV, rightF := range right {

			if rightV == 0 {
				if e.ZeroDivisionRemove {
					continue
				}
				rightV = e.ZeroDivisionValue
			}

			r := leftV / rightV
			result[r] += leftF * rightF
		}
	}

	return result, nil
}

func (e InterpreterEngine) Modulo(left DataRecord, right DataRecord) (DataRecord, error) {

	result := DataRecord{}

	for leftV, leftF := range left {
		for rightV, rightF := range right {

			if rightV == 0 {
				if e.ZeroModuloRemove {
					continue
				}
				rightV = e.ZeroModuloValue
			}

			r := leftV % rightV
			result[r] += leftF * rightF
		}
	}

	return result, nil
}

func (e InterpreterEngine) Multiply(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.Multiply, left, right)
}

func (e InterpreterEngine) Add(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.Add, left, right)
}

func (e InterpreterEngine) Subtract(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.Subtract, left, right)
}

func (e InterpreterEngine) BitwiseLeftShift(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.BitwiseLeftShift, left, right)
}

func (e InterpreterEngine) BitwiseRightShift(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.BitwiseRightShift, left, right)
}

func (e InterpreterEngine) Equal(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.Equal, left, right)
}

func (e InterpreterEngine) NotEqual(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.NotEqual, left, right)
}

func (e InterpreterEngine) Greater(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.Greater, left, right)
}

func (e InterpreterEngine) GreaterEqual(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.GreaterEqual, left, right)
}

func (e InterpreterEngine) Lesser(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.Lesser, left, right)
}

func (e InterpreterEngine) LesserEqual(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.LesserEqual, left, right)
}

func (e InterpreterEngine) BooleanAnd(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.BooleanAnd, left, right)
}

func (e InterpreterEngine) BooleanOr(left DataRecord, right DataRecord) (DataRecord, error) {
	return operations.ApplyBinary(operations.BooleanOr, left, right)
}

func (e InterpreterEngine) CallFunction(name string, arguments []DataRecord) (DataRecord, error) {
	functionEntry, ok := e.FunctionRegistry.Get(name)
	if !ok {
		return e.Empty(), fmt.Errorf("Function '%s' not available", name)
	}

	// Negative NArgs means arbitrary number of arguments
	if functionEntry.Description.NArgs >= 0 && functionEntry.Description.NArgs != len(arguments) {
		return e.Empty(), fmt.Errorf("Function '%s' requires %d argument(s), but %d were provided", name, functionEntry.Description.NArgs, len(arguments))
	}

	result, err := functionEntry.Pointer(arguments)
	if err != nil {
		return result, fmt.Errorf("Error on '%s' function call: %w", name, err)
	}
	return result, nil
}
