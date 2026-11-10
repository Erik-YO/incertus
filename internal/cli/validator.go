package cli

import (
	"fmt"
	"incertus/internal/core"
	"incertus/internal/formatters"
	post_processors "incertus/internal/postprocessor"
	"maps"
	"slices"
	"strconv"
	"strings"
)

func validateOptions(r RawOptions) (Options, error) {

	if !r.SingleExpressionSet {
		r.SingleExpression = r.ExpressionSeparator == ""
	}

	if !r.SingleExpression && !r.ForceExpressionSeparator {
		if r.ExpressionSeparator == "" {
			r.ExpressionSeparator = "\n"
		}
		tokens, err := core.Lex(r.ExpressionSeparator)
		if err == nil && len(tokens) > 0 {
			return Options{}, fmt.Errorf("The string provided as '--expression-separator' conflicts with valid expressions '%s'.\nIf you know what you are doing you can switch this check with '--force-expression-separator'", r.ExpressionSeparator)
		}
	}

	onZeroDivisionValue, onZeroDivisionRemove, err := parseOnZero("--on-division-zero", r.OnDivisionZero)
	if err != nil {
		return Options{}, err
	}

	onZeroModuloValue, onZeroModuloRemove, err := parseOnZero("--on-modulo-zero", r.OnModuloZero)
	if err != nil {
		return Options{}, err
	}

	roundSet := false
	var round int
	switch r.Round {
	case "no", "": // Default: no
	default:
		number, err := strconv.Atoi(r.Round)
		if err != nil {
			return Options{}, fmt.Errorf("Invalid string provided as '%s' (%s). If provided should be an integer or 'no'", "--round", r.Round)
		}
		round = number
		roundSet = true
	}

	formatEntry := formatters.DefaultFormatterEntry
	if r.OutputFormat != "" {
		formatEntry, err = parseOption(formatters.FormatterRegistry, "--output", r.OutputFormat)
		if err != nil {
			return Options{}, err
		}
	}
	formatFunction := formatEntry.Function

	resultsSeparator := formatEntry.DefaultResultsSeparator
	if r.ResultsSeparatorSet {
		resultsSeparator = r.ResultsSeparator
	}

	outputSort := post_processors.SortKindValue
	if r.OutputSort != "" {
		outputSort, err = parseOption(map[string]post_processors.SortKind{
			"no":    post_processors.SortKindNo,
			"value": post_processors.SortKindValue,
			"freq":  post_processors.SortKindFrequency,
		}, "--o-sort", r.OutputSort)
		if err != nil {
			return Options{}, err
		}
	}

	outputOrder := post_processors.OrderKindAsc
	if r.OutputOrder != "" {
		outputOrder, err = parseOption(map[string]post_processors.OrderKind{
			"asc":  post_processors.OrderKindAsc,
			"desc": post_processors.OrderKindDesc,
		}, "--o-order", r.OutputOrder)
		if err != nil {
			return Options{}, err
		}
	}

	outputStringValues := false
	if r.OutputStringValues != "" {
		outputStringValues, err = parseOption(map[string]bool{
			"true":  true,
			"yes":   true,
			"y":     true,
			"false": false,
			"no":    false,
			"n":     false,
		}, "--o-string-values", r.OutputStringValues)
		if err != nil {
			return Options{}, err
		}
	}

	outputNormalized := post_processors.NormalizationKind100
	if r.OutputNormalized != "" {
		outputNormalized, err = parseOption(map[string]post_processors.NormalizationKind{
			"no":  post_processors.NormalizationKindNo,
			"1":   post_processors.NormalizationKind1,
			"100": post_processors.NormalizationKind100,
		}, "--o-normalized", r.OutputNormalized)
		if err != nil {
			return Options{}, err
		}
	}

	formatterOptions := formatters.FormatterOptions{
		StringValues: outputStringValues,
		RoundSet:     roundSet,
		Round:        round,
	}

	o := Options{
		Help:             r.Help,
		HelpFunction:     r.HelpFunction,
		HelpFunctionName: r.HelpFunctionName,

		Source:   r.Source,
		FromFile: r.FromFile,

		SingleExpression:    r.SingleExpression,
		ExpressionSeparator: r.ExpressionSeparator,

		OnZeroDivisionValue:  onZeroDivisionValue,
		OnZeroDivisionRemove: onZeroDivisionRemove,

		OnZeroModuloValue:  onZeroModuloValue,
		OnZeroModuloRemove: onZeroModuloRemove,

		OutputFormatter:        formatFunction,
		OutputFormatterOptions: formatterOptions,

		OutputSort:       outputSort,
		OutputOrder:      outputOrder,
		OutputNormalized: outputNormalized,

		ResultsSeparator: resultsSeparator,

		TimeIt: r.TimeIt,
	}

	return o, nil
}

func parseOnZero(option string, provided string) (value int64, remove bool, err error) {

	switch provided {
	case "remove", "": // Default: remove
		remove = true
	case "zero":
		remove = false

	default:
		number, err := strconv.ParseInt(provided, 10, 64)
		if err != nil {
			return value, remove, fmt.Errorf("Invalid string provided as '%s' (%s). Should be an integer, 'zero' or 'remove'", option, provided)
		}
		value = number
		remove = false
	}

	return value, remove, nil
}

func parseOption[T any](choices map[string]T, option string, provided string) (T, error) {
	var selected T

	selected, ok := choices[provided]
	if ok {
		return selected, nil
	}

	keys := slices.Collect(maps.Keys(choices))
	possibleChoices := strings.Join(keys, ", ")

	return selected, fmt.Errorf("Invalid string provided as '%s' (%s). Valid values: %s", option, provided, possibleChoices)
}
