package parse

import (
	"encoding/json"
	"reflect"
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

func TestParseJSON_NumberNormalization(t *testing.T) {
	data := []byte(`{"a": 1.5, "b": [1,2,3]}`)
	v, err := Parse(data, JSON)
	if err != nil {
		t.Fatalf("parse json: %v", err)
	}
	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("type: %T", v)
	}
	if _, ok := m["a"].(json.Number); !ok {
		t.Fatalf("a should be json.Number, got %T", m["a"])
	}
}

func TestArrayIsObjects(t *testing.T) {
	if ArrayIsObjects([]any{map[string]any{"a": 1}, map[string]any{"b": 2}}) != true {
		t.Fatal("expected true")
	}
	if ArrayIsObjects([]any{1, map[string]any{"a": 1}}) != false {
		t.Fatal("expected false")
	}
}

func TestToStringKeyMap(t *testing.T) {
	in := map[any]any{"x": 1, 2: "y"}
	out := ToStringKeyMap(in)
	expKeys := []string{"x", "2"}
	for _, k := range expKeys {
		if _, ok := out[k]; !ok {
			t.Fatalf("missing key %s", k)
		}
	}
}

func TestNormalizeNested(t *testing.T) {
	in := map[any]any{"m": map[any]any{"k": 1}, "arr": []any{map[any]any{"z": 2}}}
	v := normalize(in)
	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("type %T", v)
	}
	if !reflect.DeepEqual(m["m"], map[string]any{"k": 1}) {
		t.Fatalf("nested map not normalized: %#v", m["m"])
	}
}
