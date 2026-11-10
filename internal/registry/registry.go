package registry

import (
	"fmt"
	"incertus/internal/core"
)

type FunctionDescription struct {
	Name  string
	NArgs int
	Help  string
}

type Function[T any] func([]T) (T, error)

type FunctionEntry[T any] struct {
	Pointer     Function[T]
	Description FunctionDescription
}

type FunctionRegistry[T any] struct {
	functions map[string]FunctionEntry[T]

	// Maintain insertion order for consistent iteration
	names []string
}

func NewRegistry[T any]() *FunctionRegistry[T] {
	return &FunctionRegistry[T]{
		functions: map[string]FunctionEntry[T]{},
		names:     nil,
	}
}

func (r *FunctionRegistry[T]) Register(name string, function Function[T], nArgs int, help string) error {
	if name == "" {
		return fmt.Errorf("Invalid function name, must have at least 1 character")
	}

	_, exists := r.functions[name]
	if exists {
		return fmt.Errorf("Function with name '%s' already registered", name)
	}

	if !core.IsValidIdentifier(name) {
		return fmt.Errorf("Function name '%s' is not a valid identifier", name)
	}

	r.functions[name] = FunctionEntry[T]{
		Pointer: function,
		Description: FunctionDescription{
			Name:  name,
			NArgs: nArgs,
			Help:  help,
		},
	}
	r.names = append(r.names, name)

	return nil
}

func (r FunctionRegistry[T]) Get(name string) (FunctionEntry[T], bool) {
	entry, ok := r.functions[name]
	return entry, ok
}

func (r FunctionRegistry[T]) Names() []string {
	return r.names
}
