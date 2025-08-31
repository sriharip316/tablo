package filter

import (
	"testing"

	"github.com/sriharip316/tablo/internal/flatten"
)

func TestParseCondition(t *testing.T) {
	tests := []struct {
		name        string
		expr        string
		want        Condition
		wantErr     bool
		errContains string
	}{
		{
			name: "simple equality",
			expr: "name=John",
			want: Condition{Path: "name", Operator: OpEqual, Value: "John"},
		},
		{
			name: "not equal",
			expr: "status!=active",
			want: Condition{Path: "status", Operator: OpNotEqual, Value: "active"},
		},
		{
			name: "greater than",
			expr: "age>25",
			want: Condition{Path: "age", Operator: OpGreaterThan, Value: "25"},
		},
		{
			name: "greater than equal",
			expr: "score>=90",
			want: Condition{Path: "score", Operator: OpGreaterThanEqual, Value: "90"},
		},
		{
			name: "less than",
			expr: "price<100",
			want: Condition{Path: "price", Operator: OpLessThan, Value: "100"},
		},
		{
			name: "less than equal",
			expr: "count<=5",
			want: Condition{Path: "count", Operator: OpLessThanEqual, Value: "5"},
		},
		{
			name: "contains",
			expr: "description~error",
			want: Condition{Path: "description", Operator: OpContains, Value: "error"},
		},
		{
			name: "not contains",
			expr: "message!~warning",
			want: Condition{Path: "message", Operator: OpNotContains, Value: "warning"},
		},
		{
			name: "regex match",
			expr: "email=~.*@example\\.com",
			want: Condition{Path: "email", Operator: OpMatch, Value: ".*@example\\.com"},
		},
		{
			name: "regex not match",
			expr: "phone!=~\\d{3}-\\d{3}-\\d{4}",
			want: Condition{Path: "phone", Operator: OpNotMatch, Value: "\\d{3}-\\d{3}-\\d{4}"},
		},
		{
			name: "dotted path",
			expr: "user.profile.age>18",
			want: Condition{Path: "user.profile.age", Operator: OpGreaterThan, Value: "18"},
		},
		{
			name: "whitespace handling",
			expr: " name = John Doe ",
			want: Condition{Path: "name", Operator: OpEqual, Value: "John Doe"},
		},
		{
			name:        "empty expression",
			expr:        "",
			wantErr:     true,
			errContains: "empty filter expression",
		},
		{
			name:        "no operator",
			expr:        "nameJohn",
			wantErr:     true,
			errContains: "no valid operator found",
		},
		{
			name:        "invalid regex",
			expr:        "name=~[invalid",
			wantErr:     true,
			errContains: "invalid regex pattern",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCondition(tt.expr)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseCondition() expected error but got none")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ParseCondition() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseCondition() error = %v, want nil", err)
				return
			}
			if got.Path != tt.want.Path {
				t.Errorf("ParseCondition() Path = %v, want %v", got.Path, tt.want.Path)
			}
			if got.Operator != tt.want.Operator {
				t.Errorf("ParseCondition() Operator = %v, want %v", got.Operator, tt.want.Operator)
			}
			if got.Value != tt.want.Value {
				t.Errorf("ParseCondition() Value = %v, want %v", got.Value, tt.want.Value)
			}
			// Check regex compilation for match operators
			if (tt.want.Operator == OpMatch || tt.want.Operator == OpNotMatch) && got.regex == nil {
				t.Errorf("ParseCondition() regex should be compiled for match operators")
			}
		})
	}
}

func TestParseConditions(t *testing.T) {
	tests := []struct {
		name        string
		exprs       []string
		wantCount   int
		wantErr     bool
		errContains string
	}{
		{
			name:      "multiple conditions",
			exprs:     []string{"name=John", "age>25", "status!=inactive"},
			wantCount: 3,
		},
		{
			name:      "with empty strings",
			exprs:     []string{"name=John", "", "age>25", "   "},
			wantCount: 2,
		},
		{
			name:      "empty input",
			exprs:     []string{},
			wantCount: 0,
		},
		{
			name:        "invalid condition",
			exprs:       []string{"name=John", "invalid"},
			wantErr:     true,
			errContains: "no valid operator found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseConditions(tt.exprs)
			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseConditions() expected error but got none")
					return
				}
				if tt.errContains != "" && !containsString(err.Error(), tt.errContains) {
					t.Errorf("ParseConditions() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("ParseConditions() error = %v, want nil", err)
				return
			}
			if len(got) != tt.wantCount {
				t.Errorf("ParseConditions() count = %v, want %v", len(got), tt.wantCount)
			}
		})
	}
}

