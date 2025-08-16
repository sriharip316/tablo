package render

import (
	"strings"
	"testing"

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
	got := en("abcdefgh", 5)
	if got != "abc.." {
		t.Fatalf("truncate: %q", got)
	}
}
