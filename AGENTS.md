# AGENTS.md

This file provides context for AI coding agents working on the `tablo` project.

## Project Overview

**tablo** is a CLI tool written in Go that renders CSV/JSON/JSONL/YAML data as pretty tables. It's designed for command-line data visualization with powerful features for data transformation and filtering.

### Key Capabilities

- **Multiple input formats**: JSON, JSONL, CSV, YAML, and direct stdin
- **Nested data handling**: Flatten nested objects and arrays with configurable depth
- **Column operations**: Select/exclude columns using glob patterns, rename columns
- **Row operations**: Filter with expressions, sort by multiple columns, limit output
- **Multiple output styles**: Heavy, light, ASCII, Markdown, HTML, CSV, and more
- **Formatting controls**: Custom boolean strings, float precision, null values
- **Array handling**: Flatten simple arrays to comma-separated strings or expand to rows

## Project Structure

```
tablo/
├── cmd/tablo/          # Main application entry point
├── internal/
│   ├── app/           # Application logic and CLI command handlers
│   ├── filter/        # Row filtering (--where expressions)
│   ├── flatten/       # Nested object/array flattening (--dive)
│   ├── input/         # Input reading from files/stdin/strings
│   ├── parse/         # Format parsers (JSON, CSV, YAML, JSONL)
│   ├── render/        # Table rendering and output styles
│   ├── selectors/     # Column selection/exclusion logic
│   └── sort/          # Row sorting by columns
├── demo/data/         # Sample data files for testing
├── Makefile           # Build, test, release automation
├── go.mod             # Go module dependencies
└── README.md          # User documentation
```

## Technology Stack

- **Language**: Go 1.26.0
- **CLI Framework**: cobra (command-line interface)
- **Table Rendering**: go-pretty/v6 (table formatting)
- **Data Parsing**: 
  - `gopkg.in/yaml.v3` (YAML)
  - `tidwall/jsonc` (JSON with comments)
  - stdlib `encoding/csv` and `encoding/json`
- **Text Processing**: 
  - `golang.org/x/text` (text manipulation)
  - `golang.org/x/term` (terminal capabilities)

## Development Guidelines

### Building & Testing

```bash
# Build the binary
make build

# Run tests
make test

# Run full CI suite (vet, test, etc.)
make ci

# Run the built binary
./tablo -f demo/data/list.json
```

### Code Organization

- **`cmd/tablo/`**: Contains only the main entry point
- **`internal/`**: All internal packages (not importable by external projects)
  - Each package focuses on a single responsibility
  - Packages communicate through well-defined interfaces
  - Keep parsing, filtering, and rendering concerns separate

### Versioning & Releases

- Version embedded at build time via `-ldflags`
- Development builds: `dev-<git-hash>[-dirty]`
- Release builds: `vMAJOR.MINOR.PATCH` (from git tags)

To create a release:
```bash
make TAG=v0.5.0 tag    # Create and push tag
make release           # Build multi-platform binaries
```

## Common Tasks for AI Agents

### Adding a New Output Style

1. Update `internal/render/` to add the style constant
2. Implement the style in the table rendering logic
3. Add the style to CLI flag validation in `internal/app/`
4. Update README.md examples

### Adding a New Filter Operator

1. Modify `internal/filter/` to parse the new operator
2. Implement the comparison logic
3. Add test cases for the operator
4. Document in README.md under "Row filtering"

### Adding a New Input Format

1. Create parser in `internal/parse/`
2. Register format detection in `internal/input/`
3. Add CLI flag for explicit format specification
4. Add test data in `demo/data/`
5. Document in README.md

### Enhancing Flattening Logic

1. Modify `internal/flatten/` for new flattening rules
2. Consider impact on `--dive-path` and `--max-depth` flags
3. Ensure column selectors still work with flattened paths
4. Add test cases for edge cases

### Improving Performance

- Profile with `go test -bench=. -benchmem`
- Focus on: parsing, flattening, filtering, sorting bottlenecks
- Consider lazy evaluation for large datasets
- Memory allocations matter for large table rendering

## Testing Philosophy

- Unit tests for individual packages (parsing, filtering, sorting)
- Integration tests via `cmd/tablo/` with sample data
- Test edge cases: empty input, single values, deeply nested objects
- Use table-driven tests for multiple scenarios
- Test files live next to implementation files (`*_test.go`)

## Code Style

- Follow standard Go conventions (gofmt, go vet)
- Use meaningful variable names (prefer clarity over brevity)
- Document exported functions and types
- Keep functions focused and composable
- Avoid global state; pass dependencies explicitly

## Dependencies

When adding new dependencies:
- Prefer stdlib when possible
- Choose well-maintained, popular libraries
- Consider binary size impact (this is a CLI tool)
- Update `go.mod` with `go get`
- Run `make ci` to ensure compatibility

## User-Facing Changes

When adding features:
1. Update README.md with examples
2. Add sample data to `demo/data/` if needed
3. Consider backwards compatibility
4. Update CLI help text in cobra commands
5. Consider impact on existing flags and behaviors

## Debugging Tips

- Use `tablo -i '{...}' --dive` to see flattening behavior
- Add `--index-column` to track row numbers
- Compare `--style csv` output for data verification
- Check `demo/data/` files for test scenarios
- Use `go run cmd/tablo/main.go` for quick testing

## Key Abstractions

- **Rows**: Slice of maps (`[]map[string]interface{}`)
- **Flattening**: Convert nested structures to flat key-value pairs with dotted paths
- **Selectors**: Glob patterns matching column paths
- **Filters**: Conditional expressions evaluated against row values
- **Renderers**: Transform rows into formatted table output

## Important Notes

- Column order may vary when not explicitly selected
- `--dive` changes the data structure significantly
- Filtering and sorting happen after flattening
- CSV output respects column selection but uses a different renderer
- Index columns are synthetic (added after data processing)

## Getting Help

- See README.md for comprehensive user documentation
- Check Makefile targets for common development tasks
- Refer to individual package documentation in source files
- Use `go doc` for package/function documentation
