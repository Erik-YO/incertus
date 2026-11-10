package interpreter

import "incertus/internal/registry"

type FunctionRegistry = *registry.FunctionRegistry[DataRecord]

func NewRegistry() FunctionRegistry {
	return registry.NewRegistry[DataRecord]()
}
