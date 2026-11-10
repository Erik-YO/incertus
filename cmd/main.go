package main

import (
	"incertus/internal/cli"
	"incertus/internal/interpreter"
	"incertus/internal/interpreter/functions"
	"incertus/internal/registry"
	"os"
)

func main() {

	args := os.Args[1:]

	r := interpreter.NewRegistry()
	err := mainRegister(r)
	if err != nil {
		cli.PrintCLIError(err)
		os.Exit(1)
	}

	status := cli.Main(args, r)
	os.Exit(status)
}

func mainRegister(r interpreter.FunctionRegistry) error {

	functions := []struct {
		name    string
		pointer registry.Function[interpreter.DataRecord]
		nArgs   int
		help    string
	}{
		// Dice functions
		{"dice", functions.DiceSum1, 1, "Throw a D-sided dice"},
		{"dice.sum", functions.DiceSum2, 2, "Throw N D-sided dice and add their results"},
		{"dice.best", functions.DiceBest, 2, "Keep the highest value"},
		{"dice.worst", functions.DiceWorst, 2, "Keep the lowest value"},
		// Statistics
		{"stats.mean", functions.StatMean, 1, "Calculate the mean"},
		{"stats.median", functions.StatMedian, 1, "Calculate the median value"},
	}

	for _, entry := range functions {

		err := r.Register(entry.name, entry.pointer, entry.nArgs, entry.help)
		if err != nil {
			return err
		}
	}

	return nil
}
