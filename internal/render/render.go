package render

import (
	stdjson "encoding/json"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/sriharip316/tablo/internal/flatten"
)

type Mode int

const (
	ModeRows Mode = iota
	ModeObjectKV
)

// Model represents a normalized table.
type Model struct {
	Mode Mode
	// for ModeRows
	Headers []string
	Rows    [][]any
	// for ModeObjectKV
	KV      flatten.FlatKV
	KVOrder []string
	// extra
	IndexColumn bool
}

type Options struct {
	Style          string
	ASCIIOnly      bool
	NoHeader       bool
	HeaderCase     string
	MaxColWidth    int
	WrapMode       string
	TruncateSuffix string
	NullStr        string
	BoolStr        string // format true:false
	Precision      int
	Color          string // auto|always|never (noop for now)
}

func FromPrimitiveArray(arr []any, index bool, limit int) Model {
	rows := make([][]any, 0, len(arr))
	n := len(arr)
	if limit > 0 && n > limit {
		n = limit
	}
	for i := 0; i < n; i++ {
		v := arr[i]
		rows = append(rows, []any{v})
	}
	headers := []string{"VALUE"}
	return Model{Mode: ModeRows, Headers: headers, Rows: rows, IndexColumn: index}
}

func FromFlatRows(rows []flatten.FlatKV, headers []string, index bool) Model {
	data := make([][]any, len(rows))
	for i, r := range rows {
		row := make([]any, len(headers))
		for j, h := range headers {
			if val, ok := r[h]; ok {
				row[j] = val
			} else {
				row[j] = nil
			}
		}
		data[i] = row
	}
	return Model{Mode: ModeRows, Headers: headers, Rows: data, IndexColumn: index}
}

func Render(m Model, o Options) (string, error) {
	t := table.NewWriter()

	// style
	t.SetStyle(resolveStyle(o))

	// headers
	if m.Mode == ModeObjectKV {
		// turn KV into two columns: KEY, VALUE
		headers := []string{"KEY", "VALUE"}
		if o.NoHeader {
			// no headers rendered
		} else {
			t.AppendHeader(toHeaderRow(headers, o))
		}
		// column configs for ModeObjectKV
		colCfgs := make([]table.ColumnConfig, 0, len(headers))
		for i, h := range headers {
			cfg := table.ColumnConfig{Name: headerCase(h, o.HeaderCase)}
			if o.MaxColWidth > 0 {
				cfg.WidthMax = o.MaxColWidth
				cfg.WidthMaxEnforcer = wrapEnforcer(o)
			}
			cfg.Number = i + 1
			colCfgs = append(colCfgs, cfg)
		}
		t.SetColumnConfigs(colCfgs)

		keys := m.KVOrder
		if len(keys) == 0 {
			keys = m.KV.Keys()
		}
		for _, k := range keys {
			v := m.KV[k]
			t.AppendRow(table.Row{k, formatCell(v, o)})
		}
		return chooseRender(t, o), nil
	}

	// ModeRows
	headers := m.Headers
	if !o.NoHeader {
		t.AppendHeader(toHeaderRow(headers, o))
	}
	// column configs
	colCfgs := make([]table.ColumnConfig, 0, len(headers))
	for i, h := range headers {
		cfg := table.ColumnConfig{Name: headerCase(h, o.HeaderCase)}
		if o.MaxColWidth > 0 {
			cfg.WidthMax = o.MaxColWidth
			cfg.WidthMaxEnforcer = wrapEnforcer(o)
		}
		// set numeric align right by default; handled by library if using SetAutoIndex?
		cfg.Number = i + 1
		colCfgs = append(colCfgs, cfg)
	}
	t.SetColumnConfigs(colCfgs)

	if m.IndexColumn {
		t.SetAutoIndex(true)
		t.SetIndexColumn(1)
	}

	for _, r := range m.Rows {
		row := make(table.Row, len(r))
		for i := range r {
			row[i] = formatCell(r[i], o)
		}
		t.AppendRow(row)
	}

	return chooseRender(t, o), nil
}

