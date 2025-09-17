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
	v, err := Parse([]byte(yaml), YAML, ParseOptions{})
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
	v, err := Parse(data, JSON, ParseOptions{})
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
	if _, err := Parse(bad, JSON, ParseOptions{}); err == nil {
		t.Fatal("expected JSON parse error")
	}
}

func TestParse_InvalidYAML(t *testing.T) {
	bad := []byte("a: [1,2\n")
	if _, err := Parse(bad, YAML, ParseOptions{}); err == nil {
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
	v, err := Parse([]byte(yaml), YAML, ParseOptions{})
	if err != nil {
		t.Fatalf("parse yaml: %v", err)
	}
	if _, ok := v.(map[string]any); !ok {
		t.Fatalf("expected single map, got %T", v)
	}
}

func TestParse_InvalidFormat(t *testing.T) {
	if _, err := Parse([]byte("a: 1"), Format("bogus"), ParseOptions{}); err == nil {
		t.Fatal("expected invalid format error")
	} else if err != ErrInvalidFormat {
		t.Fatalf("expected ErrInvalidFormat got %v", err)
	}
}

func TestParse_CSV(t *testing.T) {
	csvData := []byte(`name,age,city
John,30,NYC
Jane,25,LA`)
	v, err := Parse(csvData, CSV, ParseOptions{})
	if err != nil {
		t.Fatalf("parse CSV: %v", err)
	}
	arr, ok := v.([]map[string]any)
	if !ok || len(arr) != 2 {
		t.Fatalf("expected 2 objects, got %T", v)
	}
	if arr[0]["name"] != "John" {
		t.Fatalf("unexpected name: %v", arr[0]["name"])
	}
	if arr[0]["age"] != "30" {
		t.Fatalf("unexpected age: %v", arr[0]["age"])
	}
	if arr[1]["city"] != "LA" {
		t.Fatalf("unexpected city: %v", arr[1]["city"])
	}
}

func TestParse_CSVNoHeader(t *testing.T) {
	csvData := []byte(`John,30,NYC
Jane,25,LA`)
	v, err := Parse(csvData, CSV, ParseOptions{CSVNoHeader: true})
	if err != nil {
		t.Fatalf("parse CSV no header: %v", err)
	}
	arr, ok := v.([]map[string]any)
	if !ok || len(arr) != 2 {
		t.Fatalf("expected 2 objects, got %T", v)
	}
	if arr[0]["col0"] != "John" {
		t.Fatalf("unexpected col0: %v", arr[0]["col0"])
	}
	if arr[0]["col1"] != "30" {
		t.Fatalf("unexpected col1: %v", arr[0]["col1"])
	}
	if arr[1]["col2"] != "LA" {
		t.Fatalf("unexpected col2: %v", arr[1]["col2"])
	}
}

func TestDetect_CSVExtension(t *testing.T) {
	data := []byte("name,age\nJohn,30")
	d := Detector{FilePath: "data.csv", Explicit: ""}
	if got := d.Detect(data); got != CSV {
		t.Fatalf("want CSV got %v", got)
	}
}

func TestDetect_CSVSniff(t *testing.T) {
	data := []byte("name,age,city\nJohn,30,NYC")
	d := Detector{Explicit: "auto"}
	if got := d.Detect(data); got != CSV {
		t.Fatalf("want CSV got %v", got)
	}
}

func TestDetect_CSVExplicit(t *testing.T) {
	data := []byte("name,age\nJohn,30")
	d := Detector{Explicit: "csv"}
	if got := d.Detect(data); got != CSV {
		t.Fatalf("want CSV got %v", got)
	}
}

func TestParse_YAMLMultiDocWithNilDocs(t *testing.T) {
	yaml := strings.TrimSpace(`
---

---

---
b: 2
`)
	v, err := Parse([]byte(yaml), YAML, ParseOptions{})
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

func TestParse_JSONWithComments(t *testing.T) {
	// Test JSON with line comments
	jsonWithComments := []byte(`{
		// This is a line comment
		"name": "test", // End of line comment
		"age": 30,
		/* This is a block comment */
		"active": true,
		/*
		 * Multi-line block comment
		 * with multiple lines
		 */
		"data": {
			"nested": "value" // Nested comment
		}
	}`)

	v, err := Parse(jsonWithComments, JSON, ParseOptions{})
	if err != nil {
		t.Fatalf("parse JSON with comments: %v", err)
	}

	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", v)
	}

	if m["name"] != "test" {
		t.Fatalf("expected name=test, got %v", m["name"])
	}
	if !reflect.DeepEqual(m["data"], map[string]any{"nested": "value"}) {
		t.Fatalf("nested data incorrect: %#v", m["data"])
	}
}

func TestParse_JSONArrayWithComments(t *testing.T) {
	jsonArray := []byte(`[
		// First element
		{
			"id": 1,
			"name": "first" // Name comment
		},
		/* Second element */
		{
			"id": 2,
			"name": "second"
		}
		// End of array
	]`)

	v, err := Parse(jsonArray, JSON, ParseOptions{})
	if err != nil {
		t.Fatalf("parse JSON array with comments: %v", err)
	}

	arr, ok := v.([]any)
	if !ok {
		t.Fatalf("expected array, got %T", v)
	}

	if len(arr) != 2 {
		t.Fatalf("expected 2 elements, got %d", len(arr))
	}
}

func TestParse_YAMLWithComments(t *testing.T) {
	// YAML already supports comments natively, but let's test to ensure it still works
	yamlWithComments := []byte(`# Top-level comment
name: test  # End of line comment
age: 30
# Another comment
data:
  # Nested comment
  nested: value
  # Multi-line
  # comment block
  other: data
`)

	v, err := Parse(yamlWithComments, YAML, ParseOptions{})
	if err != nil {
		t.Fatalf("parse YAML with comments: %v", err)
	}

	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", v)
	}

	if m["name"] != "test" {
		t.Fatalf("expected name=test, got %v", m["name"])
	}
	if !reflect.DeepEqual(m["data"], map[string]any{"nested": "value", "other": "data"}) {
		t.Fatalf("nested data incorrect: %#v", m["data"])
	}
}

