package csv_formatter

import (
	"incertus/internal/common"
	"incertus/internal/formatters/base"
	"strings"
)

func CsvFormat(f base.FormatterOptions, s common.ResultSequence) (string, error) {
	builder := strings.Builder{}

	builder.WriteString("value,frequency\n")

	for value, freq := range s {

		builder.WriteString(base.FormatValue(value, f))
		builder.WriteRune(',')

		builder.WriteString(base.FormatFrequency(freq, f))
		builder.WriteRune('\n')
	}

	return builder.String(), nil
}
