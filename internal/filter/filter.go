package filter

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/sriharip316/tablo/internal/flatten"
)

// Operator represents comparison operators for filtering
type Operator int

const (
	OpEqual Operator = iota
	OpNotEqual
	OpGreaterThan
	OpGreaterThanEqual
	OpLessThan
	OpLessThanEqual
	OpContains
	OpNotContains
	OpMatch
	OpNotMatch
)

// String returns the string representation of the operator
func (op Operator) String() string {
	switch op {
	case OpEqual:
		return "="
	case OpNotEqual:
		return "!="
	case OpGreaterThan:
		return ">"
	case OpGreaterThanEqual:
		return ">="
	case OpLessThan:
		return "<"
	case OpLessThanEqual:
		return "<="
	case OpContains:
		return "~"
	case OpNotContains:
		return "!~"
	case OpMatch:
		return "=~"
	case OpNotMatch:
		return "!=~"
	default:
		return "unknown"
	}
}

// Condition represents a single filter condition
type Condition struct {
	Path     string
	Operator Operator
	Value    string
	regex    *regexp.Regexp // compiled regex for match operators
}

// Filter represents a collection of filter conditions
type Filter struct {
	Conditions []Condition
}

// ParseCondition parses a filter condition string like "name=John" or "age>25"
func ParseCondition(expr string) (Condition, error) {
	expr = strings.TrimSpace(expr)
	if expr == "" {
		return Condition{}, fmt.Errorf("empty filter expression")
	}

	// Try different operators in order of longest first to avoid conflicts
	operators := []struct {
		op  Operator
		str string
	}{
		{OpNotMatch, "!=~"},
		{OpGreaterThanEqual, ">="},
		{OpLessThanEqual, "<="},
		{OpNotEqual, "!="},
		{OpNotContains, "!~"},
		{OpMatch, "=~"},
		{OpEqual, "="},
		{OpGreaterThan, ">"},
		{OpLessThan, "<"},
		{OpContains, "~"},
	}

	for _, opDef := range operators {
		if idx := strings.Index(expr, opDef.str); idx > 0 {
			path := strings.TrimSpace(expr[:idx])
			value := strings.TrimSpace(expr[idx+len(opDef.str):])

			condition := Condition{
				Path:     path,
				Operator: opDef.op,
				Value:    value,
			}

			// Compile regex for match operators
			if opDef.op == OpMatch || opDef.op == OpNotMatch {
				regex, err := regexp.Compile(value)
				if err != nil {
					return Condition{}, fmt.Errorf("invalid regex pattern %q: %w", value, err)
				}
				condition.regex = regex
			}

			return condition, nil
		}
	}

	return Condition{}, fmt.Errorf("invalid filter expression %q: no valid operator found", expr)
}

// ParseConditions parses multiple filter condition strings
func ParseConditions(exprs []string) ([]Condition, error) {
	conditions := make([]Condition, 0, len(exprs))
	for _, expr := range exprs {
		expr = strings.TrimSpace(expr)
		if expr == "" {
			continue
		}
		condition, err := ParseCondition(expr)
		if err != nil {
			return nil, err
		}
		conditions = append(conditions, condition)
	}
	return conditions, nil
}

// NewFilter creates a new filter with the given conditions
func NewFilter(conditions []Condition) *Filter {
	return &Filter{Conditions: conditions}
}

// Apply applies the filter to a slice of flattened rows
func (f *Filter) Apply(rows []flatten.FlatKV) []flatten.FlatKV {
	if len(f.Conditions) == 0 {
		return rows
	}

	filtered := make([]flatten.FlatKV, 0, len(rows))
	for _, row := range rows {
		if f.matchesRow(row) {
			filtered = append(filtered, row)
		}
	}
	return filtered
}

// matchesRow checks if a row matches all filter conditions (AND logic)
func (f *Filter) matchesRow(row flatten.FlatKV) bool {
	for _, condition := range f.Conditions {
		if !f.matchesCondition(row, condition) {
			return false
		}
	}
	return true
}

