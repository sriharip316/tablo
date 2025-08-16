package parse

import (
	"strings"
	"testing"
)

func TestDetect(t *testing.T) {
	data := []byte("{\"a\":1}")
	d := Detector{Explicit: "auto", FilePath: ""}
	if got := d.Detect(data); got != JSON {
		t.Fatalf("detect json: got %v", got)
	}
	d = Detector{Explicit: "yaml"}
	if got := d.Detect([]byte("a: 1")); got != YAML {
		t.Fatalf("detect yaml explicit: got %v", got)
	}
}

func TestParseYAMLMultiDoc(t *testing.T) {
	yaml := strings.TrimSpace(`
---
a: 1
---
b: 2
`)
	v, err := Parse([]byte(yaml), YAML)
	if err != nil {
		t.Fatalf("parse yaml: %v", err)
	}
	arr, ok := v.([]any)
	if !ok || len(arr) != 2 {
		t.Fatalf("want 2 docs, got %T len=%d", v, len(arr))
	}
}
