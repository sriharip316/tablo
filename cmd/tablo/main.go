package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/sriharip316/tablo/internal/app"
)

// version is injected at build time using:
//
//	go build -ldflags="-X main.version=v1.2.3"
var version string = "dev"

// Legacy options struct for CLI flag binding
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

// toAppConfig converts CLI options to application configuration
func (opts *options) toAppConfig() app.Config {
	return app.Config{
		Input: app.InputConfig{
			File:   opts.file,
			String: opts.inStr,
			Format: opts.format,
		},
		Flatten: app.FlattenConfig{
			Enabled:            opts.dive,
			Paths:              opts.divePaths,
			MaxDepth:           opts.maxDepth,
			FlattenSimpleArray: opts.flatSimple,
		},
		Selection: app.SelectionConfig{
			SelectExpr:   opts.selectExpr,
			SelectFile:   opts.selectFile,
			ExcludeExpr:  opts.exclude,
			StrictSelect: opts.strictSel,
		},
		Output: app.OutputConfig{
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
			FilePath:       opts.outputPath,
			IndexColumn:    opts.indexColumn,
			Limit:          opts.limit,
			Color:          opts.color,
		},
		General: app.GeneralConfig{
			Quiet: opts.quiet,
		},
	}
}

func main() {
	var opts options

	root := &cobra.Command{
		Use:     "tablo",
		Version: version,
		Short:   "Render JSON/YAML as tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApp(&opts)
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

	root.SetVersionTemplate("{{.Version}}\n")

	if err := root.Execute(); err != nil {
		handleError(err, opts.quiet)
	}
}

// runApp executes the main application logic
func runApp(opts *options) error {
	config := opts.toAppConfig()
	application := app.New(config, os.Stdin)
	return application.Run()
}

// handleError processes application errors and exits with appropriate codes
func handleError(err error, quiet bool) {
	exitCode := app.GetExitCode(err)

	if !quiet {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
	}

	os.Exit(exitCode)
}
