package flatten

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
)

type Options struct {
	Enabled            bool
	MaxDepth           int // -1 unlimited
	DivePaths          []string
	FlattenSimpleArray bool
}

type FlatKV map[string]any

func (f FlatKV) Keys() []string {
	keys := make([]string, 0, len(f))
	for k := range f {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// FlattenObject flattens an object (map[string]any) respecting Options.
func FlattenObject(obj any, o Options) FlatKV {
	out := make(FlatKV)
	if !o.Enabled {
		// do not dive; stringify composite
		switch m := obj.(type) {
		case map[string]any:
			for k, v := range m {
				out[k] = maybeStringify(v, o)
			}
		default:
			// if not map, return key VALUE mapping
			out["VALUE"] = maybeStringify(obj, o)
		}
		return out
	}
	var walk func(prefix string, v any, depth int)
	walk = func(prefix string, v any, depth int) {
		if o.MaxDepth >= 0 && depth > o.MaxDepth {
			out[prefix] = stringify(v)
			return
		}
		switch vv := v.(type) {
		case map[string]any:
			for k, val := range vv {
				p := k
				if prefix != "" {
					p = prefix + "." + k
				}
				walk(p, val, depth+1)
			}
		case []any:
			// only flatten arrays of objects
			allObj := true
			for _, it := range vv {
				if _, ok := it.(map[string]any); !ok {
					allObj = false
					break
				}
			}
			if !allObj {
				if o.FlattenSimpleArray {
					out[prefix] = simpleArrayToCSV(vv)
				} else {
					out[prefix] = stringify(vv)
				}
				return
			}
			for i, it := range vv {
				p := prefix + "." + strconv.Itoa(i)
				walk(p, it, depth+1)
			}
		default:
			if prefix == "" {
				out["VALUE"] = vv
			} else {
				out[prefix] = vv
			}
		}
	}
	if m, ok := obj.(map[string]any); ok {
		// if DivePaths specified, only dive into those top-level keys
		allow := map[string]struct{}{}
		if len(o.DivePaths) > 0 {
			for _, p := range o.DivePaths {
				allow[p] = struct{}{}
			}
		}
		for k, v := range m {
			if len(allow) > 0 {
				if _, ok := allow[k]; ok {
					walk(k, v, 1)
				} else {
					// keep as is; if scalar, keep value; else stringify or CSV for simple arrays
					switch vv := v.(type) {
					case []any:
						if o.FlattenSimpleArray {
							out[k] = simpleArrayToCSV(vv)
						} else {
							out[k] = stringify(vv)
						}
					case map[string]any:
						out[k] = stringify(vv)
					default:
						out[k] = vv
					}
				}
			} else {
				walk(k, v, 1)
			}
		}
	}
	return out
}

// FlattenRows flattens an array of object-like values into rows of FlatKV.
func FlattenRows(arr []any, o Options) []FlatKV {
	rows := make([]FlatKV, 0, len(arr))
	for _, it := range arr {
		if m, ok := it.(map[string]any); ok {
			rows = append(rows, FlattenObject(m, o))
		} else {
			rows = append(rows, FlatKV{"VALUE": maybeStringify(it, o)})
		}
	}
	return rows
}

func maybeStringify(v any, o Options) any {
	switch vv := v.(type) {
	case map[string]any:
		return stringify(vv)
	case []any:
		if o.FlattenSimpleArray {
			return simpleArrayToCSV(vv)
		}
		return stringify(vv)
	default:
		return v
	}
}

func stringify(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}

func simpleArrayToCSV(v []any) string {
	parts := make([]string, len(v))
	for i, it := range v {
		switch t := it.(type) {
		case string:
			parts[i] = t
		default:
			parts[i] = stringify(t)
		}
	}
	// naive join with comma+space
	return strings.Join(parts, ", ")
}
