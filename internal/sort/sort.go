package sort

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/sriharip316/tablo/internal/flatten"
)

// Options contains configuration for sorting rows
type Options struct {
	Columns []string // Column names to sort by (with optional +/- prefix)
}

// SortColumn represents a column with its sort direction
type SortColumn struct {
	Name       string
	Descending bool
}

// Sorter handles sorting of flattened rows
type Sorter struct {
	columns []SortColumn
}

// New creates a new Sorter with the given options
func New(opts Options) *Sorter {
	columns := parseColumns(opts.Columns)
	return &Sorter{
		columns: columns,
	}
}

// parseColumns parses column specifications with optional +/- prefixes
func parseColumns(columnSpecs []string) []SortColumn {
	if len(columnSpecs) == 0 {
		return []SortColumn{}
	}

	var columns []SortColumn

	for _, spec := range columnSpecs {
		spec = strings.TrimSpace(spec)
		if spec == "" {
			continue
		}

		var column SortColumn

		// Check for explicit direction prefix
		if strings.HasPrefix(spec, "+") {
			column.Descending = false
			column.Name = spec[1:]
		} else if strings.HasPrefix(spec, "-") {
			column.Descending = true
			column.Name = spec[1:]
		} else {
			// No prefix defaults to ascending
			column.Descending = false
			column.Name = spec
		}

		columns = append(columns, column)
	}

	return columns
}

// Sort sorts the given rows by the configured columns
func (s *Sorter) Sort(rows []flatten.FlatKV) []flatten.FlatKV {
	if len(s.columns) == 0 || len(rows) <= 1 {
		return rows
	}

	// Create a copy to avoid modifying the original slice
	sorted := make([]flatten.FlatKV, len(rows))
	copy(sorted, rows)

	// Sort using Go's stable sort
	sort.SliceStable(sorted, func(i, j int) bool {
		return s.compare(sorted[i], sorted[j])
	})

	return sorted
}

// compare compares two rows based on the configured sort columns
func (s *Sorter) compare(a, b flatten.FlatKV) bool {
	for _, col := range s.columns {
		valA := a[col.Name]
		valB := b[col.Name]

		cmp := compareValues(valA, valB)
		if cmp != 0 {
			if col.Descending {
				return cmp > 0
			}
			return cmp < 0
		}
		// Values are equal, continue to next column
	}
	// All columns are equal
	return false
}

// compareValues compares two values and returns:
// -1 if a < b
//
//	0 if a == b
//	1 if a > b
func compareValues(a, b any) int {
	// Handle nil values - nil sorts before any other value
	if a == nil && b == nil {
		return 0
	}
	if a == nil {
		return -1
	}
	if b == nil {
		return 1
	}

	// Check if both values are of the same comparable type
	// Numbers
	if numA, okA := toNumber(a); okA {
		if numB, okB := toNumber(b); okB {
			if numA < numB {
				return -1
			} else if numA > numB {
				return 1
			}
			return 0
		}
	}

	// Booleans
	if boolA, okA := toBool(a); okA {
		if boolB, okB := toBool(b); okB {
			if !boolA && boolB {
				return -1
			} else if boolA && !boolB {
				return 1
			}
			return 0
		}
	}

	// For mixed types or when both are strings, use string comparison
	strA := toString(a)
	strB := toString(b)

	if strA < strB {
		return -1
	} else if strA > strB {
		return 1
	}
	return 0
}

// toNumber attempts to convert a value to a float64
func toNumber(v any) (float64, bool) {
	switch val := v.(type) {
	case float64:
		return val, true
	case float32:
		return float64(val), true
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case string:
		if f, err := strconv.ParseFloat(val, 64); err == nil {
			return f, true
		}
	}
	return 0, false
}

// toBool attempts to convert a value to a boolean
func toBool(v any) (bool, bool) {
	switch val := v.(type) {
	case bool:
		return val, true
	case string:
		switch strings.ToLower(val) {
		case "true", "t", "yes", "y", "1":
			return true, true
		case "false", "f", "no", "n", "0":
			return false, true
		}
	}
	return false, false
}

// toString converts a value to a string for comparison
func toString(v any) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return strings.ToLower(fmt.Sprintf("%v", val))
	}
}