func toHeaderRow(headers []string, o Options) table.Row {
	hr := make(table.Row, len(headers))
	for i, h := range headers {
		hr[i] = headerCase(h, o.HeaderCase)
	}
	return hr
}

func headerCase(h, mode string) string {
	switch strings.ToLower(mode) {
	case "upper":
		return strings.ToUpper(h)
	case "lower":
		return strings.ToLower(h)
	case "title":
		caser := cases.Title(language.Und)
		return caser.String(strings.ReplaceAll(h, "_", " "))
	default:
		return h
	}
}

func chooseRender(t table.Writer, o Options) string {
	switch strings.ToLower(o.Style) {
	case "markdown":
		return t.RenderMarkdown()
	default:
		return t.Render()
	}
}

func resolveStyle(o Options) table.Style {
	// default style mapping; ASCII forcing handled by StyleDefault/Box
	s := table.StyleDefault
	switch strings.ToLower(o.Style) {
	case "heavy":
		s = table.StyleBold
	case "light":
		s = table.StyleLight
	case "double":
		s = table.StyleDouble
	case "ascii":
		s = table.StyleDefault
	case "markdown":
		// style doesn't matter; choose simple
		s = table.StyleDefault
	case "compact":
		s = table.StyleLight
		s.Options.SeparateRows = false
	case "borderless":
		s = table.StyleLight
		s.Options = table.OptionsNoBordersAndSeparators
	}
	if o.ASCIIOnly {
		// force ASCII box drawing
		s.Box = table.StyleBoxDefault
	}
	// header formatting
	switch strings.ToLower(o.HeaderCase) {
	case "upper":
		s.Format.Header = text.FormatUpper
	case "lower":
		s.Format.Header = text.FormatLower
	case "title":
		s.Format.Header = text.FormatTitle
	default:
		s.Format.Header = text.FormatDefault
	}
	// width handled by column configs
	return s
}

func wrapEnforcer(o Options) table.WidthEnforcer {
	// To improve test coverage, a test case should be added in render_test.go
	// to specifically cover the scenario where o.TruncateSuffix is longer
	// than the available column width. This ensures the `suf = ""` line is hit.
	switch strings.ToLower(o.WrapMode) {
	case "char":
		return text.WrapText
	case "word":
		return text.WrapSoft
	case "off":
		fallthrough
	default:
		// truncate with suffix
		return func(s string, width int) string {
			if width <= 0 {
				return s
			}
			r := []rune(s)
			if len(r) <= width {
				return s
			}
			suf := o.TruncateSuffix
			if len([]rune(suf)) > width {
				suf = ""
			}
			keep := width - len([]rune(suf))
			keep = max(keep, 0)
			return string(r[:keep]) + suf
		}
	}
}

func formatCell(v any, o Options) any {
	if v == nil {
		if o.NullStr != "" {
			return o.NullStr
		}
		return nil
	}
	switch t := v.(type) {
	case bool:
		if o.BoolStr != "" && strings.Contains(o.BoolStr, ":") {
			parts := strings.SplitN(o.BoolStr, ":", 2)
			if t {
				return parts[0]
			}
			return parts[1]
		}
		return t
	case float64:
		if o.Precision >= 0 {
			// avoid -0
			if t == 0 {
				t = 0
			}
			return strconv.FormatFloat(t, 'f', o.Precision, 64)
		}
		return t
	case float32:
		if o.Precision >= 0 {
			f := float64(t)
			if f == 0 {
				f = 0
			}
			return strconv.FormatFloat(f, 'f', o.Precision, 64)
		}
		return t
	case stdjson.Number:
		if o.Precision >= 0 {
			if f, err := t.Float64(); err == nil {
				return strconv.FormatFloat(f, 'f', o.Precision, 64)
			}
		}
		return t.String()
	default:
		return t
	}
}
