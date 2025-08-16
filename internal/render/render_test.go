package render

import (
	"testing"

	"github.com/sriharip316/tablo/internal/flatten"
)

func TestRenderObjectKV(t *testing.T) {
	kv := flatten.FlatKV{"a": 1, "b.c": 2}
	m := Model{Mode: ModeObjectKV, KV: kv}
	out, err := Render(m, Options{Style: "ascii"})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) == 0 {
		t.Fatal("empty output")
	}
}
