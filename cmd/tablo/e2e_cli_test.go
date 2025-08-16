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
