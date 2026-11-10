package formatters

import (
	"incertus/internal/formatters/csv_formatter"
	"incertus/internal/formatters/graph_formatter"
	"incertus/internal/formatters/json_formatter"
	"incertus/internal/formatters/text_formatter"
)

var DefaultFormatterEntry formatterEntry = formatterEntry{
	Function:                text_formatter.TextFormat,
	DefaultResultsSeparator: "\n---\n",
	Description:             "As text, `value:\tfrequency` (the whitespace is a tab)",
}

var FormatterRegistry map[string]formatterEntry = map[string]formatterEntry{
	"text": DefaultFormatterEntry,
	"json": {
		Function:                json_formatter.JsonLineFormat,
		DefaultResultsSeparator: "\n",
		Description:             "As a JSON object in a single line",
	},
	"csv": {
		Function:                csv_formatter.CsvFormat,
		DefaultResultsSeparator: "\n---\n",
		Description:             "Comma separated values, `value,frequency`",
	},
	"graph": {
		Function:                graph_formatter.GraphFormat,
		DefaultResultsSeparator: "\n---\n",
		Description:             "A visual representation of the relative frequency of each value",
	},
}
