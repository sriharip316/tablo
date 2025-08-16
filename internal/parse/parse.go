package parse

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

type Format string

const (
	Auto Format = "auto"
	JSON Format = "json"
	YAML Format = "yaml"
	YML  Format = "yml"
)

type Detector struct {
	Explicit string
	FilePath string
}

func (d Detector) Detect(data []byte) Format {
	if d.Explicit != "" && d.Explicit != string(Auto) {
		switch strings.ToLower(d.Explicit) {
		case "json":
			return JSON
		case "yaml", "yml":
			return YAML
		}
	}
	// by extension
	low := strings.ToLower(d.FilePath)
	if strings.HasSuffix(low, ".json") {
		return JSON
	}
	if strings.HasSuffix(low, ".yaml") || strings.HasSuffix(low, ".yml") {
		return YAML
	}
	// sniff
	trim := bytes.TrimLeft(data, " \t\r\n")
	if len(trim) > 0 && (trim[0] == '{' || trim[0] == '[') {
		return JSON
	}
	return YAML
}

func Parse(data []byte, f Format) (any, error) {
	switch f {
	case JSON:
		var v any
		dec := json.NewDecoder(bytes.NewReader(data))
		dec.UseNumber()
		if err := dec.Decode(&v); err != nil {
			// try YAML fallback in auto-like cases
			return nil, err
		}
		return normalize(v), nil
	case YAML, YML, Auto:
		// YAML can have multi-docs
		dec := yaml.NewDecoder(bytes.NewReader(data))
		var docs []any
		for {
			var v any
			if err := dec.Decode(&v); err != nil {
				if errors.Is(err, io.EOF) {
					break
				}
				return nil, err
			}
			if v == nil {
				continue
			}
			docs = append(docs, normalize(v))
		}
		if len(docs) == 1 {
			return docs[0], nil
		}
		return docs, nil
	default:
		return nil, ErrInvalidFormat
	}
}

var ErrInvalidFormat = errors.New("invalid format")

// normalize YAML maps/ints etc.
func normalize(v any) any {
	switch t := v.(type) {
	case map[any]any:
		m := make(map[string]any, len(t))
		for k, vv := range t {
			m[fmt.Sprint(k)] = normalize(vv)
		}
		return m
	case []any:
		out := make([]any, len(t))
		for i := range t {
			out[i] = normalize(t[i])
		}
		return out
	default:
		return v
	}
}

func ToStringKeyMap(m map[any]any) map[string]any {
	res := make(map[string]any, len(m))
	for k, v := range m {
		res[fmt.Sprint(k)] = normalize(v)
	}
	return res
}

func ArrayIsObjects(arr []any) bool {
	allObjs := true
	for _, it := range arr {
		_, is1 := it.(map[string]any)
		_, is2 := it.(map[any]any)
		if !(is1 || is2) {
			allObjs = false
			break
		}
	}
	return allObjs
}