func TestParse_JSONCommentsInAutoDetect(t *testing.T) {
	// Test that JSON with comments is properly detected and parsed in auto mode
	jsonWithComments := []byte(`{
		// Auto-detect test
		"format": "json",
		"comments": true
	}`)

	// Should be detected as JSON due to opening brace
	d := Detector{Explicit: "auto"}
	format := d.Detect(jsonWithComments)
	if format != JSON {
		t.Fatalf("expected JSON detection, got %v", format)
	}

	v, err := Parse(jsonWithComments, format, ParseOptions{})
	if err != nil {
		t.Fatalf("parse auto-detected JSON with comments: %v", err)
	}

	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("expected map, got %T", v)
	}

	if m["format"] != "json" {
		t.Fatalf("expected format=json, got %v", m["format"])
	}
}

func TestDetect_JSONCExtension(t *testing.T) {
	// Test .jsonc extension is detected as JSON
	data := []byte(`{"test": "value"}`)
	d := Detector{FilePath: "test.jsonc", Explicit: ""}
	if got := d.Detect(data); got != JSON {
		t.Fatalf("want JSON for .jsonc extension, got %v", got)
	}

	// Test uppercase .JSONC extension
	d = Detector{FilePath: "test.JSONC", Explicit: ""}
	if got := d.Detect(data); got != JSON {
		t.Fatalf("want JSON for .JSONC extension, got %v", got)
	}
}

func TestDetect_LeadingCommentJSON(t *testing.T) {
	// Test JSON with leading line comment
	jsonWithLeadingComment := []byte(`// Top-level comment
{
  "name": "test",
  "value": 42
}`)
	d := Detector{}
	if got := d.Detect(jsonWithLeadingComment); got != JSON {
		t.Fatalf("want JSON for data with leading comment, got %v", got)
	}

	// Test JSON with leading block comment
	jsonWithBlockComment := []byte(`/* Block comment
   spanning multiple lines */
{"name": "test"}`)
	if got := d.Detect(jsonWithBlockComment); got != JSON {
		t.Fatalf("want JSON for data with leading block comment, got %v", got)
	}

	// Test JSON array with leading comment
	jsonArrayWithComment := []byte(`// Comment before array
[1, 2, 3]`)
	if got := d.Detect(jsonArrayWithComment); got != JSON {
		t.Fatalf("want JSON for array with leading comment, got %v", got)
	}
}

func TestDetect_JSONL(t *testing.T) {
	data := []byte(`{"name": "Alice"}
{"name": "Bob"}`)
	d := Detector{}
	if got := d.Detect(data); got != JSONL {
		t.Fatalf("want JSONL, got %v", got)
	}
}

func TestParse_JSONL(t *testing.T) {
	data := []byte(`{"id": 1}
{"id": 2}`)
	v, err := Parse(data, JSONL, ParseOptions{})
	if err != nil {
		t.Fatalf("parse JSONL: %v", err)
	}
	arr, ok := v.([]any)
	if !ok || len(arr) != 2 {
		t.Fatalf("expected 2 objects, got %T", v)
	}
}

func TestParse_JSONLWithEmptyLines(t *testing.T) {
	data := []byte(`{"id": 1}

{"id": 2}`)
	v, err := Parse(data, JSONL, ParseOptions{})
	if err != nil {
		t.Fatalf("parse JSONL with empty lines: %v", err)
	}
	arr, ok := v.([]any)
	if !ok || len(arr) != 2 {
		t.Fatalf("expected 2 objects, got %T", v)
	}
}

func TestParse_JSONLInvalid(t *testing.T) {
	data := []byte(`{"id": 1}
invalid json`)
	_, err := Parse(data, JSONL, ParseOptions{})
	if err == nil {
		t.Fatal("expected parse error for invalid JSONL")
	}
}

func TestDetect_JSONLExplicit(t *testing.T) {
	data := []byte(`{"name": "Alice"}
{"name": "Bob"}`)
	d := Detector{Explicit: "jsonl"}
	if got := d.Detect(data); got != JSONL {
		t.Fatalf("want JSONL, got %v", got)
	}
}

func TestParse_JSONLWithArrays(t *testing.T) {
	data := []byte(`[{"a":1}]
[{"a":2},{"a":3}]`)
	v, err := Parse(data, JSONL, ParseOptions{})
	if err != nil {
		t.Fatalf("parse JSONL with arrays: %v", err)
	}
	arr, ok := v.([]any)
	if !ok || len(arr) != 3 {
		t.Fatalf("expected 3 objects, got %d", len(arr))
	}
	// Check that all elements are properly flattened
	for i, item := range arr {
		obj, ok := item.(map[string]any)
		if !ok {
			t.Fatalf("item %d should be object, got %T", i, item)
		}
		expected := i + 1
		if obj["a"] != float64(expected) {
			t.Fatalf("item %d should have a=%d, got %v", i, expected, obj["a"])
		}
	}
}
