# tablo

A CLI tool to render JSON/YAML as pretty tables. It supports flattening of nested objects, selecting/excluding columns, filtering rows, and multiple output styles.

## Quick start

```bash
# From a JSON string
tablo -i '{"a":1,"b":2}'

# From a file
tablo -f demo/data/list.json

# From standard input
echo '{"a":1,"b":2}' | tablo
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

- `--select` accepts a comma-separated list of dotted paths. Use `--select-file` to load column selections from a file (one per line).
- `--index-column` adds an auto index column for row arrays.
- Use `--limit N` to restrict the number of printed rows.

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

### Sorting rows

Command:

```bash
tablo -i '[{"name":"Charlie","age":35},{"name":"Alice","age":30},{"name":"Bob","age":25}]' --sort age
```

Output:

```
┏━━━━━┳━━━━━━━━━┓
┃ age ┃ name    ┃
┣━━━━━╋━━━━━━━━━┫
┃ 25  ┃ Bob     ┃
┃ 30  ┃ Alice   ┃
┃ 35  ┃ Charlie ┃
┗━━━━━┻━━━━━━━━━┛
```

#### Per-column sort direction

You can specify sort direction for each column individually using `+` (ascending) or `-` (descending) prefixes:

```bash
# Sort by department (ascending), then by age (descending)
tablo -f data.json --sort "department,-age"

# Explicit ascending prefix (same as no prefix)
tablo -f data.json --sort "+name,-salary"

# Mixed directions with multiple columns
tablo -f data.json --sort "active,-salary,name"
```

Notes:

- `--sort column1,column2` sorts by multiple columns in order
- `--sort +column1,-column2` sorts column1 ascending, column2 descending
- Works with flattened paths (e.g., `--sort user.name,-user.age`)

<old_text line=252>
### Row sorting

Sort rows using the `--sort` flag with column names:

- `--sort 'name'` - sort by a single column
- `--sort 'name,age'` - sort by multiple columns (comma-separated)

Sorting supports different data types:
- **Numbers**: sorted numerically (e.g., 1, 2, 10, 100)
- **Strings**: sorted alphabetically 
- **Booleans**: false comes before true
- **Mixed types**: fall back to string comparison
- **Null values**: always sorted first

This works with flattened paths when using `--dive`.

Example:

```bash
tablo -f employees.json --sort 'department,age' --select 'name,department,age'
```

Output:

```
┏━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━┓
┃ name    ┃ department  ┃ age ┃
┣━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━┫
┃ Bob     ┃ Engineering ┃ 25  ┃
┃ Charlie ┃ Engineering ┃ 35  ┃
┃ David   ┃ Marketing   ┃ 28  ┃
┃ Alice   ┃ Marketing   ┃ 30  ┃
┗━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━┛
```

### CSV and HTML output

Export data as CSV for use in spreadsheet applications:

```bash
tablo -i '[{"name":"John","age":30},{"name":"Jane","age":25}]' --style csv
```

Output:

```
age,name
30,John
25,Jane
```

Generate HTML tables for use in web applications:

```bash
echo '{"user":"admin","active":true}' | tablo --dive --style html
```

Output:

```html
<table class="go-pretty-table">
  <thead>
    <tr>
      <th>KEY</th>
      <th>VALUE</th>
    </tr>
  </thead>
  <tbody>
    <tr>
      <td>active</td>
      <td>true</td>
    </tr>
    <tr>
      <td>user</td>
      <td>admin</td>
    </tr>
  </tbody>
