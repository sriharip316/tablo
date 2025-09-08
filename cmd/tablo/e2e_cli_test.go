package main

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

var cliBin string

func TestMain(m *testing.M) {
	// Build the CLI binary once for E2E tests.
	tmp := os.TempDir()
	cliBin = filepath.Join(tmp, "tablo-test-bin")
	cmd := exec.Command("go", "build", "-o", cliBin, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		// If build fails, exit tests immediately.
		os.Exit(1)
	}
	code := m.Run()
	_ = os.Remove(cliBin)
	os.Exit(code)
}

func runCLI(t *testing.T, args []string, stdin []byte) (string, string, int, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, cliBin, args...)
	var outBuf, errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	err := cmd.Run()
	exit := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exit = ee.ExitCode()
		} else {
			// e.g., context deadline exceeded
			return outBuf.String(), errBuf.String(), -1, err
		}
	}
	return outBuf.String(), errBuf.String(), exit, nil
}

func TestCLI_JSONDive(t *testing.T) {
	// Single object, dive to flatten nested keys.
	args := []string{"-i", `{"a":{"b":1}}`, "--dive", "--style", "markdown"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil {
		t.Fatalf("run err: %v; stderr=%s", err, errOut)
	}
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errOut)
	}
	if !strings.Contains(out, "a.b") || !strings.Contains(out, "1") {
		t.Fatalf("unexpected output: %q", out)
	}
}

func TestCLI_YAMLArray_Select_Limit_Index(t *testing.T) {
	// Array of objects via stdin; include specific headers, limit rows, index column.
	yaml := "- name: Alice\n  age: 30\n- name: Bob\n  age: 31\n"
	args := []string{"-F", "yaml", "--select", "name,age", "--limit", "1", "--index-column", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, []byte(yaml))
	if err != nil {
		t.Fatalf("run err: %v; stderr=%s", err, errOut)
	}
	if code != 0 {
		t.Fatalf("exit=%d stderr=%s", code, errOut)
	}
	// Expect only one data row rendered due to limit and an auto index column
	if !strings.Contains(out, "Alice") || strings.Contains(out, "Bob") {
		t.Fatalf("limit/index failed: %q", out)
	}
}

func TestCLI_PrimitiveArray(t *testing.T) {
	// Array of primitives should render single VALUE column.
	args := []string{"-i", `[1,2,3]`, "--style", "markdown", "--limit", "2"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Two rows due to limit; markdown table should contain 1 and 2, not 3.
	if !strings.Contains(out, "1") || !strings.Contains(out, "2") || strings.Contains(out, "3") {
		t.Fatalf("unexpected rows: %q", out)
	}
}

func TestCLI_OutputFile(t *testing.T) {
	// Write output to a file
	tmpDir := t.TempDir()
	outFile := filepath.Join(tmpDir, "out.txt")
	args := []string{"-i", `{"x":1}`, "--dive", "-o", outFile}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	if out != "" { // should write to file, stdout empty
		t.Fatalf("stdout not empty: %q", out)
	}
	b, rerr := os.ReadFile(outFile)
	if rerr != nil {
		t.Fatalf("read out: %v", rerr)
	}
	if !strings.Contains(string(b), "x") || !strings.HasSuffix(string(b), "\n") {
		t.Fatalf("bad file content: %q", string(b))
	}
}

func TestCLI_StrictSelect_Missing(t *testing.T) {
	args := []string{"-i", `{"a":1}`, "--dive", "--select", "x.*", "--strict-select"}
	_, errOut, code, _ := runCLI(t, args, nil)
	if code == 0 {
		t.Fatalf("expected non-zero exit; stderr=%s", errOut)
	}
	if code != 5 || !strings.Contains(errOut, "missing selected paths") {
		t.Fatalf("want code=5 and missing msg; got code=%d stderr=%q", code, errOut)
	}
}

func TestCLI_ConflictingInputs(t *testing.T) {
	tmp := t.TempDir()
	p := filepath.Join(tmp, "data.json")
	if err := os.WriteFile(p, []byte(`{"a":1}`), 0o600); err != nil {
		t.Fatal(err)
	}
	args := []string{"-i", `{"a":2}`, "-f", p}
	_, errOut, code, _ := runCLI(t, args, nil)
	if code != 2 || !strings.Contains(errOut, "conflicting inputs") {
		t.Fatalf("unexpected: code=%d stderr=%q", code, errOut)
	}
}

func TestCLI_StrictSelect_Success(t *testing.T) {
	// Single object with nested maps; strict select should succeed when all patterns match.
	args := []string{
		"-i", `{"a":{"x":{"c":1},"y":{"c":2}},"b":3,"z":4}`,
		"--dive",
		"--select", "a.*.c,b",
		"--strict-select",
		"--style", "ascii",
	}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Expect flattened keys a.x.c, a.y.c and key b; z should be excluded
	for _, want := range []string{"a.x.c", "a.y.c", "b"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing selected key %s in output: %s", want, out)
		}
	}
	if strings.Contains(out, "z") {
		t.Fatalf("unexpected key z present: %s", out)
	}
}

