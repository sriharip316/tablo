package flatten

import "testing"

func TestFlattenSimple(t *testing.T) {
	m := map[string]any{"a": 1, "b": map[string]any{"c": 2}}
	kv := FlattenObject(m, Options{Enabled: true, MaxDepth: -1})
	if kv["a"] != 1 {
		t.Fatalf("unexpected: %+v", kv)
	}
	if v, ok := kv["b.c"]; !ok || v != 2 {
		t.Fatalf("expected b.c to be 2, got %+v", kv)
	}
}
