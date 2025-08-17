package selectors

import (
	"testing"

	"github.com/sriharip316/tablo/internal/flatten"
)

func TestCompileAndMatch(t *testing.T) {
	exprs, err := CompileMany([]string{"a.*.c", "x"})
	if err != nil {
		t.Fatal(err)
	}
	keys := []string{"a.1.c", "a.2.d", "x"}
	res := ApplyToKeys(keys, exprs, nil)
	if len(res) != 2 || res[0] != "a.1.c" || res[1] != "x" {
		t.Fatalf("unexpected %v", res)
	}
}

func TestApplyToKeys_IncludeExcludeOrder(t *testing.T) {
	keys := []string{"a", "b.c", "b.d", "x"}
	inc, _ := CompileMany([]string{"b.*", "x"})
	exc, _ := CompileMany([]string{"b.d"})
	out := ApplyToKeys(keys, inc, exc)
	if len(out) != 2 || out[0] != "b.c" || out[1] != "x" {
		t.Fatalf("unexpected: %v", out)
	}
}

func TestMissingExpressions(t *testing.T) {
	keys := []string{"a.b", "c.d"}
	inc, _ := CompileMany([]string{"a.*", "x"})
	miss := MissingExpressions(keys, inc)
	if len(miss) != 1 || miss[0] != "x" {
		t.Fatalf("missing: %v", miss)
	}
}

func TestGlobToRegex(t *testing.T) {
	re := globToRegex("a.*.c")
	if re != "a\\..*\\.c" {
		t.Fatalf("got %q", re)
	}
}

// Placeholder test ensuring representative glob patterns compile successfully.
// We only use patterns that translate cleanly to regex via globToRegex.
func TestCompileMany_GlobWeirdButValid(t *testing.T) {
	cases := []string{"[abc]*", "a?b", "plain"}
	if _, err := CompileMany(cases); err != nil {
		t.Fatalf("expected all glob patterns to compile, got err=%v", err)
	}
}

func TestHeadersUnionOrder(t *testing.T) {
	rows := []flatten.FlatKV{
		{"a": 1, "b": 2},
		{"b": 3, "c": 4},
		{"a": 5, "d": 6, "c": 7},
	}
	got := HeadersUnion(rows)
	// expected order: a (from row1), b (row1), c (row2), d (row3)
	exp := []string{"a", "b", "c", "d"}
	if len(got) != len(exp) {
		t.Fatalf("len mismatch got=%v exp=%v", got, exp)
	}
	for i, k := range exp {
		if got[i] != k {
			t.Fatalf("order mismatch got=%v exp=%v", got, exp)
		}
	}
}

func TestHeadersUnionEmpty(t *testing.T) {
	if got := HeadersUnion([]flatten.FlatKV{}); len(got) != 0 {
		t.Fatalf("expected empty got=%v", got)
	}
}

func TestApplyToKeys_ExcludeOnly(t *testing.T) {
	keys := []string{"a", "b.c", "b.d", "c.e"}
	exc, err := CompileMany([]string{"b.*"})
	if err != nil {
		t.Fatalf("compile: %v", err)
	}
	out := ApplyToKeys(keys, nil, exc)
	// expect a, c.e (order preserved)
	if len(out) != 2 || out[0] != "a" || out[1] != "c.e" {
		t.Fatalf("unexpected exclude-only result: %v", out)
	}
}

func TestApplyToKeys_NoFilters(t *testing.T) {
	keys := []string{"x", "y"}
	out := ApplyToKeys(keys, nil, nil)
	if len(out) != 2 || out[0] != "x" || out[1] != "y" {
		t.Fatalf("unexpected passthrough: %v", out)
	}
}
