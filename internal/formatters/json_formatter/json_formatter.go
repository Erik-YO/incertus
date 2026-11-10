package json_formatter

import (
	"incertus/internal/common"
	"incertus/internal/formatters/base"
	"strings"
)

func JsonLineFormat(f base.FormatterOptions, s common.ResultSequence) (string, error) {
	builder := strings.Builder{}

	builder.WriteRune('{')
	printColon := false
	for value, freq := range s {
		if printColon {
			builder.WriteRune(',')
			builder.WriteRune(' ')
		} else {
			printColon = true
		}

		builder.WriteString(base.FormatValue(value, f))
		builder.WriteRune(':')
		builder.WriteRune(' ')
		builder.WriteString(base.FormatFrequency(freq, f))
	}
	builder.WriteRune('}')
	builder.WriteRune('\n')

	return builder.String(), nil
}
