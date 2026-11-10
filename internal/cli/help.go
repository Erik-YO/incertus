package cli

import (
	"fmt"
	"incertus/internal/common"
	"incertus/internal/formatters"
	"incertus/internal/interpreter"
	"incertus/internal/registry"
	"os"
	"strings"
)

func PrintError(err error) {
	println("ERROR ::", err.Error())
}

func PrintCLIError(err error) {
	println("ERROR :: CLI", err.Error())
}

func PrintHelp() {
	builder := strings.Builder{}

	fmt.Fprintf(&builder, "Usage: %s [options...] [expression | --from=FILE]", os.Args[0])

	builder.WriteString(`
Evaluate probabilistic expressions.
`)
	fmt.Fprintf(&builder, "Example: %s 'dice(20) + 3'", os.Args[0])
	builder.WriteString(`

Help:
  -h, --help                    Print this help text.
      --help-function           Print a list of registered functions.
      --help-function=<name>    Print description and usage for a function.


Line options:
      --single-expression       Treats the input as one single expression that
                                can encompass multiple lines.
                                Default if --expression-separator not provided
      --multiple-expressions    Threats the input as one expression per line.
                                Default if --expression-separator provided
      --expression-separator    When using --multiple-expressions, this string
                                is the separator between expressions.
                                Newline by default "\n"
                                This value will be validated to avoid conflicts
                                with valid expressions or operators.
      --force-expression-separator  Skip validation of --expression-separator
                                IMPORTANT: Making the value conflict with valid
                                expressions or operators can cause unexpected results.

On zero division:
      --on-division-zero=[remove|zero|N]  Policy to apply when dividing by zero.
      --on-modulo-zero=[remove|zero|N]    Policy to apply when doing a modulo by zero.
        "remove" (default) Don't count those undefined results towards the evaluation.
        "zero"   Count their frequency as if they were 0.
        N        Count their frequency as if they were an arbitrary integer N.

Rounding:
  -r, --round=[no|N]  Decimal positions to round to. Applied after normalization.

Result separator:
      --results-separator   String to use when evaluating multiple expressions
                            to indicate where the results of one ends and the
                            results of the next begin.
                            Default separator depends on --output format

Output format:
  -o, --output=[<formatter>]	Determines the output format of the evaluation results.
`)

	for name, entry := range formatters.FormatterRegistry {
		fmt.Fprintf(
			&builder,
			`        "%s"	%s
`,
			name,
			entry.Description,
		)

	}

	builder.WriteString(`
Output format modifiers:
    --o-sort=[no|value|freq]  Prints the results sorted by value, frequency or
                              not at all.
                              Default "value"
    --o-order=[asc|desc]      When -o-sort active, orders the results on
                              ascending or descending order.
                              Default "asc"
    --o-string-values=[true|yes|false|no]
                              Prints the value as a string when applicable based
                              on --output. This usually means surrounding the
                              value number with quotes '"'.
                              Default "false"
    --o-normalized=[no|1|100] Normalize the frequencies. This can output
                              decimal numbers.
                              Default "100"
        "no"    means absolute frequencies
        "1"     means normalized
        "100"   means normalized as percentage (normalized then multiplied by 100)

Debug:
  -t, --timeit  Times the steps of each execution and prints a summary at the end.
`)

	builder.WriteString(`
When FILE is '-', read from standard input.
`)

	print(builder.String())
}

func PrintHelpFunctions(r interpreter.FunctionRegistry) {
	builder := strings.Builder{}

	builder.WriteString("Available functions:\n")

	for _, name := range r.Names() {
		entry, ok := r.Get(name)
		if !ok {
			continue
		}

		builder.WriteString(printFunctionEntry(name, entry))
		builder.WriteRune('\n')
	}

	println(builder.String())
}

func PrintHelpFunction(name string, r interpreter.FunctionRegistry) {

	entry, found := r.Get(name)

	if !found {
		println(fmt.Sprintf("Function '%s' not found\nUse --help-function to get a list of available functions", name))
		return
	}

	entryStr := printFunctionEntry(name, entry)

	println(fmt.Sprintf("Usage: %s\n", entryStr))
}

func printFunctionEntry(name string, entry registry.FunctionEntry[common.DataRecord]) string {

	builder := strings.Builder{}
	builder.WriteString(name)

	builder.WriteRune('(')

	if entry.Description.NArgs > 0 {

		for idx := range entry.Description.NArgs - 1 {
			fmt.Fprintf(&builder, "ARG%d, ", idx+1)
		}

		if entry.Description.NArgs > 1 {
			fmt.Fprintf(&builder, "ARG%d", entry.Description.NArgs)

		} else {
			builder.WriteString("ARG")

		}

	} else if entry.Description.NArgs < 0 {
		builder.WriteString("ARGS...")
	}

	builder.WriteRune(')')
	builder.WriteString("\n    ")
	builder.WriteString(entry.Description.Help)
	return builder.String()
}
