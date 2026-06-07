package filter

import (
	"testing"
)

func BenchmarkParseConditionRegex(b *testing.B) {
	expr := "email=~^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseCondition(expr)
		if err != nil {
			b.Fatal(err)
		}
	}
}