// matchesCondition checks if a row matches a single condition
func (f *Filter) matchesCondition(row flatten.FlatKV, condition Condition) bool {
	value, exists := row[condition.Path]

	// Handle missing values - only equal to empty string or null
	if !exists {
		switch condition.Operator {
		case OpEqual:
			return condition.Value == "" || condition.Value == "null"
		case OpNotEqual:
			return condition.Value != "" && condition.Value != "null"
		default:
			return false
		}
	}

	return f.compareValues(value, condition)
}

// compareValues compares a row value against a condition
func (f *Filter) compareValues(rowValue any, condition Condition) bool {
	// Convert row value to string for comparison
	rowStr := f.valueToString(rowValue)
	condStr := condition.Value

	switch condition.Operator {
	case OpEqual:
		return f.equalComparison(rowValue, condStr)
	case OpNotEqual:
		return !f.equalComparison(rowValue, condStr)
	case OpGreaterThan:
		return f.numericComparison(rowValue, condStr) > 0
	case OpGreaterThanEqual:
		return f.numericComparison(rowValue, condStr) >= 0
	case OpLessThan:
		return f.numericComparison(rowValue, condStr) < 0
	case OpLessThanEqual:
		return f.numericComparison(rowValue, condStr) <= 0
	case OpContains:
		return strings.Contains(rowStr, condStr)
	case OpNotContains:
		return !strings.Contains(rowStr, condStr)
	case OpMatch:
		return condition.regex != nil && condition.regex.MatchString(rowStr)
	case OpNotMatch:
		return condition.regex != nil && !condition.regex.MatchString(rowStr)
	default:
		return false
	}
}

// equalComparison performs type-aware equality comparison
func (f *Filter) equalComparison(rowValue any, condStr string) bool {
	// Handle null/nil values
	if rowValue == nil {
		return condStr == "null" || condStr == ""
	}

	// Handle boolean values
	if b, ok := rowValue.(bool); ok {
		if condBool, err := strconv.ParseBool(condStr); err == nil {
			return b == condBool
		}
		// Also allow string representation
		return (b && condStr == "true") || (!b && condStr == "false")
	}

	// Handle numeric values
	if f.isNumeric(rowValue) && f.isNumericString(condStr) {
		return f.numericComparison(rowValue, condStr) == 0
	}

	// Default to string comparison
	return f.valueToString(rowValue) == condStr
}

// numericComparison compares numeric values, returns -1, 0, or 1
func (f *Filter) numericComparison(rowValue any, condStr string) int {
	rowNum := f.toFloat64(rowValue)
	condNum, err := strconv.ParseFloat(condStr, 64)
	if err != nil {
		// If condition is not numeric, fall back to string comparison
		rowStr := f.valueToString(rowValue)
		if rowStr < condStr {
			return -1
		} else if rowStr > condStr {
			return 1
		}
		return 0
	}

	if rowNum < condNum {
		return -1
	} else if rowNum > condNum {
		return 1
	}
	return 0
}

// valueToString converts any value to its string representation
func (f *Filter) valueToString(value any) string {
	if value == nil {
		return ""
	}
	switch v := value.(type) {
	case string:
		return v
	case bool:
		return strconv.FormatBool(v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float32:
		return strconv.FormatFloat(float64(v), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	default:
		return fmt.Sprintf("%v", v)
	}
}

// isNumeric checks if a value is numeric
func (f *Filter) isNumeric(value any) bool {
	switch value.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return true
	case float32, float64:
		return true
	default:
		return false
	}
}

// isNumericString checks if a string represents a numeric value
func (f *Filter) isNumericString(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

// toFloat64 converts numeric values to float64
func (f *Filter) toFloat64(value any) float64 {
	switch v := value.(type) {
	case int:
		return float64(v)
	case int8:
		return float64(v)
	case int16:
		return float64(v)
	case int32:
		return float64(v)
	case int64:
		return float64(v)
	case uint:
		return float64(v)
	case uint8:
		return float64(v)
	case uint16:
		return float64(v)
	case uint32:
		return float64(v)
	case uint64:
		return float64(v)
	case float32:
		return float64(v)
	case float64:
		return v
	default:
		// Try to parse as string
		if s := f.valueToString(value); s != "" {
			if f, err := strconv.ParseFloat(s, 64); err == nil {
				return f
			}
		}
		return 0
	}
}
