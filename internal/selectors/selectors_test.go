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
