package cli

import (
	"fmt"
	"slices"
	"strings"
)

type cliReader struct {
	args    []string
	current int
	state   RawOptions
}

const (
	flagPrefix         = "-"
	flagValueSeparator = "="
)

func readRawOptions(args []string) (RawOptions, error) {

	expandedArgs := expandArguments(args)

	reader := cliReader{args: expandedArgs, current: 0}

	var err error
	var running bool = true

	for running {
		running, err = reader.next()
		if err != nil {
			return reader.state, err
		}
	}

	return reader.state, nil
}

func (r *cliReader) next() (keepRunning bool, err error) {

	if r.current >= len(r.args) {
		return false, nil
	}

	if r.match("--help", "-h") {
		r.state.Help = true
		return false, nil
	}

	if r.match("--help-function") {
		r.state.HelpFunction = true
		name, found := r.fetch()
		if found {
			r.state.HelpFunctionName = name
		}
		return false, nil
	}

	if r.match("--from", "-f") {
		filename, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--from' must provide file path")
		}
		r.state.FromFile = true
		r.state.Source = filename
		return true, nil
	}

	if r.match("--single-expression") {
		r.state.SingleExpression = true
		r.state.SingleExpressionSet = true
		return true, nil
	}

	if r.match("--multiple-expressions") {
		if r.state.SingleExpression {
			return false, fmt.Errorf("Argument '--multiple-expressions' is incompatible with '--single-expression'")
		}
		r.state.SingleExpressionSet = true
		return true, nil
	}

	if r.match("--expression-separator") {
		if r.state.ExpressionSeparator != "" {
			return false, fmt.Errorf("Duplicated argument '--expression-separator'|'--force-expression-separator'")
		}

		sep, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--expression-separator' must provide separator string between different expressions")
		}
		r.state.ExpressionSeparator = sep
		r.state.ForceExpressionSeparator = false
		return true, nil
	}

	if r.match("--force-expression-separator") {
		if r.state.ExpressionSeparator != "" {
			return false, fmt.Errorf("Duplicated argument '--expression-separator'|'--force-expression-separator'")
		}

		sep, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--force-expression-separator' must provide separator string between different expressions")
		}
		r.state.ExpressionSeparator = sep
		r.state.ForceExpressionSeparator = true
		return true, nil
	}

	if r.match("--on-division-zero") {
		onDivisionZero, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--on-division-zero' must provide policy to apply when dividing by zero")
		}
		r.state.OnDivisionZero = onDivisionZero
		return true, nil
	}

	if r.match("--on-modulo-zero") {
		onModuloZero, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--on-modulo-zero' must provide policy to apply when dividing by zero")
		}
		r.state.OnModuloZero = onModuloZero
		return true, nil
	}

	if r.match("--round", "-r") {
		round, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--round' must provide the decimal positions to round to")
		}
		r.state.Round = round
		return true, nil
	}

	if r.match("--results-separator") {
		sep, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--results-separator' must provide separator string between different results")
		}
		r.state.ResultsSeparator = sep
		r.state.ResultsSeparatorSet = true
		return true, nil
	}

	if r.match("--output", "-o") {
		output, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--output' must provide the target format")
		}
		r.state.OutputFormat = output
		return true, nil
	}

	if r.match("--output-sort", "--o-sort") {
		outputSort, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--output-sort' must indicate how to sort the results of the evaluations (by value, by frequency or not at all)")
		}
		r.state.OutputSort = outputSort
		return true, nil
	}

	if r.match("--output-order", "--o-order") {
		outputOrder, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--output-order' must indicate how to order the results of the evaluations (ascending or descending)")
		}
		r.state.OutputOrder = outputOrder
		return true, nil
	}

	if r.match("--output-string-values", "--o-string-values") {
		stringValues, found := r.fetch()
		if !found {
			// Accepts true|yes|false|no, it is false by default
			// on flag but no value, it serves as activation signal
			stringValues = "true"
		}
		r.state.OutputStringValues = stringValues
		return true, nil
	}

	if r.match("--output-normalized", "--o-normalized") {
		normalized, found := r.fetch()
		if !found {
			return false, fmt.Errorf("Argument '--output-normalized' must provide normalization value (1, 100, or 'no'|'false')")
		}
		r.state.OutputNormalized = normalized
		return true, nil
	}

	if r.match("--timeit", "-t") {
		r.state.TimeIt = true
		return true, nil
	}

	// Last argument only if file source not already provided
	if r.current+1 == len(r.args) && !r.state.FromFile {
		expressions := r.args[r.current]
		r.current++
		r.state.Source = expressions
		return false, nil
	}

	// Error handling
	return false, fmt.Errorf("Unrecognized argument '%s'", r.args[r.current])
}

// Captures the next on coincidence
func (r *cliReader) match(flags ...string) (found bool) {
	if r.current >= len(r.args) {
		return false
	}

	if slices.Contains(flags, r.args[r.current]) {
		r.current++
		return true
	}
	return false
}

// Captures the next if there is one
func (r *cliReader) fetch() (arg string, found bool) {
	if r.current >= len(r.args) {
		return "", false
	}
	arg = r.args[r.current]
	r.current++
	return arg, true
}

func expandArguments(args []string) []string {

	// Minimal size we are going to need
	result := make([]string, 0, len(args))

	for _, arg := range args {
		if !strings.HasPrefix(arg, flagPrefix) {
			// Not a flag, add as is
			result = append(result, arg)
			continue
		}

		// Flag, split by '='
		before, after, found := strings.Cut(arg, flagValueSeparator)

		result = append(result, before)
		if found {
			result = append(result, after)
		}

	}

	return result
}
