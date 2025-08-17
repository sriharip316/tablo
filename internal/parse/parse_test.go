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

func TestDetect_FileExtensionCaseInsensitive(t *testing.T) {
	data := []byte("a: 1")
	d := Detector{FilePath: "DATA.YML", Explicit: ""}
	if got := d.Detect(data); got != YAML {
		t.Fatalf("want YAML got %v", got)
	}
}

func TestDetect_ExplicitInvalidFallsBackToSniff(t *testing.T) {
	data := []byte(" { \"a\": 1 } ")
	d := Detector{Explicit: "weird", FilePath: ""}
	if got := d.Detect(data); got != JSON {
		t.Fatalf("want JSON sniff got %v", got)
	}
}

func TestDetect_SniffJSONWithLeadingSpaces(t *testing.T) {
	data := []byte("   [1,2,3]")
	d := Detector{}
	if got := d.Detect(data); got != JSON {
		t.Fatalf("want JSON got %v", got)
	}
}

func TestParse_InvalidJSON(t *testing.T) {
	bad := []byte("{\"a\": 1")
	if _, err := Parse(bad, JSON); err == nil {
		t.Fatal("expected JSON parse error")
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	bad := []byte("a: [1,2\n")
	if _, err := Parse(bad, YAML); err == nil {
		t.Fatal("expected YAML parse error")
	}
}

func TestParse_YAMLSkipEmptyDocs(t *testing.T) {
	yaml := strings.TrimSpace(`
---
# comment only
---
a: 1
---
`)
	v, err := Parse([]byte(yaml), YAML)
	if err != nil {
		t.Fatalf("parse yaml: %v", err)
	}
	if _, ok := v.(map[string]any); !ok {
		t.Fatalf("expected single map, got %T", v)
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	if _, err := Parse([]byte("a: 1"), Format("bogus")); err == nil {
		t.Fatal("expected invalid format error")
	} else if err != ErrInvalidFormat {
		t.Fatalf("expected ErrInvalidFormat got %v", err)
	}
}

func TestParse_YAMLMultiDocWithNilDocs(t *testing.T) {
	yaml := strings.TrimSpace(`
---

---

---
b: 2
`)
	v, err := Parse([]byte(yaml), YAML)
	if err != nil {
		t.Fatalf("parse yaml: %v", err)
	}
	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("expected single map got %T", v)
	}
	if m["b"] != 2 {
		t.Fatalf("unexpected doc content: %#v", m)
	}
}
