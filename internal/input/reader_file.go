package input

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

type InputReaderFileObj struct {
	file      *bufio.Reader
	separator []byte
	buff      []byte
}

func (r *InputReaderFileObj) ReadOne() (string, error) {
	sepLen := len(r.separator)
	for {
		b, err := r.file.ReadByte()
		if err != nil {
			if len(r.buff) < 1 {
				break
			}
			value := string(r.buff)
			r.buff = nil // Free buffer memory after last read
			return value, nil
		}

		r.buff = append(r.buff, b)

		buffLen := len(r.buff)
		if buffLen < sepLen {
			continue
		}

		if bytes.Equal(r.separator, r.buff[buffLen-sepLen:]) {
			value := string(r.buff[:buffLen-sepLen])
			r.buff = r.buff[:0] // keep the buffer memory but reset content
			return value, nil
		}
	}

	return "", io.EOF
}

func NewInputReaderFileObj(
	file *os.File,
	expressionSeparator string,
) *InputReaderFileObj {
	separator := []byte(expressionSeparator)
	return &InputReaderFileObj{
		file:      bufio.NewReader(file),
		separator: separator,
		buff:      nil,
	}
}
