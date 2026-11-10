package cli

import (
	"incertus/internal/formatters"
	post_processors "incertus/internal/postprocessor"
)

type RawOptions struct {
	Help             bool
	HelpFunction     bool
	HelpFunctionName string

	// Input options
	Source   string
	FromFile bool

	// Line options
	SingleExpression         bool
	SingleExpressionSet      bool
	ExpressionSeparator      string // Validate to avoid parsing conflicts
	ForceExpressionSeparator bool

	// Zero division
	OnDivisionZero string
	OnModuloZero   string

	Round string

	// Output format
	OutputFormat        string
	OutputSort          string
	OutputOrder         string
	OutputStringValues  string
	OutputNormalized    string
	ResultsSeparator    string
	ResultsSeparatorSet bool

	// Debug
	TimeIt bool
}

type Options struct {
	Help             bool
	HelpFunction     bool
	HelpFunctionName string

	// Input options
	Source   string
	FromFile bool

	// Line options
	SingleExpression    bool
	ExpressionSeparator string // Already validated (or not)

	// Zero division
	OnZeroDivisionRemove bool
	OnZeroDivisionValue  int64
	OnZeroModuloRemove   bool
	OnZeroModuloValue    int64

	// Output format
	OutputFormatter        formatters.FormatterFunc
	OutputFormatterOptions formatters.FormatterOptions

	OutputSort       post_processors.SortKind
	OutputOrder      post_processors.OrderKind
	OutputNormalized post_processors.NormalizationKind
	ResultsSeparator string

	// Debug
	TimeIt bool
}
