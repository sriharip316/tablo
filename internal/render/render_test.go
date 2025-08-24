package render

import (
	stdjson "encoding/json"
	"strings"
	"testing"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/sriharip316/tablo/internal/flatten"
)

func TestRenderObjectKV(t *testing.T) {
	kv := flatten.FlatKV{"a": 1, "b.c": 2}
	m := Model{Mode: ModeObjectKV, KV: kv}
	out, err := Render(m, Options{Style: "ascii"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Fatal("empty output")
	}
}

func TestRenderRows_WithIndexAndFormat(t *testing.T) {
	rows := []flatten.FlatKV{{"a": 1.2345, "b": true}, {"b": false}}
	m := FromFlatRows(rows, []string{"a", "b"}, true)
	out, err := Render(m, Options{Style: "ascii", Precision: 2, BoolStr: "Y:N", NullStr: "null"})
	if err != nil {
		t.Fatal(err)
	}
	// Expect first column to be auto index numbers 1 and 2
	if !strings.Contains(out, "| 1 ") || !strings.Contains(out, "| 2 ") {
		t.Fatalf("expected index values present: \n%s", out)
	}
	if !strings.Contains(out, "1.23") { // precision applied
		t.Fatalf("precision not applied: \n%s", out)
	}
	if !strings.Contains(out, "Y") || !strings.Contains(out, "N") {
		t.Fatalf("bool formatting missing: \n%s", out)
	}
}

func TestRender_PrimitiveArrayLimit(t *testing.T) {
	m := FromPrimitiveArray([]any{1, 2, 3}, false, 2)
	out, err := Render(m, Options{Style: "markdown"})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(out, "\n") < 2 {
		t.Fatalf("unexpected markdown output: %q", out)
	}
}

func TestHeaderCase(t *testing.T) {
	if headerCase("hello_world", "title") != "Hello World" {
		t.Fatalf("title case failed")
	}
}

func TestWrapEnforcer_Truncate(t *testing.T) {
	en := wrapEnforcer(Options{WrapMode: "off", TruncateSuffix: ".."})

	// Test case: basic truncation
	got := en("abcdefgh", 5)
	if got != "abc.." {
		t.Fatalf("truncate: %q, want %q", got, "abc..")
	}

	// Test case: width <= 0 should return original string
	got = en("abcdefgh", 0)
	if got != "abcdefgh" {
		t.Fatalf("width 0: %q, want %q", got, "abcdefgh")
	}
	got = en("abcdefgh", -5)
	if got != "abcdefgh" {
		t.Fatalf("width -5: %q, want %q", got, "abcdefgh")
	}

	// Test case: len(s) <= width should return original string
	got = en("abc", 5)
	if got != "abc" {
		t.Fatalf("len <= width: %q, want %q", got, "abc")
	}

	// Test case: len(suffix) > width, should not add suffix, truncate if needed
	en = wrapEnforcer(Options{WrapMode: "off", TruncateSuffix: "longsuffix"})
	got = en("abc", 2) // width 2, suffix "longsuffix" (len 10). Should truncate to "ab"
	if got != "ab" {
		t.Fatalf("suffix longer than width, string fits: %q, want %q", got, "ab")
	}
	got = en("abcdefgh", 5) // width 5, suffix "longsuffix" (len 10). Should truncate to "abcde"
	if got != "abcde" {
		t.Fatalf("suffix longer than width, string longer than width: %q, want %q", got, "abcde")
	}

	// Test cases for "char" and "word" wrap modes (only ensure they return different enforcers)
	enChar := wrapEnforcer(Options{WrapMode: "char"})
	enWord := wrapEnforcer(Options{WrapMode: "word"})
	enOff := wrapEnforcer(Options{WrapMode: "off"})

	// We can't directly compare function pointers, so we compare behavior for a simple case.
	// This mainly checks that it's not the "off" mode's truncation for "char" and "word".
	// The actual wrapping logic is handled by the `go-pretty` library's `text.WrapText` and `text.WrapSoft`.
	if enChar("test", 2) == enOff("test", 2) {
		t.Fatalf("char wrap enforcer appears to be same as off wrap enforcer")
	}
	if enWord("test", 2) == enOff("test", 2) {
		t.Fatalf("word wrap enforcer appears to be same as off wrap enforcer")
	}
}

func TestResolveStyleVariants(t *testing.T) {
	// heavy -> StyleBold (sanity: header formatter default)
	sHeavy := resolveStyle(Options{Style: "heavy"})
	if sHeavy.Name != table.StyleBold.Name {
		t.Fatalf("heavy style mismatch: got %s want %s", sHeavy.Name, table.StyleBold.Name)
	}

	// light
	sLight := resolveStyle(Options{Style: "light"})
	if sLight.Name != table.StyleLight.Name {
		t.Fatalf("light style mismatch: got %s", sLight.Name)
	}

	// double
	sDouble := resolveStyle(Options{Style: "double"})
	if sDouble.Name != table.StyleDouble.Name {
		t.Fatalf("double style mismatch: got %s", sDouble.Name)
	}

	// ascii (defaults to StyleDefault)
	sASCII := resolveStyle(Options{Style: "ascii"})
	if sASCII.Name != table.StyleDefault.Name {
		t.Fatalf("ascii style mismatch: got %s", sASCII.Name)
	}

	// markdown (should still use default base)
	sMarkdown := resolveStyle(Options{Style: "markdown"})
	if sMarkdown.Name != table.StyleDefault.Name {
		t.Fatalf("markdown base style mismatch")
	}

	// compact -> StyleLight w/ SeparateRows disabled
	sCompact := resolveStyle(Options{Style: "compact"})
	if sCompact.Name != table.StyleLight.Name {
		t.Fatalf("compact base mismatch")
	}
	if sCompact.Options.SeparateRows {
		t.Fatalf("compact should disable SeparateRows")
	}

	// borderless -> StyleLight + no borders/separators
	sBorderless := resolveStyle(Options{Style: "borderless"})
	if sBorderless.Name != table.StyleLight.Name {
		t.Fatalf("borderless base mismatch")
	}
	if sBorderless.Options.DrawBorder {
		t.Fatalf("borderless should not draw border")
	}

	// ASCIIOnly override sets Box to StyleBoxDefault
	sHeavyASCII := resolveStyle(Options{Style: "heavy", ASCIIOnly: true})
	if sHeavyASCII.Box.TopLeft != table.StyleBoxDefault.TopLeft {
		t.Fatalf("ASCIIOnly override failed: %+v", sHeavyASCII.Box)
	}

	// HeaderCase variations: ensure explicit upper case selection
	sUpper := resolveStyle(Options{Style: "ascii", HeaderCase: "upper"})
	if sUpper.Format.Header != text.FormatUpper {
		t.Fatalf("expected upper header format, got %+v", sUpper.Format.Header)
	}
}

func TestFormatCellBranches(t *testing.T) {
	// nil value with explicit NullStr
	if v := formatCell(nil, Options{NullStr: "<nil>"}); v != "<nil>" {
		t.Fatalf("nil formatting failed: %v", v)
	}

	// bool with mapping
	if v := formatCell(true, Options{BoolStr: "Y:N"}); v != "Y" {
		t.Fatalf("bool mapping true failed: %v", v)
	}
	if v := formatCell(false, Options{BoolStr: "Y:N"}); v != "N" {
		t.Fatalf("bool mapping false failed: %v", v)
	}

	// bool without mapping
	if v := formatCell(true, Options{}); v != true {
		t.Fatalf("bool default failed: %v (%T)", v, v)
	}

	// float64 with precision
	if v := formatCell(1.23456, Options{Precision: 2}); v != "1.23" {
		t.Fatalf("float64 precision failed: %v", v)
	}

	// float32 with precision
	if v := formatCell(float32(3.14159), Options{Precision: 3}); v != "3.142" {
		t.Fatalf("float32 precision failed: %v", v)
	}

	// json.Number with precision
	num := stdjson.Number("2.71828")
	if v := formatCell(num, Options{Precision: 3}); v != "2.718" {
		t.Fatalf("json.Number precision failed: %v", v)
	}

	// json.Number without precision (Precision = -1) should return original string
	if v := formatCell(stdjson.Number("42"), Options{Precision: -1}); v != "42" {
		t.Fatalf("json.Number default failed: %v", v)
	}

	// default branch (string passthrough)
	if v := formatCell("hello", Options{}); v != "hello" {
		t.Fatalf("string passthrough failed: %v", v)
	}
}

func TestRenderObjectKV_WithWrapping(t *testing.T) {
	kv := flatten.FlatKV{"name": "This is a very long text that should be wrapped when the column width is limited", "value": 42}
	m := Model{Mode: ModeObjectKV, KV: kv}

	// Test word wrapping
	out, err := Render(m, Options{Style: "ascii", MaxColWidth: 20, WrapMode: "word"})
	if err != nil {
		t.Fatal(err)
	}
	// Should contain wrapped text across multiple lines
	if !strings.Contains(out, "This is a very long") {
		t.Fatalf("expected wrapped text, got: %s", out)
	}
	// Should have multiple lines for the long text
	lines := strings.Split(out, "\n")
	longValueLines := 0
	for _, line := range lines {
		if strings.Contains(line, "text that should be") || strings.Contains(line, "wrapped when the") {
			longValueLines++
		}
	}
	if longValueLines == 0 {
		t.Fatalf("expected wrapped lines for long text, got: %s", out)
	}

	// Test truncation mode
	outTrunc, err := Render(m, Options{Style: "ascii", MaxColWidth: 20, WrapMode: "off", TruncateSuffix: "..."})
	if err != nil {
		t.Fatal(err)
	}
	// Should contain truncated text with suffix
	if !strings.Contains(outTrunc, "...") {
		t.Fatalf("expected truncated text with suffix, got: %s", outTrunc)
	}
}

func TestHeaderCaseAdditional(t *testing.T) {
	if got := headerCase("Mixed_Case", "lower"); got != "mixed_case" {
		t.Fatalf("lower headerCase failed: %s", got)
	}
	if got := headerCase("Mixed_Case", "upper"); got != "MIXED_CASE" {
		t.Fatalf("upper headerCase failed: %s", got)
	}
	// original fallback
	if got := headerCase("Already", "original"); got != "Already" {
		t.Fatalf("original headerCase failed: %s", got)
	}
}