</table>
```

### Row filtering

Filter rows using the `--where` flag with condition expressions:

- `--where 'name=John'` - equality comparison
- `--where 'age>25'` - numeric comparison (`>`, `>=`, `<`, `<=`)
- `--where 'active=true'` - boolean comparison
- `--where 'name~pattern'` - string contains (`~` for contains, `!~` for not contains)
- `--where 'email=~.*@example\.com'` - regex matching (`=~` for match, `!=~` for not match)

Multiple `--where` flags are combined using AND logic. This works with flattened paths when using `--dive`.

Example:

```bash
tablo -f employees.json --where 'department=Engineering' --where 'salary>75000' --select 'name,salary'
```

Output:

```
┏━━━━━━━━━┳━━━━━━━━┓
┃ name    ┃ salary ┃
┣━━━━━━━━━╋━━━━━━━━┫
┃ Bob     ┃ 85000  ┃
┃ Charlie ┃ 80000  ┃
┃ Frank   ┃ 90000  ┃
┗━━━━━━━━━┻━━━━━━━━┛
```

### Row sorting with per-column direction

Sort rows using the `--sort` flag with column names and optional direction prefixes:

- `--sort 'name'` - sort by a single column
- `--sort 'name,age'` - sort by multiple columns (comma-separated)
- `--sort '+name,-age'` - sort by name ascending, then age descending

Direction prefixes:
- `+column` or `column` - ascending order (default)
- `-column` - descending order

Sorting supports different data types:
- **Numbers**: sorted numerically (e.g., 1, 2, 10, 100)
- **Strings**: sorted alphabetically 
- **Booleans**: false comes before true
- **Mixed types**: fall back to string comparison
- **Null values**: always sorted first

This works with flattened paths when using `--dive`.

Example:

```bash
tablo -f employees.json --sort 'department,-age' --select 'name,department,age'
```

Output:

```
┏━━━━━━━━━┳━━━━━━━━━━━━━┳━━━━━┓
┃ name    ┃ department  ┃ age ┃
┣━━━━━━━━━╋━━━━━━━━━━━━━╋━━━━━┫
┃ Charlie ┃ Engineering ┃ 35  ┃
┃ Bob     ┃ Engineering ┃ 25  ┃
┃ Alice   ┃ Marketing   ┃ 30  ┃
┃ David   ┃ Marketing   ┃ 28  ┃
┗━━━━━━━━━┻━━━━━━━━━━━━━┻━━━━━┛
```

### Formatting options (booleans, precision, null)

You can customize formatting when rendering rows:

- `--bool-str 'Y:N'` to render booleans as custom strings.
- `--precision 2` to format floats with 2 decimal places.
- `--null-str null` to display missing values as the literal `null`.

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
  - `--dive-path k1 --dive-path k2` dives only into the listed top-level keys.
- `--max-depth N` limits flattening depth (`-1` = unlimited).

## Output styles

Choose a table style with `--style`:

- `heavy` (default), `light`, `double`, `ascii`, `markdown`, `compact`, `borderless`, `html`, `csv`.
- Force ASCII borders with `--ascii` (applies only to table styles).

## Selecting/excluding columns

Use dotted path expressions with glob support for each segment (`*` and `?`). Examples:

- Include: `--select 'user.*.name,meta.id'`
- Exclude: `--exclude 'debug.*'`
- Strict mode: `--strict-select` fails if any selected path is missing.

## Versioning & Releases

- Stable releases are tagged with semantic versions: `vMAJOR.MINOR.PATCH`.
- Binaries built from an exact tag report that tag (e.g., `v0.4.0`).
- Non-tag builds report a development identifier: `dev-<short-hash>`.
- A `-dirty` suffix is appended if there are uncommitted changes.

To create a new release:

```
# ensure clean working tree and tests pass
make ci
# choose the next version and create a tag
make TAG=v0.5.0 tag
# build multi-platform artifacts (automatically detects version from tag)
make release
```

The release process:

- `make tag` validates the working tree is clean and creates/pushes the git tag
- `make release` runs `release-check` to validate git state and builds for multiple platforms
- Release artifacts are built for: linux/amd64, linux/arm64, darwin/amd64, darwin/arm64, and windows/amd64.
- All binaries are placed in `dist/` along with a `sha256sums.txt` file.

## Future Enhancements

- JSON Lines (NDJSON) streaming mode.
- Colorization and themes.
- Custom key ordering configs.
- JMESPath/JQ-like query language support.
- Interactive mode for paging large tables.
- Plugin architecture for custom renderers.
