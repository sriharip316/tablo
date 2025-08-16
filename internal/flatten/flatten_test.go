package flatten

import (
	"encoding/json"
	"reflect"
	"testing"
)

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

func TestFlatten_Disabled(t *testing.T) {
	obj := map[string]any{"a": 1, "b": map[string]any{"c": 2}, "d": []any{"x", 2}}
	kv := FlattenObject(obj, Options{Enabled: false})
	if _, ok := kv["a"]; !ok {
		t.Fatalf("expected key a")
	}
	if s, ok := kv["b"].(string); !ok || s == "" {
		t.Fatalf("expected stringified b, got %T %v", kv["b"], kv["b"])
	}
}

func TestFlatten_WithDepthAndArrays(t *testing.T) {
	obj := map[string]any{
		"a":     map[string]any{"b": map[string]any{"c": 3}},
		"arr":   []any{map[string]any{"x": 1}, map[string]any{"x": 2}},
		"prims": []any{"p", 1},
	}
	kv := FlattenObject(obj, Options{Enabled: true, MaxDepth: 2})
	// depth is counted from top-level=1, so c at depth=3 exceeds and becomes string at a.b.c
	if v, ok := kv["a.b.c"].(string); !ok || v == "" {
		t.Fatalf("expected a.b.c stringified, got: %T %v", kv["a.b.c"], kv["a.b.c"])
	}
	if _, ok := kv["arr.0.x"]; !ok {
		t.Fatalf("missing arr.0.x")
	}
	if _, ok := kv["prims"].(string); !ok {
		t.Fatalf("expected prims stringified")
	}
}

func TestFlatten_SimpleArrayCSV(t *testing.T) {
	obj := map[string]any{"tags": []any{"a", "b", 3}}
	kv := FlattenObject(obj, Options{Enabled: true, FlattenSimpleArray: true, MaxDepth: -1})
	if kv["tags"] != "a, b, 3" {
		t.Fatalf("unexpected: %v", kv["tags"])
	}
}

func TestFlattenRows_Mixed(t *testing.T) {
	arr := []any{
		map[string]any{"a": 1},
		2,
	}
	rows := FlattenRows(arr, Options{Enabled: true})
	if len(rows) != 2 {
		t.Fatalf("want 2 rows")
	}
	if rows[1]["VALUE"] != 2 {
		t.Fatalf("second row VALUE mismatch: %+v", rows[1])
	}
}

func TestFlatKV_KeysSorted(t *testing.T) {
	kv := FlatKV{"b": 1, "a": 2}
	got := kv.Keys()
	exp := []string{"a", "b"}
	if !reflect.DeepEqual(got, exp) {
		t.Fatalf("keys order: %v", got)
	}
}

func TestStringifyStable(t *testing.T) {
	m := map[string]any{"a": 1}
	s := stringify(m)
	var back map[string]any
	if err := json.Unmarshal([]byte(s), &back); err != nil {
		t.Fatalf("bad json: %v", err)
	}
}
