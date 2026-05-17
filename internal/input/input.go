package input

import (
	"errors"
	"io"
	"os"
)

var ErrLimitExceeded = errors.New("input size exceeds limit")

type Reader struct {
	inStr    string
	file     string
	stdin    io.Reader
	maxBytes int64
}

func NewReader(inStr, file string, stdin io.Reader, maxBytes int64) *Reader {
	return &Reader{
		inStr:    inStr,
		file:     file,
		stdin:    stdin,
		maxBytes: maxBytes,
	}
}

func (r *Reader) Read() ([]byte, error) {
	if r.inStr != "" {
		if r.maxBytes > 0 && int64(len(r.inStr)) > r.maxBytes {
			return nil, ErrLimitExceeded
		}
		return []byte(r.inStr), nil
	}

	if r.file != "" {
		f, err := os.Open(r.file)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		if r.maxBytes > 0 {
			info, err := f.Stat()
			if err == nil && info.Size() > r.maxBytes {
				return nil, ErrLimitExceeded
			}
		}
		return r.readLimited(f)
	}

	if r.stdin == nil {
		return nil, errors.New("no input available")
	}
	return r.readLimited(r.stdin)
}

func (r *Reader) readLimited(reader io.Reader) ([]byte, error) {
	if r.maxBytes <= 0 {
		return io.ReadAll(reader)
	}

	// Read up to maxBytes + 1 to detect if limit is exceeded
	lr := io.LimitReader(reader, r.maxBytes+1)
	data, err := io.ReadAll(lr)
	if err != nil {
		return nil, err
	}

	if int64(len(data)) > r.maxBytes {
		return nil, ErrLimitExceeded
	}

	return data, nil
}
