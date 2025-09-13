package app

import (
	"io"
	"os"

	"github.com/sriharip316/tablo/internal/filter"
	"github.com/sriharip316/tablo/internal/flatten"
	"github.com/sriharip316/tablo/internal/input"
	"github.com/sriharip316/tablo/internal/parse"
	"github.com/sriharip316/tablo/internal/render"
	"github.com/sriharip316/tablo/internal/selectors"
	"github.com/sriharip316/tablo/internal/sort"
)

// Config groups application configuration into logical sections
type Config struct {
	Input     InputConfig
	Flatten   FlattenConfig
	Selection SelectionConfig
	Filter    FilterConfig
	Sort      SortConfig
	Output    OutputConfig
	General   GeneralConfig
}

type InputConfig struct {
	File        string
	String      string
	Format      string
	CSVNoHeader bool
}

type FlattenConfig struct {
	Enabled            bool
	Paths              []string
	MaxDepth           int
	FlattenSimpleArray bool
}

type SelectionConfig struct {
	SelectExpr   string
	SelectFile   string
	ExcludeExpr  string
	StrictSelect bool
}

type FilterConfig struct {
	WhereExprs []string
}

type SortConfig struct {
	Columns []string
}

type OutputConfig struct {
	Style          string
	ASCIIOnly      bool
	NoHeader       bool
	HeaderCase     string
	MaxColWidth    int
	WrapMode       string
	TruncateSuffix string
	NullStr        string
	BoolStr        string
	Precision      int
	FilePath       string
	IndexColumn    bool
	Limit          int
	Color          string
}

type GeneralConfig struct {
	Quiet bool
}

// Application encapsulates the core application logic
type Application struct {
	config Config
	stdin  io.Reader
}

// New creates a new Application instance
func New(config Config, stdin io.Reader) *Application {
	return &Application{
		config: config,
		stdin:  stdin,
	}
}

// Run executes the main application logic
func (app *Application) Run() error {
	// Validate configuration
	if err := app.validateConfig(); err != nil {
		return err
	}

	// Read input
	data, err := app.readInput()
	if err != nil {
		return NewError(ErrCodeInput, "failed to read input", err)
	}

	// Parse data
	parsed, err := app.parseData(data)
	if err != nil {
		return NewError(ErrCodeParse, "failed to parse input", err)
	}

	// Process data (flatten, select, etc.)
	model, err := app.processData(parsed)
	if err != nil {
		return NewError(ErrCodeProcessing, "failed to process data", err)
	}

	// Render output
	output, err := app.renderOutput(model)
	if err != nil {
		return NewError(ErrCodeRender, "failed to render output", err)
	}

	// Write output
	if err := app.writeOutput(output); err != nil {
		return NewError(ErrCodeOutput, "failed to write output", err)
	}

	return nil
}

func (app *Application) validateConfig() error {
	if app.config.Input.String != "" && app.config.Input.File != "" {
		return NewError(ErrCodeUsage, "conflicting inputs: --input and --file cannot be used together", nil)
	}
	return nil
}

func (app *Application) readInput() ([]byte, error) {
	reader := input.NewReader(app.config.Input.String, app.config.Input.File, app.stdin)
	return reader.Read()
}

func (app *Application) parseData(data []byte) (any, error) {
	detector := parse.Detector{
		Explicit: app.config.Input.Format,
		FilePath: app.config.Input.File,
	}
	format := detector.Detect(data)
	opts := parse.ParseOptions{
		CSVNoHeader: app.config.Input.CSVNoHeader,
	}
	return parse.Parse(data, format, opts)
}

func (app *Application) processData(parsed any) (render.Model, error) {
	// Normalize data structure
	normalized := app.normalizeData(parsed)

	// Apply flattening
	flattenOpts := flatten.Options{
		Enabled:            app.config.Flatten.Enabled || len(app.config.Flatten.Paths) > 0,
		MaxDepth:           app.config.Flatten.MaxDepth,
		DivePaths:          app.config.Flatten.Paths,
		FlattenSimpleArray: app.config.Flatten.FlattenSimpleArray,
	}

	// Determine processing mode based on data structure
	switch data := normalized.(type) {
	case map[string]any:
		return app.processObject(data, flattenOpts)
	case []any:
		return app.processArray(data, flattenOpts)
	default:
		// Single primitive value
		return render.FromPrimitiveArray([]any{data}, app.config.Output.IndexColumn, app.config.Output.Limit), nil
	}
}

