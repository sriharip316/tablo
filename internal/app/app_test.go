package app

import (
	"os"
	"path/filepath"
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

func TestNormalizeData_SliceOfAny(t *testing.T) {
	app := &Application{}

	// Test []any normalization
	input := []any{map[any]any{"key": "value"}, "foo", []map[string]any{{"bar": "baz"}}}
	normalized := app.normalizeData(input)

	result, ok := normalized.([]any)
	if !ok {
		t.Fatalf("expected []any, got %T", normalized)
	}

	if len(result) != 3 {
		t.Fatalf("expected 3 elements, got %d", len(result))
	}

	first, ok := result[0].(map[string]any)
	if !ok {
		t.Fatalf("expected first element to be map[string]any, got %T", result[0])
	}

	if first["key"] != "value" {
		t.Errorf("expected value 'value', got %v", first["key"])
	}

	if result[1] != "foo" {
		t.Errorf("expected second element to be 'foo', got %v", result[1])
	}

	third, ok := result[2].([]any)
	if !ok {
		t.Fatalf("expected third element to be []any, got %T", result[2])
	}

	fourth, ok := third[0].(map[string]any)
	if !ok {
		t.Fatalf("expected fourth element to be map[string]any, got %T", third[0])
	}

	if fourth["bar"] != "baz" {
		t.Errorf("expected value 'baz', got %v", fourth["bar"])
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

	joined = joinStrings([]string{}, ", ")
	if joined != "" {
		t.Errorf("joinStrings empty slice: got %s, want ''", joined)
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

func TestRun_NoInput_Error(t *testing.T) {
	// No input string, no file, and nil stdin should error at readInput()
	app := New(Config{}, nil)
	err := app.Run()
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var ae *AppError
	if !AsAppError(err, &ae) || ae.Code != ErrCodeInput {
		t.Fatalf("expected input error, got: %#v", err)
	}
	if !strings.Contains(ae.Message, "failed to read input") {
		t.Errorf("unexpected message: %s", ae.Message)
	}
}

func TestRun_ParseError_JSON(t *testing.T) {
	cfg := Config{Input: InputConfig{String: "{\"a\": }", Format: "json"}}
	app := New(cfg, nil)
	err := app.Run()
	var ae *AppError
	if !AsAppError(err, &ae) || ae.Code != ErrCodeParse {
		t.Fatalf("expected parse error, got: %#v", err)
	}
}

func TestRun_ProcessError_StrictSelectMissing(t *testing.T) {
	cfg := Config{
		Input:     InputConfig{String: `{"name":"n"}`, Format: "json"},
		Selection: SelectionConfig{SelectExpr: "missing", StrictSelect: true},
		Output:    OutputConfig{Style: "ascii"},
	}
	app := New(cfg, nil)
	err := app.Run()
	var ae *AppError
	if !AsAppError(err, &ae) || ae.Code != ErrCodeProcessing {
		t.Fatalf("expected processing error, got: %#v", err)
	}
	// Ensure original selection error message bubbles in chain for visibility
	if ae.Cause == nil {
		t.Fatalf("expected wrapped cause error, got nil")
	}
}

func TestRun_WriteOutput_Error(t *testing.T) {
	// Use a path whose parent directory does not exist -> os.Create should fail with ENOENT
	impossiblePath := filepath.Join("/this/path/should/not/exist", "tablo_out.txt")
	cfg := Config{
		Input:  InputConfig{String: `{"k":1}`, Format: "json"},
		Output: OutputConfig{FilePath: impossiblePath},
	}
	app := New(cfg, nil)
	err := app.Run()
	var ae *AppError
	if !AsAppError(err, &ae) || ae.Code != ErrCodeOutput {
		t.Fatalf("expected output error, got: %#v", err)
	}
}

func TestProcessArray_LimitOne_AsObject(t *testing.T) {
	app := New(Config{Output: OutputConfig{Limit: 1}}, nil)
	arr := []any{
		map[string]any{"a": 1, "b": 2},
		map[string]any{"a": 3, "c": 4},
	}
	model, err := app.processArray(arr, flatten.Options{Enabled: false})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if model.Mode != 1 { // ModeObjectKV
		t.Fatalf("expected object KV mode, got %v", model.Mode)
	}
}

func TestProcessData_Primitive(t *testing.T) {
	app := New(Config{}, nil)
	m, err := app.processData(123)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if m.Mode != 0 { // ModeRows
		t.Fatalf("expected ModeRows for primitive, got %v", m.Mode)
	}
	if len(m.Rows) != 1 || len(m.Rows[0]) != 1 || m.Rows[0][0] != 123 {
		t.Fatalf("unexpected model rows: %+v", m.Rows)
	}
}

func TestCompileSelectors_WithFileAndExclude(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "selectors.txt")
	content := "# comment\n\n name \n age \n# another\naddress.*\n"
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp selectors file: %v", err)
	}

	cfg := Config{
		Selection: SelectionConfig{
			SelectExpr:  "city",
			SelectFile:  path,
			ExcludeExpr: "age",
		},
	}
	app := New(cfg, nil)

	keys := []string{"name", "age", "city", "address.street"}
	filtered, err := app.applySelection(keys)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	joined := strings.Join(filtered, ",")
	// Expect order: includes in order of include-exprs across input order.
	// From file: name, age, address.* (age excluded). From expr: city.
	// So expected roughly: name, city, address.street
	if !strings.Contains(joined, "name") || !strings.Contains(joined, "city") || !strings.Contains(joined, "address.street") {
		t.Fatalf("unexpected filtered keys: %v", filtered)
	}
	for _, k := range filtered {
		if k == "age" {
			t.Fatalf("exclude did not apply; got %v", filtered)
		}
	}
}

func TestApplication_ApplyRowFiltering(t *testing.T) {
	tests := []struct {
		name        string
		whereExprs  []string
		inputData   string
		wantCount   int
		wantErr     bool
		errContains string
	}{
		{
			name:       "no filters",
			whereExprs: []string{},
			inputData:  `[{"name":"John","age":30},{"name":"Jane","age":25}]`,
			wantCount:  2,
		},
		{
			name:       "equality filter",
			whereExprs: []string{"name=John"},
			inputData:  `[{"name":"John","age":30},{"name":"Jane","age":25}]`,
			wantCount:  1,
		},
		{
			name:       "numeric filter",
			whereExprs: []string{"age>25"},
			inputData:  `[{"name":"John","age":30},{"name":"Jane","age":25}]`,
			wantCount:  1,
		},
		{
			name:       "multiple filters (AND)",
			whereExprs: []string{"age>=25", "name!=John"},
			inputData:  `[{"name":"John","age":30},{"name":"Jane","age":25},{"name":"Bob","age":20}]`,
			wantCount:  1,
		},
		{
			name:        "invalid filter",
			whereExprs:  []string{"invalid_expr"},
			inputData:   `[{"name":"John","age":30}]`,
			wantErr:     true,
			errContains: "invalid filter condition",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := Config{
				Filter: FilterConfig{
					WhereExprs: tt.whereExprs,
				},
			}

			application := New(config, strings.NewReader(tt.inputData))
			err := application.Run()

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error but got none")
					return
				}
				if tt.errContains != "" && !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestApplication_FilteringWithFlattening(t *testing.T) {
	input := `[
		{"user":{"name":"John","profile":{"age":30}},"active":true},
		{"user":{"name":"Jane","profile":{"age":25}},"active":false}
	]`

	config := Config{
		Flatten: FlattenConfig{
			Enabled: true,
		},
		Filter: FilterConfig{
			WhereExprs: []string{"user.profile.age>25"},
		},
		Output: OutputConfig{
			Style: "csv",
		},
	}

	application := New(config, strings.NewReader(input))

	// Since we can't easily capture stdout in this test, we'll just verify no error
	err := application.Run()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestApplication_FilteringWithSelection(t *testing.T) {
	input := `[
		{"name":"John","age":30,"city":"NYC","active":true},
		{"name":"Jane","age":25,"city":"LA","active":false}
	]`

	config := Config{
		Selection: SelectionConfig{
			SelectExpr: "name,age",
		},
		Filter: FilterConfig{
			WhereExprs: []string{"active=true"},
		},
		Output: OutputConfig{
			Style: "csv",
		},
	}

	application := New(config, strings.NewReader(input))
	err := application.Run()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestReadSelectFile_Basic(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "sel.txt")
	data := "\n# comment\n name \n \n# c2\nfoo.bar\n"
	if err := os.WriteFile(path, []byte(data), 0o644); err != nil {
		t.Fatalf("failed to write selectors: %v", err)
	}
	app := New(Config{Selection: SelectionConfig{SelectFile: path}}, nil)
	got, err := app.readSelectFile()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 || got[0] != "name" || got[1] != "foo.bar" {
		t.Fatalf("unexpected patterns: %v", got)
	}
}

func TestReadSelectFile_Error(t *testing.T) {
	app := New(Config{Selection: SelectionConfig{SelectFile: filepath.Join("/nope", "missing.txt")}}, nil)
	_, err := app.readSelectFile()
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestWriteOutput_ToFile_AppendsNewline(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "out.txt")
	app := New(Config{Output: OutputConfig{FilePath: path}}, nil)
	if err := app.writeOutput("hello"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read written file: %v", err)
	}
	if string(b) != "hello\n" {
		t.Fatalf("unexpected file content: %q", string(b))
	}
}

func TestSplitLines(t *testing.T) {
	// empty
	if res := splitLines(""); res != nil {
		t.Fatalf("expected nil for empty, got %v", res)
	}
	// multi-line
	res := splitLines("a\nb\nc")
	if len(res) != 3 || res[0] != "a" || res[2] != "c" {
		t.Fatalf("unexpected split result: %v", res)
	}
}

func TestNormalizeData_NoChangeForArrayAndMap(t *testing.T) {
	app := New(Config{}, nil)
	arr := []any{1, 2}
	if got := app.normalizeData(arr); got == nil {
		t.Fatal("unexpected nil")
	} else if _, ok := got.([]any); !ok {
		t.Fatalf("expected []any, got %T", got)
	}
	m := map[string]any{"a": 1}
	if got := app.normalizeData(m); got == nil {
		t.Fatal("unexpected nil")
	} else if _, ok := got.(map[string]any); !ok {
		t.Fatalf("expected map[string]any, got %T", got)
	}
}

func TestApplication_SortingWithPerColumnDirection(t *testing.T) {
	tests := []struct {
		name          string
		sortColumns   []string
		inputRows     []flatten.FlatKV
		expectedOrder []string // names in expected order
	}{
		{
			name:        "sort by name ascending, age descending",
			sortColumns: []string{"name", "-age"},
			inputRows: []flatten.FlatKV{
				{"name": "Alice", "age": 25},
				{"name": "Bob", "age": 30},
				{"name": "Alice", "age": 35},
				{"name": "Bob", "age": 20},
			},
			expectedOrder: []string{"Alice", "Alice", "Bob", "Bob"}, // Alice (35, 25), Bob (30, 20)
		},
		{
			name:        "sort with explicit ascending prefix",
			sortColumns: []string{"+name"},
			inputRows: []flatten.FlatKV{
				{"name": "Charlie"},
				{"name": "Alice"},
				{"name": "Bob"},
			},
			expectedOrder: []string{"Alice", "Bob", "Charlie"},
		},
		{
			name:        "sort with explicit descending prefix",
			sortColumns: []string{"-name"},
			inputRows: []flatten.FlatKV{
				{"name": "Alice"},
				{"name": "Charlie"},
				{"name": "Bob"},
			},
			expectedOrder: []string{"Charlie", "Bob", "Alice"},
		},
		{
			name:        "sort with no prefix defaults to ascending",
			sortColumns: []string{"name"},
			inputRows: []flatten.FlatKV{
				{"name": "Charlie"},
				{"name": "Alice"},
				{"name": "Bob"},
			},
			expectedOrder: []string{"Alice", "Bob", "Charlie"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				config: Config{
					Sort: SortConfig{
						Columns: tt.sortColumns,
					},
				},
			}

			result := app.applySorting(tt.inputRows)

			if len(result) != len(tt.expectedOrder) {
				t.Errorf("expected %d rows, got %d", len(tt.expectedOrder), len(result))
				return
			}

			for i, expectedName := range tt.expectedOrder {
				if result[i]["name"] != expectedName {
					t.Errorf("at position %d: expected name %s, got %s", i, expectedName, result[i]["name"])
				}
			}
		})
	}
}

func TestSplitCommaStringInSorting(t *testing.T) {
	app := &Application{
		config: Config{
			Sort: SortConfig{
				Columns: []string{"name,age", "-department"},
			},
		},
	}

	rows := []flatten.FlatKV{
		{"name": "Bob", "age": 25, "department": "Sales"},
		{"name": "Alice", "age": 30, "department": "Engineering"},
	}

	result := app.applySorting(rows)

	// Should be sorted by: name (asc), age (asc), department (desc)
	// Alice should come before Bob
	if len(result) != 2 {
		t.Errorf("expected 2 rows, got %d", len(result))
		return
	}

	if result[0]["name"] != "Alice" {
		t.Errorf("expected first row name to be Alice, got %s", result[0]["name"])
	}

	if result[1]["name"] != "Bob" {
		t.Errorf("expected second row name to be Bob, got %s", result[1]["name"])
	}
}
