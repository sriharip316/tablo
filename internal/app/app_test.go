package app

import (
	"strings"
	"testing"

	"github.com/sriharip316/tablo/internal/flatten"
)

func TestApplication_BasicFlow(t *testing.T) {
	config := Config{
		Input: InputConfig{
			String: `{"name": "test", "value": 42}`,
			Format: "json",
		},
		Output: OutputConfig{
			Style: "ascii",
		},
	}

	app := New(config, nil)
	err := app.Run()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestApplication_ConflictingInputs(t *testing.T) {
	config := Config{
		Input: InputConfig{
			String: `{"test": true}`,
			File:   "test.json",
		},
	}

	app := New(config, nil)
	err := app.Run()

	var appErr *AppError
	if !AsAppError(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}

	if appErr.Code != ErrCodeUsage {
		t.Errorf("expected usage error, got %v", appErr.Code)
	}
}

func TestApplication_ProcessObject(t *testing.T) {
	config := Config{
		Flatten: FlattenConfig{
			Enabled: true,
		},
		Output: OutputConfig{
			Style: "ascii",
		},
	}

	app := New(config, nil)

	obj := map[string]any{
		"a": map[string]any{"b": 1},
		"c": 2,
	}

	flattenOpts := flatten.Options{
		Enabled:  true,
		MaxDepth: -1,
	}

	model, err := app.processObject(obj, flattenOpts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if model.Mode != 1 { // ModeObjectKV
		t.Errorf("expected object KV mode, got %v", model.Mode)
	}
}

func TestApplication_ProcessArray(t *testing.T) {
	config := Config{
		Output: OutputConfig{
			IndexColumn: true,
		},
	}

	app := New(config, nil)

	arr := []any{
		map[string]any{"name": "Alice", "age": 30},
		map[string]any{"name": "Bob", "age": 25},
	}

	flattenOpts := flatten.Options{
		Enabled: false,
	}

	model, err := app.processArray(arr, flattenOpts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(model.Headers) == 0 {
		t.Error("expected headers to be populated")
	}
	if len(model.Rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(model.Rows))
	}
	if !model.IndexColumn {
		t.Error("expected index column to be enabled")
	}
}

func TestApplication_ProcessPrimitiveArray(t *testing.T) {
	config := Config{
		Output: OutputConfig{
			Limit: 2,
		},
	}

	app := New(config, nil)

	arr := []any{1, 2, 3, 4}

	flattenOpts := flatten.Options{
		Enabled: false,
	}

	model, err := app.processArray(arr, flattenOpts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(model.Headers) != 1 || model.Headers[0] != "VALUE" {
		t.Errorf("expected single VALUE header, got %v", model.Headers)
	}
	if len(model.Rows) != 2 {
		t.Errorf("expected limit to be applied, got %d rows", len(model.Rows))
	}
}

func TestApplication_ApplySelection(t *testing.T) {
	config := Config{
		Selection: SelectionConfig{
			SelectExpr: "name,age",
		},
	}

	app := New(config, nil)

	keys := []string{"name", "age", "city", "country"}

	filtered, err := app.applySelection(keys)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(filtered) != 2 {
		t.Errorf("expected 2 filtered keys, got %d", len(filtered))
	}
	if filtered[0] != "name" || filtered[1] != "age" {
		t.Errorf("unexpected filtered keys: %v", filtered)
	}
}

func TestApplication_StrictSelection(t *testing.T) {
	config := Config{
		Selection: SelectionConfig{
			SelectExpr:   "missing",
			StrictSelect: true,
		},
	}

	app := New(config, nil)

	keys := []string{"name", "age"}

	_, err := app.applySelection(keys)

	var appErr *AppError
	if !AsAppError(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}

	if appErr.Code != ErrCodeSelection {
		t.Errorf("expected selection error, got %v", appErr.Code)
	}

	if !strings.Contains(appErr.Message, "missing selected paths") {
		t.Errorf("unexpected error message: %s", appErr.Message)
	}
}

func TestApplication_NormalizeData(t *testing.T) {
	app := &Application{}

	// Test map[any]any normalization
	input := map[any]any{"key": "value", 123: "number"}
	normalized := app.normalizeData(input)

	result, ok := normalized.(map[string]any)
	if !ok {
		t.Fatalf("expected map[string]any, got %T", normalized)
	}

	if result["key"] != "value" {
		t.Errorf("expected value 'value', got %v", result["key"])
	}
	if result["123"] != "number" {
		t.Errorf("expected value 'number', got %v", result["123"])
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test splitCommaString
	result := splitCommaString("a, b,, c ")
	expected := []string{"a", "b", "c"}
	if len(result) != len(expected) {
		t.Errorf("splitCommaString length mismatch: got %d, want %d", len(result), len(expected))
	}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("splitCommaString[%d]: got %s, want %s", i, result[i], v)
		}
	}

	// Test joinStrings
	joined := joinStrings([]string{"a", "b", "c"}, ", ")
	if joined != "a, b, c" {
		t.Errorf("joinStrings: got %s, want 'a, b, c'", joined)
	}

	// Test trimSpace
	trimmed := trimSpace("  hello world  ")
	if trimmed != "hello world" {
		t.Errorf("trimSpace: got '%s', want 'hello world'", trimmed)
	}

	// Test endsWithNewline
	if !endsWithNewline("test\n") {
		t.Error("endsWithNewline: expected true for string ending with newline")
	}
	if endsWithNewline("test") {
		t.Error("endsWithNewline: expected false for string not ending with newline")
	}
}

func TestApplication_WriteOutput(t *testing.T) {
	config := Config{
		Output: OutputConfig{
			FilePath: "", // stdout
		},
	}

	app := New(config, nil)

	// Test that newline is added
	err := app.writeOutput("test output")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test with output that already has newline
	err = app.writeOutput("test output\n")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestConfig_Defaults(t *testing.T) {
	// Test that default values are properly handled
	config := Config{}

	// Should not panic with empty config
	app := New(config, strings.NewReader(`{"test": true}`))
	if app == nil {
		t.Fatal("expected non-nil application")
	}
}

func TestSplitString(t *testing.T) {
	// Test edge cases for splitString helper
	result := splitString("", ",")
	if result != nil {
		t.Errorf("splitString empty string: got %v, want nil", result)
	}

	result = splitString("single", ",")
	if len(result) != 1 || result[0] != "single" {
		t.Errorf("splitString single item: got %v, want [single]", result)
	}

	result = splitString("a,,b", ",")
	if len(result) != 3 || result[1] != "" {
		t.Errorf("splitString with empty: got %v, want [a  b]", result)
	}
}

func TestIsSpace(t *testing.T) {
	testCases := []struct {
		char     byte
		expected bool
	}{
		{' ', true},
		{'\t', true},
		{'\n', true},
		{'\r', true},
		{'a', false},
		{'1', false},
		{'!', false},
	}

	for _, tc := range testCases {
		result := isSpace(tc.char)
		if result != tc.expected {
			t.Errorf("isSpace(%q): got %v, want %v", tc.char, result, tc.expected)
		}
	}
}
