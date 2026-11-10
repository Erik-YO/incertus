package cli

import (
	"fmt"
	"incertus/internal/core"
	"incertus/internal/formatters"
	"incertus/internal/input"
	"incertus/internal/interpreter"
	post_processors "incertus/internal/postprocessor"
	"io"
	"time"
)

func Main(args []string, registry interpreter.FunctionRegistry) int {
	var err error

	raw, err := readRawOptions(args)
	if err != nil {
		PrintCLIError(err)
		return 1
	}

	opts, err := validateOptions(raw)
	if err != nil {
		PrintCLIError(err)
		return 1
	}

	if opts.Help {
		PrintHelp()
		return 0
	}

	if opts.HelpFunction {
		if opts.HelpFunctionName == "" {
			PrintHelpFunctions(registry)
		} else {
			PrintHelpFunction(opts.HelpFunctionName, registry)
		}
		return 0
	}

	inputReader, err := input.CreateReader(
		opts.Source,
		opts.FromFile,
		opts.SingleExpression,
		opts.ExpressionSeparator,
	)
	if err != nil {
		PrintCLIError(err)
		return 1
	}

	engine := interpreter.InterpreterEngine{
		FunctionRegistry:   registry,
		ZeroDivisionRemove: opts.OnZeroDivisionRemove,
		ZeroDivisionValue:  opts.OnZeroDivisionValue,
		ZeroModuloRemove:   opts.OnZeroModuloRemove,
		ZeroModuloValue:    opts.OnZeroModuloValue,
	}

	processor := post_processors.PostProcessor{
		Sort:          opts.OutputSort,
		Order:         opts.OutputOrder,
		Normalization: opts.OutputNormalized,
	}

	isFirst := true
	var input string
	for {

		input, err = inputReader.ReadOne()

		if err == io.EOF {
			return 0
		}

		if err != nil {
			PrintError(err)
			break
		}

		if input == "" {
			continue
		}

		if isFirst {
			isFirst = false
		} else {
			print(opts.ResultsSeparator)
		}

		err = ProcessOne(
			input,
			engine,
			processor,
			opts.OutputFormatter,
			opts.OutputFormatterOptions,
			opts.TimeIt,
		)

		if err != nil {
			PrintError(err)
			break
		}

	}

	return 1
}

func ProcessOne(
	input string,
	engine interpreter.InterpreterEngine,
	processor post_processors.PostProcessor,
	formatter formatters.FormatterFunc,
	formatterOptions formatters.FormatterOptions,
	timeIt bool,
) error {
	var timeStart, timeLex, timeParse, timeProcess, timePostProcess, timeFormat time.Time

	timeStart = time.Now()
	tokens, err := core.Lex(input)
	if err != nil {
		return err
	}
	timeLex = time.Now()

	ast, err := core.Parse[interpreter.DataRecord](tokens)
	if err != nil {
		return err
	}
	timeParse = time.Now()

	result, err := ast.Process(engine)
	if err != nil {
		return err
	}
	timeProcess = time.Now()

	resultSequence := processor.Process(result)
	timePostProcess = time.Now()

	out, err := formatter(formatterOptions, resultSequence)
	timeFormat = time.Now()

	if err != nil {
		return err
	}

	fmt.Println(out)
	if timeIt {
		fmt.Printf(`Lexing  	%.3f ms
Parsing  	%.3f ms
Evaluation	%.3f ms
Postprocess	%.3f ms
Formatting	%.3f ms
Total		%.3f ms
`,
			float64(timeLex.Sub(timeStart).Nanoseconds())/1000000,
			float64(timeParse.Sub(timeLex).Nanoseconds())/1000000,
			float64(timeProcess.Sub(timeParse).Nanoseconds())/1000000,
			float64(timePostProcess.Sub(timeProcess).Nanoseconds())/1000000,
			float64(timeFormat.Sub(timePostProcess).Nanoseconds())/1000000,
			float64(timeFormat.Sub(timeStart).Nanoseconds())/1000000,
		)
	}

	return nil
}
