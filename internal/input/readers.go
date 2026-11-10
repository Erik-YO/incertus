package input

import (
	"fmt"
	"io"
	"os"
	"strings"
)

func CreateReader(
	source string,
	fromFile bool,
	singleExpressions bool,
	expressionSeparator string,
) (InputReader, error) {

	if !fromFile {
		return CreateInputReaderStatic(source, singleExpressions, expressionSeparator), nil
	}

	if source == "-" {
		return CreateInputReaderStdin(source, singleExpressions, expressionSeparator)
	}

	return CreateInputReaderFile(source, singleExpressions, expressionSeparator)
}

func CreateInputReaderStatic(
	source string,
	singleExpression bool,
	expressionSeparator string,
) *InputReaderStatic {

	var lines []string
	if singleExpression {
		lines = []string{source}
	} else {
		lines = strings.Split(source, expressionSeparator)
	}

	return NewInputReaderStatic(lines)
}

func CreateInputReaderStdin(
	source string,
	singleExpression bool,
	expressionSeparator string,
) (InputReader, error) {

	if singleExpression {
		contentBytes, err := io.ReadAll(os.Stdin)
		if err != nil {
			return nil, fmt.Errorf("Error reading from stdin")
		}
		content := string(contentBytes)
		return NewInputReaderStatic([]string{content}), nil
	}

	return NewInputReaderFileObj(os.Stdin, expressionSeparator), nil
}

func CreateInputReaderFile(
	source string,
	singleExpression bool,
	expressionSeparator string,
) (InputReader, error) {

	if singleExpression {
		contentBytes, err := os.ReadFile(source)
		if err != nil {
			return nil, fmt.Errorf("Error opening the file at '%s'", source)
		}
		content := string(contentBytes)
		return NewInputReaderStatic([]string{content}), nil
	}

	file, err := os.Open(source)
	if err != nil {
		return nil, fmt.Errorf("Error opening the file at '%s'", source)
	}

	return NewInputReaderFileObj(file, expressionSeparator), nil
}
