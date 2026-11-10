package graph_formatter

import (
	"incertus/internal/common"
	"incertus/internal/formatters/base"
	"math"
	"strings"
)

const MaxWidth = 50

func GraphFormat(f base.FormatterOptions, s common.ResultSequence) (string, error) {
	builder := strings.Builder{}

	var a []struct {
		v int64
		f float64
	}
	maxF := 0.
	for v, f := range s {
		a = append(a, struct {
			v int64
			f float64
		}{v, f})
		if f > maxF {
			maxF = f
		}
	}

	for _, vf := range a {
		value, freq := vf.v, vf.f

		builder.WriteString(base.FormatValue(value, f))
		builder.WriteRune('\t')

		builder.WriteString(base.FormatFrequency(freq, f))

		if maxF > 0 {
			builder.WriteRune('\t')
			w := int64(math.Round((freq * MaxWidth) / maxF))
			for range w {
				builder.WriteRune('█')
			}
		}

		builder.WriteRune('\n')
	}

	return builder.String(), nil
}
