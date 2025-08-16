package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/sriharip316/tablo/internal/flatten"
	"github.com/sriharip316/tablo/internal/input"
	"github.com/sriharip316/tablo/internal/parse"
	"github.com/sriharip316/tablo/internal/render"
	"github.com/sriharip316/tablo/internal/selectors"
)

var version = "0.1.0"

type options struct {
	// input
	file   string
	inStr  string
	format string

	// flatten
	dive       bool
	divePaths  []string
	maxDepth   int
	flatSimple bool

	// selection
	selectExpr string
	selectFile string
	exclude    string
	strictSel  bool

	// output formatting
	style          string
	asciiOnly      bool
	noHeader       bool
	headerCase     string
	maxColWidth    int
	wrap           string
	truncateSuffix string
	nullStr        string
	boolStr        string
	precision      int
	outputPath     string
	indexColumn    bool
	limit          int
	color          string

	// general
	quiet bool
}

func main() {
	var opts options

	root := &cobra.Command{
		Use:   "tablo",
		Short: "Render JSON/YAML as tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			// input precedence checks
			if opts.inStr != "" && opts.file != "" {
				return cliErr(2, "conflicting inputs: --input and --file")
			}
			// read
			rdr := input.NewReader(opts.inStr, opts.file, os.Stdin)
			data, rerr := rdr.Read()
			if rerr != nil {
				return cliErr(3, rerr.Error())
			}

			// detect & parse
			det := parse.Detector{Explicit: opts.format, FilePath: opts.file}
			fmtKind := det.Detect(data)
			parsed, perr := parse.Parse(data, fmtKind)
			if perr != nil {
				return cliErr(4, perr.Error())
			}

			// for YAML multi-docs, Parse already returns []any
			// build rows/columns
			var singleObject bool

			switch v := parsed.(type) {
			case []any:
				_ = v
			case map[string]any:
				singleObject = true
			case map[any]any:
				singleObject = true
				v2 := parse.ToStringKeyMap(v)
				parsed = v2
			default:
				vrows := []any{v}
				parsed = vrows
			}

			// flatten if requested
			fopts := flatten.Options{
				Enabled:            opts.dive || len(opts.divePaths) > 0,
				MaxDepth:           opts.maxDepth,
				DivePaths:          opts.divePaths,
				FlattenSimpleArray: opts.flatSimple,
			}

			var tableModel render.Model

			if singleObject {
				// object → key/value table
				flat := flatten.FlattenObject(parsed, fopts)
				// selection on keys
				inc, exc, selErr := compileSelections(opts)
				if selErr != nil {
					return cliErr(2, selErr.Error())
				}
				keys := selectors.ApplyToKeys(flat.Keys(), inc, exc)
				if opts.strictSel {
					missing := selectors.MissingExpressions(keys, inc)
					if len(missing) > 0 {
						return cliErr(5, "missing selected paths: "+strings.Join(missing, ", "))
					}
				}
				tableModel = render.Model{Mode: render.ModeObjectKV, KV: flat, KVOrder: keys}
			} else {
				// array or others
				switch v := parsed.(type) {
				case []any:
					// either array of maps or primitives
					isObjects := parse.ArrayIsObjects(v)
					if !isObjects {
						// primitives
						tableModel = render.FromPrimitiveArray(v, opts.indexColumn, opts.limit)
					} else {
						// objects
						flatRows := flatten.FlattenRows(v, fopts)
						inc, exc, selErr := compileSelections(opts)
						if selErr != nil {
							return cliErr(2, selErr.Error())
						}
						hdrs := selectors.HeadersUnion(flatRows)
						if len(inc) > 0 {
							hdrs = selectors.ApplyToKeys(hdrs, inc, nil)
						}
						if len(exc) > 0 {
							hdrs = selectors.ApplyToKeys(hdrs, nil, exc)
						}
						if opts.strictSel && len(inc) > 0 {
							missing := selectors.MissingExpressions(hdrs, inc)
							if len(missing) > 0 {
								return cliErr(5, "missing selected paths: "+strings.Join(missing, ", "))
							}
						}
						if opts.limit > 0 && len(flatRows) > opts.limit {
							flatRows = flatRows[:opts.limit]
						}
						tableModel = render.FromFlatRows(flatRows, hdrs, opts.indexColumn)
					}
				default:
					// scalar → same as primitives array of length 1
					tableModel = render.FromPrimitiveArray([]any{parsed}, opts.indexColumn, opts.limit)
				}
			}

			// render
			ro := render.Options{
				Style:          opts.style,
				ASCIIOnly:      opts.asciiOnly,
				NoHeader:       opts.noHeader,
				HeaderCase:     opts.headerCase,
				MaxColWidth:    opts.maxColWidth,
				WrapMode:       opts.wrap,
				TruncateSuffix: opts.truncateSuffix,
				NullStr:        opts.nullStr,
				BoolStr:        opts.boolStr,
				Precision:      opts.precision,
				Color:          opts.color,
			}
			out, rerr := render.Render(tableModel, ro)
			if rerr != nil {
				return cliErr(2, rerr.Error())
			}

			// write output
			var w io.Writer = os.Stdout
			if opts.outputPath != "" {
				f, ferr := os.Create(opts.outputPath)
				if ferr != nil {
					return cliErr(4, ferr.Error())
				}
				defer f.Close()
				w = f
			}
			// ensure a trailing newline after the table
			if !strings.HasSuffix(out, "\n") {
				out += "\n"
			}
			_, _ = io.WriteString(w, out)
			return nil
		},
	}

	// flags
	// input
	root.Flags().StringVarP(&opts.file, "file", "f", "", "Path to input file")
	root.Flags().StringVarP(&opts.inStr, "input", "i", "", "Raw input string")
	root.Flags().StringVarP(&opts.format, "format", "F", "auto", "Input format: auto|json|yaml|yml")

	// flatten
	root.Flags().BoolVarP(&opts.dive, "dive", "d", false, "Enable flattening of nested objects and arrays of objects")
	root.Flags().StringSliceVarP(&opts.divePaths, "dive-path", "D", nil, "Dive only into listed top-level paths (repeatable)")
	root.Flags().IntVarP(&opts.maxDepth, "max-depth", "m", -1, "Maximum depth to dive; -1 = unlimited")
	root.Flags().BoolVar(&opts.flatSimple, "flatten-simple-arrays", false, "Flatten arrays of primitives to comma-separated strings")

	// selection
	root.Flags().StringVarP(&opts.selectExpr, "select", "s", "", "Comma-separated dotted path expressions to include")
	root.Flags().StringVar(&opts.selectFile, "select-file", "", "Path to file containing one path expression per line")
	root.Flags().StringVarP(&opts.exclude, "exclude", "E", "", "Comma-separated dotted path expressions to exclude")
	root.Flags().BoolVar(&opts.strictSel, "strict-select", false, "Error when any selected path does not exist")

	// output formatting
	root.Flags().StringVar(&opts.style, "style", "heavy", "Table style: heavy|light|double|ascii|markdown|compact|borderless")
	root.Flags().BoolVar(&opts.asciiOnly, "ascii", false, "Force ASCII borders")
	root.Flags().BoolVar(&opts.noHeader, "no-header", false, "Omit header row")
	root.Flags().StringVar(&opts.headerCase, "header-case", "original", "Header case: original|upper|lower|title")
	root.Flags().IntVar(&opts.maxColWidth, "max-col-width", 0, "Max column width; 0 = no limit")
	root.Flags().StringVar(&opts.wrap, "wrap", "off", "Cell wrapping: off|word|char")
	root.Flags().StringVar(&opts.truncateSuffix, "truncate-suffix", "…", "Suffix for truncated cells")
	root.Flags().StringVar(&opts.nullStr, "null-str", "null", "String to represent null values")
	root.Flags().StringVar(&opts.boolStr, "bool-str", "true:false", "Booleans mapping true:false")
	root.Flags().IntVar(&opts.precision, "precision", -1, "Decimal precision for floats; -1 = no change")
	root.Flags().StringVarP(&opts.outputPath, "output", "o", "", "Write output to file instead of stdout")
	root.Flags().BoolVar(&opts.indexColumn, "index-column", false, "Include INDEX column for arrays")
	root.Flags().IntVar(&opts.limit, "limit", 0, "Limit number of rows printed; 0 = all")
	root.Flags().StringVar(&opts.color, "color", "auto", "Colorize output: auto|always|never")

	// general
	root.Flags().BoolVar(&opts.quiet, "quiet", false, "Suppress non-error logging")

	root.Version = version
	root.SetVersionTemplate("{{.Version}}\n")

	if err := root.Execute(); err != nil {
		var ce *cliError
		if errors.As(err, &ce) {
			if !opts.quiet {
				_, _ = fmt.Fprintln(os.Stderr, ce.msg)
			}
			os.Exit(ce.code)
		}
		if !opts.quiet {
			_, _ = fmt.Fprintln(os.Stderr, err.Error())
		}
		os.Exit(1)
	}
}

func compileSelections(opts options) (include []selectors.Expr, exclude []selectors.Expr, err error) {
	var inc []string
	if opts.selectExpr != "" {
		inc = append(inc, splitComma(opts.selectExpr)...)
	}
	if opts.selectFile != "" {
		b, rerr := os.ReadFile(opts.selectFile)
		if rerr != nil {
			return nil, nil, rerr
		}
		for _, line := range strings.Split(string(b), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			inc = append(inc, line)
		}
	}
	var exc []string
	if opts.exclude != "" {
		exc = append(exc, splitComma(opts.exclude)...)
	}
	incx, ierr := selectors.CompileMany(inc)
	if ierr != nil {
		return nil, nil, ierr
	}
	excx, eerr := selectors.CompileMany(exc)
	if eerr != nil {
		return nil, nil, eerr
	}
	return incx, excx, nil
}

func splitComma(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	return out
}

type cliError struct {
	code int
	msg  string
}

func (e *cliError) Error() string { return e.msg }

func cliErr(code int, msg string) error { return &cliError{code: code, msg: msg} }
