package app

// Default values for application configuration
const (
	// Input defaults
	DefaultFormat = "auto"

	// Flattening defaults
	DefaultMaxDepth = -1 // unlimited

	// Output defaults
	DefaultStyle          = "heavy"
	DefaultHeaderCase     = "original"
	DefaultMaxColWidth    = 0 // no limit
	DefaultWrapMode       = "off"
	DefaultTruncateSuffix = "…"
	DefaultNullStr        = "null"
	DefaultBoolStr        = "true:false"
	DefaultPrecision      = -1 // no change
	DefaultLimit          = 0  // no limit
	DefaultColor          = "auto"

	// Selection defaults
	DefaultStrictSelect = false

	// General defaults
	DefaultQuiet = false
)

// Style constants
const (
	StyleHeavy      = "heavy"
	StyleLight      = "light"
	StyleDouble     = "double"
	StyleASCII      = "ascii"
	StyleMarkdown   = "markdown"
	StyleCompact    = "compact"
	StyleBorderless = "borderless"
	StyleHTML       = "html"
	StyleCSV        = "csv"
)

// Format constants
const (
	FormatAuto  = "auto"
	FormatJSON  = "json"
	FormatYAML  = "yaml"
	FormatYML   = "yml"
	FormatCSV   = "csv"
	FormatJSONL = "jsonl"
)

// Header case constants
const (
	HeaderCaseOriginal = "original"
	HeaderCaseUpper    = "upper"
	HeaderCaseLower    = "lower"
	HeaderCaseTitle    = "title"
)

// Wrap mode constants
const (
	WrapModeOff  = "off"
	WrapModeWord = "word"
	WrapModeChar = "char"
)

// Color constants
const (
	ColorAuto   = "auto"
	ColorAlways = "always"
	ColorNever  = "never"
)

// File extensions
const (
	ExtJSON  = ".json"
	ExtJSONC = ".jsonc"
	ExtYAML  = ".yaml"
	ExtYML   = ".yml"
)

// Special column names
const (
	ColumnNameValue = "VALUE"
	ColumnNameKey   = "KEY"
	ColumnNameIndex = "INDEX"
)

// Comment prefixes
const (
	CommentPrefix         = "#"
	JSONLineComment       = "//"
	JSONBlockCommentStart = "/*"
	JSONBlockCommentEnd   = "*/"
)

// Separators and delimiters
const (
	PathSeparator    = "."
	CommaSeparator   = ","
	ColonSeparator   = ":"
	NewlineSeparator = "\n"
)

// Validation limits
const (
	MaxInputSizeMB    = 50
	MaxInputSizeBytes = MaxInputSizeMB * 1024 * 1024
	MaxColumnWidth    = 1000
	MaxPrecision      = 20
	MaxDepthLimit     = 100
)

// Messages
const (
	MsgNoInputAvailable     = "no input available"
	MsgInvalidFormat        = "invalid format"
	MsgConflictingInputs    = "conflicting inputs: --input and --file cannot be used together"
	MsgMissingSelectedPaths = "missing selected paths"
)
