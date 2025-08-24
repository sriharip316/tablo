# tablo

A CLI tool to render JSON/YAML as pretty tables. It can flatten nested objects, select columns, format booleans/floats, and output in multiple styles (heavy/light/double/ascii/markdown).

## Install / Build

- Requires Go 1.22+
- Easiest (uses version injection & reproducible flags):
  - `make build` (outputs `bin/tablo`)
  - `make install` (installs to `GOBIN`/`GOPATH/bin`)
- Direct with Go (no Makefile conveniences):
  - `go build -ldflags "-s -w" ./cmd/tablo`
- Run without installing:
  - `go run ./cmd/tablo --help`
- Show version:
  - `tablo --version`

### Makefile Highlights

The provided `Makefile` supports:

| Task | Command | Notes |
|------|---------|-------|
| Build local binary | `make build` | Injects version via `-ldflags -X main.version=...` |
| Install | `make install` | Same flags as build |
| Lint | `make lint` | Uses `golangci-lint` if available, else `go vet` |
| Test | `make test` | Race detector enabled |
| Coverage | `make cover` / `make cover-html` | Enforces minimum (configurable via `MIN_COVER`) |
| Multi-platform release | `make release` | Binaries in `dist/` + `sha256sums.txt` |
| Tag a release | `make VERSION=1.2.3 tag` | Creates & pushes `v1.2.3` |
| CI pipeline | `make ci` | tidy + lint + test + cover |

### Version Resolution

`tablo --version` reports (in priority order):
1. Injected build-time value (from `-ldflags -X main.version=...`)
2. Latest Git tag (`vX.Y.Z`); if workspace dirty: `vX.Y.Z-dirty`
3. Otherwise a development build: `dev-<short-hash>` (dirty adds `-dirty`)
4. Final fallback: `dev`

Examples:
- Tagged clean commit: `v0.3.1`
- Tagged with uncommitted changes: `v0.3.1-dirty`
- Untagged clean commit (hash `a1b2c3d`): `dev-a1b2c3d`
- Untagged dirty: `dev-a1b2c3d-dirty`

To force a version (e.g. during packaging):
```
make VERSION=1.4.0 build
./bin/tablo --version   # -> 1.4.0
```

## Quick start

```bash
# From a JSON string
tablo -i '{"a":1,"b":2}'

# From a file
tablo -f data.json

# From piped input
echo '{"a":1,"b":2}' | tablo -i -
```

## Comment Support

Both JSON and YAML input support comments:

```bash
# JSON with comments (JSONC format)
tablo -i '{
  "name": "test", // Line comment
  /* Block comment */
  "value": 42
}'

# YAML with comments (native support)
tablo -i '
# Top-level comment
name: test  # End of line comment
value: 42
'

# JSONC files with .jsonc extension
tablo -f data.jsonc
```

Comments are supported across all input methods (inline, file, stdin). Files with `.jsonc` extension are automatically detected as JSON with comments.

## Examples

### Flatten a JSON object to key/value pairs

Command:

```bash
tablo -i '{"a":{"b":1},"tags":["x","y",3]}' --dive --flatten-simple-arrays
```

Output:

```
┏━━━━━━┳━━━━━━━━━┓
┃ KEY  ┃ VALUE   ┃
┣━━━━━━╋━━━━━━━━━┫
┃ a.b  ┃ 1       ┃
┃ tags ┃ x, y, 3 ┃
┗━━━━━━┻━━━━━━━━━┛
```

Notes:

- `--dive` flattens nested objects (e.g., `a.b`).
- `--flatten-simple-arrays` converts arrays of primitives into a comma-separated string.

### YAML array of objects with selected columns and index

Command:

```bash
tablo -F yaml --index-column --select 'name,age' --style ascii <<'YAML'
- name: Alice
  age: 30
- name: Bob
  age: 31
YAML
```

Output:

```
+---+-------+-----+
|   | name  | age |
+---+-------+-----+
| 1 | Alice |  30 |
| 2 | Bob   |  31 |
+---+-------+-----+
```

Notes:

- `--select` accepts a comma-separated list of dotted paths. Use `--select-file` to load from a file (one per line).
- `--index-column` adds an auto index column for row arrays.
- Use `--limit N` to restrict the number of rows printed.

### Array of primitives

Command:

```bash
tablo -i '[1,2,3,4]' --limit 3 --style markdown
```

Output:

```markdown
| VALUE |
| ----- |
| 1     |
| 2     |
| 3     |
```

### Formatting options (booleans, precision, null)

When rendering rows, you can customize formatting:

- `--bool-str 'Y:N'` to render booleans as custom strings.
- `--precision 2` to format floats with 2 decimal places.
- `--null-str null` to show missing values as the literal `null`.

Example:

```bash
tablo -i '[{"a":1.2345,"b":true},{"b":false}]' --style ascii --precision 2 --bool-str 'Y:N' --index-column
```

Output:

```
+---+------+---+
|   | a    | b |
+---+------+---+
| 1 | 1.23 | Y |
| 2 | null | N |
+---+------+---+
```

## Flattening controls

- `--dive` enables flattening of nested objects and arrays of objects.
- `--dive-path k1 --dive-path k2` only dives into the listed top-level keys.
- `--max-depth N` limits flattening depth (`-1` = unlimited).

## Output styles

Choose a table style with `--style`:

- `heavy` (default), `light`, `double`, `ascii`, `markdown`, `compact`, `borderless`.
- Force ASCII borders with `--ascii`.

## Selecting/excluding columns

Use dotted path expressions with glob support per segment (`*` and `?`). Examples:

- Include: `--select 'user.*.name,meta.id'`
- Exclude: `--exclude 'debug.*'`
- Strict mode: `--strict-select` fails if any selected path is missing.

## Versioning & Releases

- Stable releases are tagged with semantic versions: `vMAJOR.MINOR.PATCH`.
- Binaries built from an exact tag report that tag (e.g., `v0.4.0`).
- Non-tag builds report a development identifier: `dev-<short-hash>`.
- A `-dirty` suffix is appended when there are uncommitted changes.

To create a new release:
```
# ensure clean working tree and tests pass
make ci
# choose next version
make VERSION=0.5.0 tag
# build multi-platform artifacts
make VERSION=0.5.0 release
```

You can verify the version embedded in a built binary:
```
./bin/tablo --version
```

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
