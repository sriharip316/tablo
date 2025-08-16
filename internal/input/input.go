package input

import (
	"errors"
	"io"
	"os"
)

type Reader struct {
	inStr string
	file  string
	stdin io.Reader
}

func NewReader(inStr, file string, stdin io.Reader) *Reader {
	return &Reader{inStr: inStr, file: file, stdin: stdin}
}

func (r *Reader) Read() ([]byte, error) {
	if r.inStr != "" {
		return []byte(r.inStr), nil
	}
	if r.file != "" {
		b, err := os.ReadFile(r.file)
		if err != nil {
			return nil, err
		}
		return b, nil
	}
	if r.stdin == nil {
		return nil, errors.New("no input available")
	}
	return io.ReadAll(r.stdin)
}
