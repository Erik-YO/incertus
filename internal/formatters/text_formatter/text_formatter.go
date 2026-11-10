package text_formatter

import (
	"incertus/internal/common"
	"incertus/internal/formatters/base"
	"strings"
)

func TextFormat(f base.FormatterOptions, s common.ResultSequence) (string, error) {
	builder := strings.Builder{}

	for value, freq := range s {

		builder.WriteString(base.FormatValue(value, f))
		builder.WriteRune(':')
		builder.WriteRune('\t')

		builder.WriteString(base.FormatFrequency(freq, f))
		builder.WriteRune('\n')
	}

	return builder.String(), nil
}