func (app *Application) normalizeData(parsed any) any {
	switch v := parsed.(type) {
	case []any:
		for i, item := range v {
			v[i] = app.normalizeData(item)
		}
		return v
	case map[string]any:
		for key, val := range v {
			v[key] = app.normalizeData(val)
		}
		return v
	case map[any]any:
		stringMap := parse.ToStringKeyMap(v)
		for key, val := range stringMap {
			stringMap[key] = app.normalizeData(val)
		}
		return stringMap
	case []map[string]any:
		// Convert []map[string]any to []any for CSV support
		result := make([]any, len(v))
		for i, m := range v {
			result[i] = app.normalizeData(m)
		}
		return result
	default:
		return v
	}
}

func (app *Application) processObject(obj map[string]any, flattenOpts flatten.Options) (render.Model, error) {
	flattened := flatten.FlattenObject(obj, flattenOpts)

	// Apply selection
	keys, err := app.applySelection(flattened.Keys())
	if err != nil {
		return render.Model{}, err
	}

	return render.Model{
		Mode:    render.ModeObjectKV,
		KV:      flattened,
		KVOrder: keys,
	}, nil
}

func (app *Application) processArray(arr []any, flattenOpts flatten.Options) (render.Model, error) {
	// Check if array contains objects
	if !parse.ArrayIsObjects(arr) {
		return render.FromPrimitiveArray(arr, app.config.Output.IndexColumn, app.config.Output.Limit), nil
	}

	// Process array of objects
	flatRows := flatten.FlattenRows(arr, flattenOpts)

	// Apply row filtering
	filteredRows, err := app.applyRowFiltering(flatRows)
	if err != nil {
		return render.Model{}, err
	}

	// Apply sorting
	sortedRows := app.applySorting(filteredRows)

	// Apply limit after filtering and sorting
	if app.config.Output.Limit > 0 && len(sortedRows) > app.config.Output.Limit {
		sortedRows = sortedRows[:app.config.Output.Limit]
	}

	// Get union of headers
	headers := selectors.HeadersUnion(filteredRows)

	// Apply selection
	filteredHeaders, err := app.applySelection(headers)
	if err != nil {
		return render.Model{}, err
	}

	// If limit is 1, treat it as single object
	if app.config.Output.Limit == 1 {
		return render.Model{
			Mode:    render.ModeObjectKV,
			KV:      sortedRows[0],
			KVOrder: filteredHeaders,
		}, nil
	}

	return render.FromFlatRows(sortedRows, filteredHeaders, app.config.Output.IndexColumn), nil
}

func (app *Application) applySelection(keys []string) ([]string, error) {
	include, exclude, err := app.compileSelectors()
	if err != nil {
		return nil, err
	}

	filtered := selectors.ApplyToKeys(keys, include, exclude)

	// Check strict selection requirements
	if app.config.Selection.StrictSelect && len(include) > 0 {
		missing := selectors.MissingExpressions(filtered, include)
		if len(missing) > 0 {
			return nil, NewError(ErrCodeSelection, "missing selected paths: "+joinStrings(missing, ", "), nil)
		}
	}

	return filtered, nil
}

func (app *Application) compileSelectors() (include, exclude []selectors.Expr, err error) {
	// Compile include selectors
	var includePatterns []string
	if app.config.Selection.SelectExpr != "" {
		includePatterns = append(includePatterns, splitCommaString(app.config.Selection.SelectExpr)...)
	}
	if app.config.Selection.SelectFile != "" {
		filePatterns, err := app.readSelectFile()
		if err != nil {
			return nil, nil, err
		}
		includePatterns = append(includePatterns, filePatterns...)
	}

	include, err = selectors.CompileMany(includePatterns)
	if err != nil {
		return nil, nil, NewError(ErrCodeUsage, "invalid include selector", err)
	}

	// Compile exclude selectors
	var excludePatterns []string
	if app.config.Selection.ExcludeExpr != "" {
		excludePatterns = splitCommaString(app.config.Selection.ExcludeExpr)
	}

	exclude, err = selectors.CompileMany(excludePatterns)
	if err != nil {
		return nil, nil, NewError(ErrCodeUsage, "invalid exclude selector", err)
	}

	return include, exclude, nil
}

