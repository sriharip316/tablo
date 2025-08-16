# tablo

A CLI tool to render JSON/YAML as pretty tables. It can flatten nested objects, select columns, format booleans/floats, and output in multiple styles (heavy/light/double/ascii/markdown).

## Install / Build

- Requires Go 1.22+
- Build locally:
  - `go build ./cmd/tablo`
- Or run without installing:
  - `go run ./cmd/tablo --help`

## Quick start

```bash
# From a JSON string
tablo -i '{"a":1,"b":2}'

# From a file
tablo -f data.json

# From piped input
echo '{"a":1,"b":2}' | tablo -i -
```

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
