# tablo – A CLI to Render JSON/YAML as Tables

## Overview

**tablo** is a cross-platform CLI tool that renders structured data (JSON or YAML) into human-readable tables. It supports single objects, arrays of objects, nested objects, and arrays of objects via explicit flattening ("dive") rules. The tool is modular, well-tested, and easy to extend.

### Primary Use Cases

- Inspecting API responses quickly.
- CLI pipelines that transform JSON/YAML to table form.
- Selective projection and flattening to compare nested data across array items.

## Naming

- Canonical CLI name: `tablo`

## Scope and Goals

### Must-haves

- **Input Sources:**
  - String argument (`--input`)
  - File path (`--file`)
  - Stdin (default)
- **Input Formats:** JSON or YAML (auto-detected unless specified)
- **Rendering:**
  - Single object → 2-column key/value table.
  - Array of objects → tabular layout with header row.
  - Array of primitives → single column (`VALUE`), optionally with index.
- **Flattening ("dive")**
  - Nested objects and arrays of objects can be flattened into dotted keys (e.g., `A.B.C`, `Arr.0.Item`).
  - Simple arrays and arrays of simple arrays are not flattened unless explicitly requested.
  - Max depth control for flattening.
- **Selection:**
  - Include only specified paths in output (supports dotted paths and wildcards).
  - Exclude paths after inclusion.
- **Table Styles:** Multiple styles, including Unicode box-drawing, ASCII, markdown, compact, and borderless.
- **Output:**
  - Write to stdout (default) or file (`--output`).
  - Header case transformation, column width, wrapping/truncation, null/bool formatting, float precision.
- **Error Handling:**
  - Clear exit codes and messages for parse errors, invalid flags, IO errors, selection errors.
- **Modularity:**
  - Idiomatic Go, small cohesive packages, clear interfaces, unit and integration tests.
- **Documentation:**
  - README, usage, examples, man page, contributing guide.

### Non-goals (v1)

- Editing/transformation of data beyond selection, flattening, and rendering.
- Query languages beyond simple wildcard selection (e.g., JQ expressions).
- Streaming large JSONL as rows (may be considered later).

## CLI Specification

### Command Synopsis

```
tablo [options] [--file PATH | --input STRING | -] [--]
```

### Flags and Options

#### Input

- `--file`, `-f`: Path to input file.
- `--input`, `-i`: Raw input string (JSON or YAML).
- `--format`, `-F`: Input format (`auto`, `json`, `yaml`, `yml`). Default: `auto`.
- Precedence: `--input` > `--file` > stdin.

#### Flattening ("dive")

- `--dive`, `-d`: Enable flattening of nested objects and arrays of objects into dotted paths.
- `--dive-path`, `-D`: Dive only into listed top-level keys/paths (dot notation, repeatable).
- `--max-depth`, `-m`: Maximum depth to dive. Default: unlimited.
- `--flatten-simple-arrays`: Convert arrays of primitives to comma-separated strings in cells.

#### Selection

- `--select`, `-s`: Comma-separated list of dotted path expressions to include. Supports wildcards (`*`, `_`, `?`).
- `--select-file`: Path to file containing one path expression per line (merged with `--select`).
- `--exclude`, `-E`: Comma-separated list of paths to exclude (applied after inclusion).
- `--strict-select`: Error when any selected path does not exist in the data. Default: false.

#### Output Formatting

- `--style`: Table style (`heavy`, `light`, `double`, `ascii`, `markdown`, `compact`, `borderless`). Default: `heavy`.
- `--ascii`: Force ASCII borders.
- `--no-header`: Omit header row.
- `--header-case`: Header case transformation (`original`, `upper`, `lower`, `title`).
- `--max-col-width`: Wrap/truncate cells exceeding width. Default: 0 (no limit).
- `--wrap`: Cell wrapping mode (`off`, `word`, `char`).
- `--truncate-suffix`: Suffix for truncated cells. Default: `…` (ellipsis).
- `--null-str`: String to represent null values. Default: `null`.
- `--bool-str`: Booleans as literals or custom mapping (`true:false`). Default: `true:false`.
- `--precision`: Decimal precision for floats. Default: -1 (no change).
- `--output`, `-o`: Write output to file instead of stdout.
- `--index-column`: For arrays, include an `INDEX` column.
- `--limit`: Limit number of rows printed.
- `--color`: Colorize output (`auto`, `always`, `never`). Default: `auto`.

#### General

- `--version`, `-v`: Print version and exit.
- `--help`, `-h`: Help.
- `--quiet`: Suppress non-error logging.

### Exit Codes

- 0: Success
- 2: Invalid CLI usage (flags, conflicts, missing arguments, parse error)
- 3: Unable to read input (file errors, stdin issues)
- 4: Input format detection or parse failure
- 5: Selection/flattening errors (only with `--strict-select` or internal unexpected errors)

## Input Format Detection

- With `--format`: Explicitly parse as specified; on failure, exit with code 4.
- With `--format=auto` (default):
  - If `--file` is provided, use file extension; fallback to content sniffing if unknown.
  - If `--input` or stdin, sniff content: `{` or `[` → try JSON first; else try YAML.
- YAML multi-document: Treated as array of documents by default.

## Rendering Rules

### Fundamental Types and Stringification

