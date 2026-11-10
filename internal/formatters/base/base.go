package base

import (
	"fmt"
	"math"
	"strconv"
)

type FormatterOptions struct {
	StringValues bool
	RoundSet     bool
	Round        int
}

func FormatValue(v int64, opts FormatterOptions) string {
	if opts.StringValues {
		return fmt.Sprintf("\"%d\"", v)
	}
	return fmt.Sprintf("%d", v)
}

func FormatFrequency(f float64, opts FormatterOptions) string {

	if opts.RoundSet {

		if opts.Round == 0 {
			return fmt.Sprintf("%d", uint64(math.Trunc(f)))
		}

		if opts.Round > 0 {
			return strconv.FormatFloat(f, 'f', opts.Round, 64)
		}
		ratio := math.Pow(10, float64(-opts.Round))
		f = math.Trunc(f/ratio) * ratio
		return fmt.Sprintf("%d", uint64(f))
	}

	if f == math.Trunc(f) {
		// No decimal dot
		return fmt.Sprintf("%d", uint64(f))
	}
	// -1 means just enough decimals to represent the number, without trailing zeros
	return strconv.FormatFloat(f, 'f', -1, 64)
}