func TestCLI_PrimitiveArray_IndexColumn(t *testing.T) {
	args := []string{"-i", `[10,20]`, "--index-column", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Expect index numbers 1 and 2 and corresponding values 10 and 20
	if !strings.Contains(out, "| 1 ") || !strings.Contains(out, "| 2 ") {
		t.Fatalf("missing index column values: %s", out)
	}
	if !strings.Contains(out, "10") || !strings.Contains(out, "20") {
		t.Fatalf("missing primitive values: %s", out)
	}
}

func TestCLI_JSONWithComments_Inline(t *testing.T) {
	// Test JSON with comments via inline input
	jsonInput := `{
		"name": "test", // End of line comment
		/* Block comment */
		"value": 42,
		"data": {
			"nested": true // Nested comment
		}
	}`
	args := []string{"-i", jsonInput, "--dive", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Verify all fields are parsed correctly despite comments
	if !strings.Contains(out, "test") || !strings.Contains(out, "42") || !strings.Contains(out, "true") {
		t.Fatalf("JSON with comments not parsed correctly: %s", out)
	}
	if !strings.Contains(out, "data.nested") {
		t.Fatalf("nested data not flattened: %s", out)
	}
}

func TestCLI_JSONWithComments_File(t *testing.T) {
	// Create a temporary file with JSON comments
	tmpDir := t.TempDir()
	jsonFile := filepath.Join(tmpDir, "test.json")
	jsonContent := `{
		// This is a test file
		"project": "tablo",
		"features": [
			{
				"name": "comments", // Feature name
				"enabled": true
			}
		]
		/* End of file */
	}`
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0o600); err != nil {
		t.Fatal(err)
	}

	args := []string{"-f", jsonFile, "--dive", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	if !strings.Contains(out, "tablo") || !strings.Contains(out, "comments") {
		t.Fatalf("JSON file with comments not parsed correctly: %s", out)
	}
}

func TestCLI_YAMLWithComments_Stdin(t *testing.T) {
	// Test YAML with comments via stdin
	yamlInput := `# Top level comment
name: test  # End of line comment
data:
  # Nested comment
  value: 123
  active: true
`
	args := []string{"--format", "yaml", "--dive", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, []byte(yamlInput))
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Verify YAML with comments is parsed correctly
	if !strings.Contains(out, "test") || !strings.Contains(out, "123") || !strings.Contains(out, "true") {
		t.Fatalf("YAML with comments not parsed correctly: %s", out)
	}
	if !strings.Contains(out, "data.value") || !strings.Contains(out, "data.active") {
		t.Fatalf("nested YAML data not flattened: %s", out)
	}
}

func TestCLI_JSONCommentsArray(t *testing.T) {
	// Test JSON array with comments
	jsonArray := `[
		// First item
		{"id": 1, "name": "first"},
		/* Second item */
		{"id": 2, "name": "second"}
		// End of array
	]`
	args := []string{"-i", jsonArray, "--style", "ascii", "--index-column"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should render as table with 2 rows
	if !strings.Contains(out, "first") || !strings.Contains(out, "second") {
		t.Fatalf("JSON array with comments not parsed correctly: %s", out)
	}
	// Should have index column
	if !strings.Contains(out, "| 1 ") || !strings.Contains(out, "| 2 ") {
		t.Fatalf("index column missing: %s", out)
	}
}

func TestCLI_JSONCExtension(t *testing.T) {
	// Test that .jsonc files are properly detected and parsed as JSON
	tmpDir := t.TempDir()
	jsoncFile := filepath.Join(tmpDir, "test.jsonc")
	jsoncContent := `// Top-level comment
{
	"name": "jsonc test",
	/* Block comment */
	"features": [
		{
			"name": "extension support", // Feature comment
			"enabled": true
		}
	]
}`
	if err := os.WriteFile(jsoncFile, []byte(jsoncContent), 0o600); err != nil {
		t.Fatal(err)
	}

	args := []string{"-f", jsoncFile, "--dive", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	if !strings.Contains(out, "jsonc test") || !strings.Contains(out, "extension support") {
		t.Fatalf(".jsonc file not parsed correctly: %s", out)
	}
	if !strings.Contains(out, "features.0.name") || !strings.Contains(out, "features.0.enabled") {
		t.Fatalf(".jsonc file not flattened correctly: %s", out)
	}
}

func TestCLI_HTMLOutput(t *testing.T) {
	jsonInput := `{"name": "John", "age": 30}`
	args := []string{"-i", jsonInput, "--dive", "--style", "html"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Verify HTML output format
	if !strings.Contains(out, "<table") || !strings.Contains(out, "</table>") {
		t.Fatalf("expected HTML table tags, got: %s", out)
	}
	if !strings.Contains(out, "<th>") || !strings.Contains(out, "<td>") {
		t.Fatalf("expected HTML th/td tags, got: %s", out)
	}
}

func TestCLI_CSVOutput(t *testing.T) {
	jsonInput := `[{"name": "John", "age": 30}, {"name": "Jane", "age": 25}]`
	args := []string{"-i", jsonInput, "--style", "csv"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Verify CSV output format
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < 2 {
		t.Fatalf("expected at least header and one data row in CSV output: %s", out)
	}
	// Check for comma-separated values
	if !strings.Contains(lines[0], ",") {
		t.Fatalf("expected comma-separated header in CSV output: %s", lines[0])
	}
	for i := 1; i < len(lines); i++ {
		if !strings.Contains(lines[i], ",") {
			t.Fatalf("expected comma-separated data in CSV output line %d: %s", i, lines[i])
		}
	}
}

func TestCLI_FilteringBasic(t *testing.T) {
	// Test basic filtering with equality
	jsonInput := `[{"name": "John", "age": 30}, {"name": "Jane", "age": 25}]`
	args := []string{"-i", jsonInput, "--where", "name=John", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should only contain John's row
	if !strings.Contains(out, "John") || strings.Contains(out, "Jane") {
		t.Fatalf("filter name=John failed: %s", out)
	}
}

func TestCLI_FilteringNumeric(t *testing.T) {
	// Test numeric comparison filtering
	jsonInput := `[{"name": "John", "age": 30}, {"name": "Jane", "age": 25}, {"name": "Bob", "age": 35}]`
	args := []string{"-i", jsonInput, "--where", "age>30", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should only contain Bob's row (age 35)
	if !strings.Contains(out, "Bob") || strings.Contains(out, "John") || strings.Contains(out, "Jane") {
		t.Fatalf("filter age>30 failed: %s", out)
	}
}

func TestCLI_FilteringBoolean(t *testing.T) {
	// Test boolean filtering
	jsonInput := `[{"name": "John", "active": true}, {"name": "Jane", "active": false}]`
	args := []string{"-i", jsonInput, "--where", "active=true", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should only contain John's row
	if !strings.Contains(out, "John") || strings.Contains(out, "Jane") {
		t.Fatalf("filter active=true failed: %s", out)
	}
}

func TestCLI_FilteringMultiple(t *testing.T) {
	// Test multiple filter conditions (AND logic)
	jsonInput := `[{"name": "John", "age": 30, "active": true}, {"name": "Jane", "age": 25, "active": false}, {"name": "Bob", "age": 35, "active": true}]`
	args := []string{"-i", jsonInput, "--where", "active=true", "--where", "age>=30", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should contain John and Bob (both active=true and age>=30)
	if !strings.Contains(out, "John") || !strings.Contains(out, "Bob") || strings.Contains(out, "Jane") {
		t.Fatalf("multiple filters failed: %s", out)
	}
}

func TestCLI_FilteringContains(t *testing.T) {
	// Test string contains filtering
	jsonInput := `[{"name": "Johnson", "city": "NYC"}, {"name": "Jane", "city": "LA"}, {"name": "Bob", "city": "Boston"}]`
	args := []string{"-i", jsonInput, "--where", "name~o", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should contain Johnson and Bob (both have 'o' in name)
	if !strings.Contains(out, "Johnson") || !strings.Contains(out, "Bob") || strings.Contains(out, "Jane") {
		t.Fatalf("filter name~o failed: %s", out)
	}
}

func TestCLI_FilteringRegex(t *testing.T) {
	// Test regex filtering
	jsonInput := `[{"email": "john@example.com"}, {"email": "jane@test.org"}, {"email": "bob@example.net"}]`
	args := []string{"-i", jsonInput, "--where", "email=~.*@example\\.(com|net)", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should contain john@example.com and bob@example.net
	if !strings.Contains(out, "john@example.com") || !strings.Contains(out, "bob@example.net") || strings.Contains(out, "jane@test.org") {
		t.Fatalf("regex filter failed: %s", out)
	}
}

func TestCLI_FilteringWithFlattening(t *testing.T) {
	// Test filtering with nested data (flattening)
	jsonInput := `[{"user": {"name": "John", "profile": {"age": 30}}, "active": true}, {"user": {"name": "Jane", "profile": {"age": 25}}, "active": false}]`
	args := []string{"-i", jsonInput, "--dive", "--where", "user.profile.age>25", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should only contain John's row (age 30 > 25)
	if !strings.Contains(out, "John") || strings.Contains(out, "Jane") {
		t.Fatalf("filter with flattening failed: %s", out)
	}
}

func TestCLI_FilteringWithSelection(t *testing.T) {
	// Test filtering combined with column selection
	jsonInput := `[{"name": "John", "age": 30, "city": "NYC", "active": true}, {"name": "Jane", "age": 25, "city": "LA", "active": false}]`
	args := []string{"-i", jsonInput, "--where", "active=true", "--select", "name,age", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should contain only John's row and only name, age columns
	if !strings.Contains(out, "John") || !strings.Contains(out, "30") || strings.Contains(out, "Jane") {
		t.Fatalf("filter with selection failed: %s", out)
	}
	if strings.Contains(out, "NYC") || strings.Contains(out, "active") {
		t.Fatalf("column selection failed with filtering: %s", out)
	}
}

func TestCLI_FilteringInvalidExpression(t *testing.T) {
	// Test error handling for invalid filter expressions
	jsonInput := `[{"name": "John", "age": 30}]`
	args := []string{"-i", jsonInput, "--where", "invalid_expression", "--style", "ascii"}
	_, errOut, code, _ := runCLI(t, args, nil)
	if code == 0 {
		t.Fatalf("expected non-zero exit code for invalid filter expression")
	}
	if !strings.Contains(errOut, "invalid filter condition") || !strings.Contains(errOut, "no valid operator found") {
		t.Fatalf("expected specific error message, got: %s", errOut)
	}
}

func TestCLI_FilteringNoMatches(t *testing.T) {
	// Test filtering that matches no rows
	jsonInput := `[{"name": "John", "age": 30}, {"name": "Jane", "age": 25}]`
	args := []string{"-i", jsonInput, "--where", "age>50", "--style", "ascii"}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should have no data rows
	if strings.Contains(out, "John") || strings.Contains(out, "Jane") || strings.Contains(out, "30") || strings.Contains(out, "25") {
		t.Fatalf("filter should match no rows: %s", out)
	}
	// When no rows match, output should be empty (this is expected behavior)
	if strings.TrimSpace(out) != "" {
		t.Fatalf("expected empty output when no rows match, got: %s", out)
	}
}

func TestCLI_NoInput_ShowsHelp(t *testing.T) {
	// Test that help is displayed when no input is provided
	args := []string{}
	out, errOut, code, err := runCLI(t, args, nil)
	if err != nil || code != 0 {
		t.Fatalf("err=%v code=%d stderr=%s", err, code, errOut)
	}
	// Should display help text
	if !strings.Contains(out, "Usage:") {
		t.Fatalf("expected help text with 'Usage:', got: %s", out)
	}
	if !strings.Contains(out, "Render JSON/YAML as tables") {
		t.Fatalf("expected help text with description, got: %s", out)
	}
	if !strings.Contains(out, "Flags:") {
		t.Fatalf("expected help text with 'Flags:', got: %s", out)
	}
}
