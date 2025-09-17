package parse

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/tidwall/jsonc"
	"gopkg.in/yaml.v3"
)

type Format string

const (
	Auto  Format = "auto"
	JSON  Format = "json"
	YAML  Format = "yaml"
	YML   Format = "yml"
	CSV   Format = "csv"
	JSONL Format = "jsonl"
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
		case "csv":
			return CSV
		case "jsonl":
			return JSONL
		}
	}
	// by extension
	low := strings.ToLower(d.FilePath)
	if strings.HasSuffix(low, ".json") || strings.HasSuffix(low, ".jsonc") {
		return JSON
	}
	if strings.HasSuffix(low, ".yaml") || strings.HasSuffix(low, ".yml") {
		return YAML
	}
	if strings.HasSuffix(low, ".csv") {
		return CSV
	}
	// sniff
	trim := bytes.TrimLeft(data, " \t\r\n")
	if len(trim) > 0 && (trim[0] == '{' || trim[0] == '[') {
		// Check if this might be JSONL (multiple lines)
		lines := strings.Split(string(trim), "\n")
		if len(lines) > 1 {
			allJSON := true
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "{") && !strings.HasPrefix(line, "[") {
					allJSON = false
					break
				}
			}
			if allJSON {
				return JSONL
			}
		}
		return JSON
	}
	// Check for JSON with leading comments
	if len(trim) > 0 && (trim[0] == '/' && len(trim) > 1 && (trim[1] == '/' || trim[1] == '*')) {
		// Strip comments and check if JSON follows
		cleanJSON := jsonc.ToJSON(data)
		cleanTrim := bytes.TrimLeft(cleanJSON, " \t\r\n")
		if len(cleanTrim) > 0 && (cleanTrim[0] == '{' || cleanTrim[0] == '[') {
			return JSON
		}
	}
	// Check for CSV by looking for comma-separated values in first line
	if len(trim) > 0 {
		firstLine := strings.Split(string(trim), "\n")[0]
		if strings.Contains(firstLine, ",") && !strings.Contains(firstLine, "{") && !strings.Contains(firstLine, ":") {
			return CSV
		}
	}
	return YAML
}

type ParseOptions struct {
	CSVNoHeader bool
}

func Parse(data []byte, f Format, opts ParseOptions) (any, error) {
	switch f {
	case JSON:
		// Strip comments from JSON using jsonc library
		cleanJSON := jsonc.ToJSON(data)
		var v any
		dec := json.NewDecoder(bytes.NewReader(cleanJSON))
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
	case CSV:
		return parseCSV(data, opts.CSVNoHeader)
	case JSONL:
		return parseJSONL(data)
	default:
		return nil, ErrInvalidFormat
	}
}

// parseCSV converts CSV data to []map[string]any
func parseCSV(data []byte, noHeader bool) (any, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return []map[string]any{}, nil
	}

	var headers []string
	var startRow int

	if noHeader {
		// Generate column headers as col0, col1, etc.
		headers = make([]string, len(records[0]))
		for i := range headers {
			headers[i] = fmt.Sprintf("col%d", i)
		}
		startRow = 0
	} else {
		headers = records[0]
		startRow = 1
	}

	var result []map[string]any
	for i := startRow; i < len(records); i++ {
		row := records[i]
		obj := make(map[string]any)
		for j, value := range row {
			if j < len(headers) {
				obj[headers[j]] = value
			}
		}
		result = append(result, obj)
	}
	return result, nil
}

var ErrInvalidFormat = errors.New("invalid format")

// parseJSONL converts JSON Lines data to []any
func parseJSONL(data []byte) (any, error) {
	lines := strings.Split(string(data), "\n")
	var result []any
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var v any
		if err := json.Unmarshal([]byte(line), &v); err != nil {
			return nil, err
		}
		// If the parsed value is an array, flatten it into individual objects
		if arr, ok := v.([]any); ok {
			for _, item := range arr {
				result = append(result, normalize(item))
			}
		} else {
			result = append(result, normalize(v))
		}
	}
	return result, nil
}

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
		if !is1 && !is2 {
			allObjs = false
			break
		}
	}
	return allObjs
}
