package formatters

import (
	"incertus/internal/common"
	"incertus/internal/formatters/base"
)

type FormatterOptions = base.FormatterOptions

type FormatterFunc func(f base.FormatterOptions, s common.ResultSequence) (string, error)

type formatterEntry struct {
	Function                FormatterFunc
	DefaultResultsSeparator string
	Description             string
}