- Numbers: Render using provided precision if set; otherwise standard formatting.
- Booleans: Render using `--bool-str` mapping or true/false.
- Strings: Render as-is.
- Null: Render as `--null-str`.
- Non-scalar composite values (when not dived): Render as compact JSON string.

### Single Object

- Without `--dive`: Two-column table (`KEY`, `VALUE`). Nested objects/arrays as JSON string.
- With `--dive`: Flatten nested objects/arrays into dotted paths.

### Array of Objects

- Without `--dive`: Header row = union of top-level keys. Missing keys → empty cells.
- With `--dive`: Header = union of flattened keys. Arrays of objects included via dotted indices. Simple arrays remain stringified.
- Column order: Top-level keys first in natural order, then children sorted lexicographically; `--select`/`--include` overrides order.

### Array of Primitives

- Single column (`VALUE`), optionally with index.

### Selection Semantics

- `--select`/`--include`: Comma-separated dotted path expressions with wildcards. Only those columns/keys shown, in specified order.
- `--exclude`: Remove columns after inclusion.
- `--strict-select`: Error if any selected path resolves to zero matches.

### Alignment

- Header: Left aligned.
- Values: Numeric right aligned; others left aligned.

## Dive Rules in Detail

- Dive applies to nested objects and arrays of objects (flattened as `A.B.C`, `Arr.0.Item`).
- Does not apply to simple arrays or arrays of simple arrays unless explicitly requested.
- Max depth: Counts segments across object fields and array indices; stops flattening beyond max depth.

## Table Styles

- Supported: `heavy`, `light`, `double`, `ascii`, `markdown`, `compact`, `borderless`.
- `--ascii` overrides style to ASCII.
- Column width controls via `--max-col-width`, wrapping/truncation.

## Error Handling

- Conflicting inputs: exit code 2 with clear message.
- Invalid flag values: exit code 2.
- Read/parse errors: exit code 3/4 with diagnostic details.
- Selection mismatch with `--strict-select`: exit code 5 listing missing expressions.
- Dive logic respects rules; unsupported types are not dived, no error.

## Edge Cases & Rules

- Empty input: Print helpful message and exit with code 2 (unless `--quiet`).
- Mixed-type arrays: If all elements are not objects, render as single column; if mixed, promote primitives to row with only `VALUE` column.
- Deeply nested/large arrays: `--limit` applies; `--quiet`/`--no-header` reduce output.
- YAML anchors/aliases: Handled by YAML parser.

## Non-functional Requirements

- Performance: Handle inputs up to ~50 MB efficiently.
- Stability: Deterministic output for same input/flags.
- Portability: Linux primary target; CI builds for macOS/Windows.
- Code quality: Idiomatic Go, clear interfaces, high test coverage.

## Architecture & Modules

- Go version: ≥ 1.21
- Suggested dependencies:
  - CLI: github.com/spf13/cobra
  - Table rendering: github.com/jedib0t/go-pretty/table
  - YAML: gopkg.in/yaml.v3
  - JSON: encoding/json
- Package layout:
  - cmd/tablo: Cobra root command, flag parsing, wiring.
  - internal/input: Input reading, format detection.
  - internal/parse: JSON/YAML parsing.
  - internal/flatten: Flattening utilities, path traversal.
  - internal/selectors: Path expression matching, include/exclude logic.
  - internal/render: Table rendering, styles, alignment, formatting.
  - internal/model: Table model, alignment/type tags.
  - internal/util: Helpers, error types, stringification.
  - tests/: Golden test data, integration harness.

### Key Interfaces

- `type Reader interface { Read() ([]byte, error) }`
- `type Parser interface { Parse([]byte) (any, error) }`
- `type Flattener interface { Flatten(any, opts) (FlatMap or Rows, error) }`
- `type Selector interface { Filter(paths []PathExpr, data) (data, error) }`
- `type Renderer interface { Render(model.Table, Style, Options) (string, error) }`

## Testing Strategy

- **Unit tests:**
  - Format detection, parsing, flattening, selection, rendering, error cases.
- **Integration tests:**
  - CLI end-to-end with fixtures (stdin, file, string input; all flag combinations).
  - Golden outputs for deterministic verification.
- **Coverage targets:**
  - Unit: ≥80% across internal packages.
  - Integration: Cover all critical paths and flag combinations.

## Documentation & Deliverables

- README: Installation, usage, flags, examples, troubleshooting.
- Man page: Generated from Cobra help text.
- CONTRIBUTING: Development workflow, tests, code style, CI instructions.
- Examples directory: Curated inputs and expected outputs.
- CHANGELOG and versioning guidelines (SemVer recommended).

## Acceptance Criteria

- CLI binary `tablo` builds and installs.
- Example inputs produce correct outputs (automated integration tests).
- Flags work as specified.
- Tests pass locally and in CI.
- Documentation is present and `--help` is useful.
- Error cases have informative messages and non-zero exit codes.

## Future Enhancements (Out of Scope for v1)

- Row filtering (e.g., `--where path=value`).
- Sorting rows by column(s).
- JSON Lines (NDJSON) streaming mode.
- CSV/TSV export.
- Colorization and themes.
- Custom key ordering configs.
- JMESPath/JQ-like query language support.
- Interactive mode for paging large tables.
- Plugin architecture for custom renderers.
