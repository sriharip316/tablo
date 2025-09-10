package input

import (
	"errors"
	"os"
	"strings"
	"testing"
)

func TestReader_InStr(t *testing.T) {
	r := NewReader("hello", "", nil)
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
	r := NewReader("", f.Name(), nil)
	b, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if string(b) != "filedata" {
		t.Fatalf("got %q", string(b))
	}
}

func TestReader_Stdin(t *testing.T) {
	r := NewReader("", "", strings.NewReader("stdin"))
	b, err := r.Read()
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if string(b) != "stdin" {
		t.Fatalf("got %q", string(b))
	}
}

func TestReader_NoInputError(t *testing.T) {
	r := NewReader("", "", nil)
	_, err := r.Read()
	if err == nil {
		t.Fatalf("expected error")
	}
	if !errors.Is(err, err) { // ensure it's non-nil and propagated
		t.Logf("err: %v", err)
	}
}
