package flatten

import (
	"encoding/json"
	"reflect"
	"strings"
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

// Test flattening when DivePaths restrict which top-level keys are flattened.
func TestFlatten_DivePathsSelective(t *testing.T) {
	obj := map[string]any{
		"a": map[string]any{
			"x": 1,
			"y": map[string]any{"z": 2},
		},
		"b": map[string]any{
			"q": 3,
		},
		"tags": []any{"p", "q", 9},
	}

	kv := FlattenObject(obj, Options{
		Enabled:            true,
		DivePaths:          []string{"a"},
		MaxDepth:           -1,
		FlattenSimpleArray: true,
	})

	// Expect flattened keys for "a"
	if kv["a.x"] != 1 {
		t.Fatalf("missing a.x %+v", kv)
	}
	if kv["a.y.z"] != 2 {
		t.Fatalf("missing a.y.z %+v", kv)
	}

	// "b" should NOT be flattened because only "a" is allowed; should be stringified JSON
	bv, ok := kv["b"].(string)
	if !ok || !strings.Contains(bv, `"q":3`) {
		t.Fatalf("expected stringified b, got %#v", kv["b"])
	}

	// tags (simple array under non-dive key) should be CSV because FlattenSimpleArray=true
	if kv["tags"] != "p, q, 9" {
		t.Fatalf("expected CSV tags, got %#v", kv["tags"])
	}
}

// Test MaxDepth exact edge: when depth limit is reached, nested content is stringified.
func TestFlatten_MaxDepthEdge(t *testing.T) {
	obj := map[string]any{
		"top": map[string]any{
			"mid": map[string]any{
				"leaf": 42,
			},
		},
	}

	// MaxDepth=1: depth=1 for top, depth=2 for mid exceeds limit so mid map stringified.
	kv := FlattenObject(obj, Options{Enabled: true, MaxDepth: 1})
	// Expect no fully flattened leaf; got stringified mid map at key top.mid
	if _, ok := kv["top.mid.leaf"]; ok {
		t.Fatalf("did not expect key top.mid.leaf at MaxDepth=1: %+v", kv)
	}
	mid, ok := kv["top.mid"].(string)
	if !ok || !strings.Contains(mid, `"leaf":42`) {
		t.Fatalf("expected stringified mid map, got %#v", kv["top.mid"])
	}

	// MaxDepth=2 should allow flattening one more level producing top.mid.leaf (value may be raw int 42 or string "42"
	// because depth>MaxDepth stringification can JSON-encode the scalar).
	kv2 := FlattenObject(obj, Options{Enabled: true, MaxDepth: 2})
	vRaw, ok := kv2["top.mid.leaf"]
	if !ok {
		t.Fatalf("expected flattened leaf key at depth=2, got %+v", kv2)
	}
	if vRaw != 42 && vRaw != "42" {
		t.Fatalf("expected leaf value 42 (int or string), got %#v (full map: %+v)", vRaw, kv2)
	}
}

// Test MaxDepth=0 edge: enabled flattening but zero depth means no flattening beyond top-level call.
func TestFlatten_MaxDepthZero(t *testing.T) {
	obj := map[string]any{
		"a": map[string]any{"b": 1},
		"x": 7,
	}
	kv := FlattenObject(obj, Options{Enabled: true, MaxDepth: 0})
	// "a" should remain stringified, not "a.b"
	if _, exists := kv["a.b"]; exists {
		t.Fatalf("did not expect a.b with MaxDepth=0: %+v", kv)
	}
	if s, ok := kv["a"].(string); !ok || !strings.Contains(s, `"b":1`) {
		t.Fatalf("expected stringified a, got %#v", kv["a"])
	}
	if kv["x"] != "7" {
		t.Fatalf("expected stringified scalar x, got %#v", kv["x"])
	}
}

// Test that when Enabled=false but DivePaths provided, still no flattening occurs.
func TestFlatten_DivePathsIgnoredWhenDisabled(t *testing.T) {
	obj := map[string]any{
		"a": map[string]any{"b": 1},
	}
	kv := FlattenObject(obj, Options{Enabled: false, DivePaths: []string{"a"}})
	if _, ok := kv["a.b"]; ok {
		t.Fatalf("unexpected flattened key when disabled: %+v", kv)
	}
	if s, ok := kv["a"].(string); !ok || !strings.Contains(s, `"b":1`) {
		t.Fatalf("expected stringified a when disabled, got %#v", kv["a"])
	}
}

// Test that arrays of objects are flattened only for allowed dive paths.
func TestFlatten_DivePaths_ArrayOfObjects(t *testing.T) {
	obj := map[string]any{
		"keep": []any{
			map[string]any{"x": 1},
			map[string]any{"x": 2},
		},
		"dive": []any{
			map[string]any{"y": 3},
			map[string]any{"y": 4},
		},
	}

	kv := FlattenObject(obj, Options{Enabled: true, DivePaths: []string{"dive"}, MaxDepth: -1})
	// Expect flattened indices for dive.*
	if kv["dive.0.y"] != 3 || kv["dive.1.y"] != 4 {
		t.Fatalf("missing flattened dive indices: %+v", kv)
	}
	// keep should be stringified (array of objects but not in DivePaths)
	if s, ok := kv["keep"].(string); !ok || !strings.Contains(s, `"x":1`) {
		t.Fatalf("expected stringified keep array, got %#v", kv["keep"])
	}
}