func TestFilter_Apply(t *testing.T) {
	// Create test data
	rows := []flatten.FlatKV{
		{
			"name":   "John",
			"age":    30,
			"active": true,
			"score":  85.5,
		},
		{
			"name":   "Jane",
			"age":    25,
			"active": false,
			"score":  92.0,
		},
		{
			"name":   "Bob",
			"age":    35,
			"active": true,
			"score":  78.0,
		},
		{
			"name":   "Alice",
			"age":    28,
			"active": true,
			"score":  95.5,
		},
	}

	tests := []struct {
		name       string
		conditions []string
		wantCount  int
		wantNames  []string
	}{
		{
			name:       "no conditions",
			conditions: []string{},
			wantCount:  4,
			wantNames:  []string{"John", "Jane", "Bob", "Alice"},
		},
		{
			name:       "equality filter",
			conditions: []string{"name=John"},
			wantCount:  1,
			wantNames:  []string{"John"},
		},
		{
			name:       "numeric greater than",
			conditions: []string{"age>30"},
			wantCount:  1,
			wantNames:  []string{"Bob"},
		},
		{
			name:       "numeric greater than equal",
			conditions: []string{"age>=30"},
			wantCount:  2,
			wantNames:  []string{"John", "Bob"},
		},
		{
			name:       "boolean filter",
			conditions: []string{"active=true"},
			wantCount:  3,
			wantNames:  []string{"John", "Bob", "Alice"},
		},
		{
			name:       "multiple conditions (AND)",
			conditions: []string{"active=true", "age>=30"},
			wantCount:  2,
			wantNames:  []string{"John", "Bob"},
		},
		{
			name:       "float comparison",
			conditions: []string{"score>90"},
			wantCount:  2,
			wantNames:  []string{"Jane", "Alice"},
		},
		{
			name:       "not equal",
			conditions: []string{"name!=John"},
			wantCount:  3,
			wantNames:  []string{"Jane", "Bob", "Alice"},
		},
		{
			name:       "contains",
			conditions: []string{"name~o"},
			wantCount:  2,
			wantNames:  []string{"John", "Bob"},
		},
		{
			name:       "less than",
			conditions: []string{"age<30"},
			wantCount:  2,
			wantNames:  []string{"Jane", "Alice"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions, err := ParseConditions(tt.conditions)
			if err != nil {
				t.Fatalf("ParseConditions() error = %v", err)
			}

			filter := NewFilter(conditions)
			result := filter.Apply(rows)

			if len(result) != tt.wantCount {
				t.Errorf("Filter.Apply() count = %v, want %v", len(result), tt.wantCount)
			}

			// Check that the right rows are returned
			gotNames := make([]string, len(result))
			for i, row := range result {
				gotNames[i] = row["name"].(string)
			}

			if !slicesEqual(gotNames, tt.wantNames) {
				t.Errorf("Filter.Apply() names = %v, want %v", gotNames, tt.wantNames)
			}
		})
	}
}

func TestFilter_MissingFields(t *testing.T) {
	rows := []flatten.FlatKV{
		{
			"name": "John",
			"age":  30,
		},
		{
			"name": "Jane",
			// age missing
		},
	}

	tests := []struct {
		name      string
		condition string
		wantCount int
		wantNames []string
	}{
		{
			name:      "equal to missing field",
			condition: "age=",
			wantCount: 1,
			wantNames: []string{"Jane"},
		},
		{
			name:      "not equal to missing field",
			condition: "age!=",
			wantCount: 1,
			wantNames: []string{"John"},
		},
		{
			name:      "numeric comparison with missing field",
			condition: "age>25",
			wantCount: 1,
			wantNames: []string{"John"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions, err := ParseConditions([]string{tt.condition})
			if err != nil {
				t.Fatalf("ParseConditions() error = %v", err)
			}

			filter := NewFilter(conditions)
			result := filter.Apply(rows)

			if len(result) != tt.wantCount {
				t.Errorf("Filter.Apply() count = %v, want %v", len(result), tt.wantCount)
			}

			gotNames := make([]string, len(result))
			for i, row := range result {
				gotNames[i] = row["name"].(string)
			}

			if !slicesEqual(gotNames, tt.wantNames) {
				t.Errorf("Filter.Apply() names = %v, want %v", gotNames, tt.wantNames)
			}
		})
	}
}

func TestFilter_RegexMatching(t *testing.T) {
	rows := []flatten.FlatKV{
		{
			"email": "john@example.com",
		},
		{
			"email": "jane@test.org",
		},
		{
			"email": "bob@example.net",
		},
	}

	tests := []struct {
		name      string
		condition string
		wantCount int
	}{
		{
			name:      "regex match",
			condition: "email=~.*@example\\.(com|net)",
			wantCount: 2,
		},
		{
			name:      "regex not match",
			condition: "email!=~.*@example\\.com",
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conditions, err := ParseConditions([]string{tt.condition})
			if err != nil {
				t.Fatalf("ParseConditions() error = %v", err)
			}

			filter := NewFilter(conditions)
			result := filter.Apply(rows)

			if len(result) != tt.wantCount {
				t.Errorf("Filter.Apply() count = %v, want %v", len(result), tt.wantCount)
			}
		})
	}
}

func TestOperator_String(t *testing.T) {
	tests := []struct {
		op   Operator
		want string
	}{
		{OpEqual, "="},
		{OpNotEqual, "!="},
		{OpGreaterThan, ">"},
		{OpGreaterThanEqual, ">="},
		{OpLessThan, "<"},
		{OpLessThanEqual, "<="},
		{OpContains, "~"},
		{OpNotContains, "!~"},
		{OpMatch, "=~"},
		{OpNotMatch, "!=~"},
		{Operator(999), "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := tt.op.String()
			if got != tt.want {
				t.Errorf("Operator.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper functions
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr) >= 0
}

func findSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func slicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
