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

func main() {
	var config app.Config

	root := &cobra.Command{
		Use:     "tablo",
		Version: version,
		Short:   "Render JSON/YAML as tables",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runApp(&config)
		},
	}

	// input
	root.Flags().StringVarP(&config.Input.File, "file", "f", "", "Path to input file")
	root.Flags().StringVarP(&config.Input.String, "input", "i", "", "Raw input string")
	root.Flags().StringVarP(&config.Input.Format, "format", "F", "auto", "Input format: auto|json|yaml|yml")

	// flatten
	root.Flags().BoolVarP(&config.Flatten.Enabled, "dive", "d", false, "Enable flattening of nested objects and arrays of objects")
	root.Flags().StringSliceVarP(&config.Flatten.Paths, "dive-path", "D", nil, "Dive only into listed top-level paths (repeatable)")
	root.Flags().IntVarP(&config.Flatten.MaxDepth, "max-depth", "m", -1, "Maximum depth to dive; -1 = unlimited")
	root.Flags().BoolVar(&config.Flatten.FlattenSimpleArray, "flatten-simple-arrays", false, "Flatten arrays of primitives to comma-separated strings")

	// selection
	root.Flags().StringVarP(&config.Selection.SelectExpr, "select", "s", "", "Comma-separated dotted path expressions to include")
	root.Flags().StringVar(&config.Selection.SelectFile, "select-file", "", "Path to file containing one path expression per line")
	root.Flags().StringVarP(&config.Selection.ExcludeExpr, "exclude", "E", "", "Comma-separated dotted path expressions to exclude")
	root.Flags().BoolVar(&config.Selection.StrictSelect, "strict-select", false, "Error when any selected path does not exist")

	// output formatting
	root.Flags().StringVar(&config.Output.Style, "style", "heavy", "Table style: heavy|light|double|ascii|markdown|compact|borderless|html|csv")
	root.Flags().BoolVar(&config.Output.ASCIIOnly, "ascii", false, "Force ASCII borders")
	root.Flags().BoolVar(&config.Output.NoHeader, "no-header", false, "Omit header row")
	root.Flags().StringVar(&config.Output.HeaderCase, "header-case", "original", "Header case: original|upper|lower|title")
	root.Flags().IntVar(&config.Output.MaxColWidth, "max-col-width", 0, "Max column width; 0 = no limit")
	root.Flags().StringVar(&config.Output.WrapMode, "wrap", "off", "Cell wrapping: off|word|char")
	root.Flags().StringVar(&config.Output.TruncateSuffix, "truncate-suffix", "…", "Suffix for truncated cells")
	root.Flags().StringVar(&config.Output.NullStr, "null-str", "null", "String to represent null values")
	root.Flags().StringVar(&config.Output.BoolStr, "bool-str", "true:false", "Booleans mapping true:false")
	root.Flags().IntVar(&config.Output.Precision, "precision", -1, "Decimal precision for floats; -1 = no change")
	root.Flags().StringVarP(&config.Output.FilePath, "output", "o", "", "Write output to file instead of stdout")
	root.Flags().BoolVar(&config.Output.IndexColumn, "index-column", false, "Include INDEX column for arrays")
	root.Flags().IntVar(&config.Output.Limit, "limit", 0, "Limit number of rows printed; 0 = all")
	root.Flags().StringVar(&config.Output.Color, "color", "auto", "Colorize output: auto|always|never")

	// general
	root.Flags().BoolVar(&config.General.Quiet, "quiet", false, "Suppress non-error logging")

	root.SetVersionTemplate("{{.Version}}\n")

	if err := root.Execute(); err != nil {
		handleError(err, config.General.Quiet)
	}
}

// runApp executes the main application logic
func runApp(config *app.Config) error {
	application := app.New(*config, os.Stdin)
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
