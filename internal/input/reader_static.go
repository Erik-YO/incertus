package input

import (
	"io"
)

type InputReaderStatic struct {
	Lines []string
	next  int
}

func (r *InputReaderStatic) ReadOne() (string, error) {
	if r.next >= len(r.Lines) {
		return "", io.EOF
	}
	value := r.Lines[r.next]
	r.next++
	return value, nil
}

func NewInputReaderStatic(lines []string) *InputReaderStatic {
	return &InputReaderStatic{Lines: lines, next: 0}
}