func (app *Application) readSelectFile() ([]string, error) {
	data, err := os.ReadFile(app.config.Selection.SelectFile)
	if err != nil {
		return nil, err
	}

	var patterns []string
	for _, line := range splitLines(string(data)) {
		line = trimSpace(line)
		if line == "" || line[0] == '#' {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, nil
}

func (app *Application) renderOutput(model render.Model) (string, error) {
	opts := render.Options{
		Style:          app.config.Output.Style,
		ASCIIOnly:      app.config.Output.ASCIIOnly,
		NoHeader:       app.config.Output.NoHeader,
		HeaderCase:     app.config.Output.HeaderCase,
		MaxColWidth:    app.config.Output.MaxColWidth,
		WrapMode:       app.config.Output.WrapMode,
		TruncateSuffix: app.config.Output.TruncateSuffix,
		NullStr:        app.config.Output.NullStr,
		BoolStr:        app.config.Output.BoolStr,
		Precision:      app.config.Output.Precision,
		Color:          app.config.Output.Color,
	}

	return render.Render(model, opts)
}

func (app *Application) applyRowFiltering(rows []flatten.FlatKV) ([]flatten.FlatKV, error) {
	if len(app.config.Filter.WhereExprs) == 0 {
		return rows, nil
	}

	conditions, err := filter.ParseConditions(app.config.Filter.WhereExprs)
	if err != nil {
		return nil, NewError(ErrCodeUsage, "invalid filter condition", err)
	}

	rowFilter := filter.NewFilter(conditions)
	return rowFilter.Apply(rows), nil
}

func (app *Application) applySorting(rows []flatten.FlatKV) []flatten.FlatKV {
	if len(app.config.Sort.Columns) == 0 {
		return rows
	}

	// Parse comma-separated column specifications
	var expandedColumns []string
	for _, col := range app.config.Sort.Columns {
		expandedColumns = append(expandedColumns, splitCommaString(col)...)
	}

	sortOpts := sort.Options{
		Columns: expandedColumns,
	}

	sorter := sort.New(sortOpts)
	return sorter.Sort(rows)
}

func (app *Application) writeOutput(output string) error {
	var writer io.Writer = os.Stdout

	if app.config.Output.FilePath != "" {
		file, err := os.Create(app.config.Output.FilePath)
		if err != nil {
			return err
		}
		defer func() { _ = file.Close() }()
		writer = file
	}

	if !endsWithNewline(output) {
		output += "\n"
	}

	_, err := io.WriteString(writer, output)
	return err
}

// Helper functions
func splitCommaString(s string) []string {
	var result []string
	for _, part := range splitString(s, ",") {
		part = trimSpace(part)
		if part != "" {
			result = append(result, part)
		}
	}
	return result
}

func splitLines(s string) []string {
	return splitString(s, "\n")
}

func splitString(s, sep string) []string {
	if s == "" {
		return nil
	}
	var result []string
	start := 0
	for i := 0; i < len(s); i++ {
		if i+len(sep) <= len(s) && s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}
	result = append(result, s[start:])
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && isSpace(s[start]) {
		start++
	}
	for end > start && isSpace(s[end-1]) {
		end--
	}

	return s[start:end]
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func joinStrings(strings []string, sep string) string {
	if len(strings) == 0 {
		return ""
	}
	if len(strings) == 1 {
		return strings[0]
	}

	totalLen := len(sep) * (len(strings) - 1)
	for _, s := range strings {
		totalLen += len(s)
	}

	result := make([]byte, 0, totalLen)
	result = append(result, strings[0]...)
	for _, s := range strings[1:] {
		result = append(result, sep...)
		result = append(result, s...)
	}

	return string(result)
}

func endsWithNewline(s string) bool {
	return len(s) > 0 && s[len(s)-1] == '\n'
}
