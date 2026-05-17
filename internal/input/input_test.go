package input

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestReader_InStr(t *testing.T) {
	r := NewReader("hello", "", nil, 10)
	b, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if string(b) != "hello" {
		t.Fatalf("got %q", string(b))
	}
}

func TestReader_File(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tablo-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	_, _ = f.WriteString("filedata")
	r := NewReader("", f.Name(), nil, 10)
	b, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if string(b) != "filedata" {
		t.Fatalf("got %q", string(b))
	}
}

func TestReader_Stdin(t *testing.T) {
	r := NewReader("", "", strings.NewReader("stdin"), 10)
	b, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if string(b) != "stdin" {
		t.Fatalf("got %q", string(b))
	}
}

func TestReader_NoInputError(t *testing.T) {
	r := NewReader("", "", nil, 10)
	_, err := r.Read()
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestReader_LimitExceeded_InStr(t *testing.T) {
	r := NewReader("too long", "", nil, 5)
	_, err := r.Read()
	if !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("expected ErrLimitExceeded, got %v", err)
	}
}

func TestReader_LimitExceeded_File(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "tablo-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = f.Close() }()
	_, _ = f.WriteString("too long file data")
	r := NewReader("", f.Name(), nil, 5)
	_, err = r.Read()
	if !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("expected ErrLimitExceeded, got %v", err)
	}
}

func TestReader_LimitExceeded_Stdin(t *testing.T) {
	r := NewReader("", "", strings.NewReader("too long stdin"), 5)
	_, err := r.Read()
	if !errors.Is(err, ErrLimitExceeded) {
		t.Fatalf("expected ErrLimitExceeded, got %v", err)
	}
}
