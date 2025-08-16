package selectors

import "testing"

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
