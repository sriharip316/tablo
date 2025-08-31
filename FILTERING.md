# Row Filtering Examples

This document demonstrates the row filtering capabilities of `tablo` using the `--where` flag.

## Sample Data

The examples below use this employee dataset:

```json
[
  {"name": "John", "age": 30, "department": "Engineering", "salary": 75000, "active": true, "email": "john@company.com"},
  {"name": "Jane", "age": 25, "department": "Marketing", "salary": 65000, "active": true, "email": "jane@company.com"},
  {"name": "Bob", "age": 35, "department": "Engineering", "salary": 85000, "active": false, "email": "bob@company.com"},
  {"name": "Alice", "age": 28, "department": "Sales", "salary": 70000, "active": true, "email": "alice@company.com"},
  {"name": "Charlie", "age": 32, "department": "Engineering", "salary": 80000, "active": true, "email": "charlie@company.com"},
  {"name": "Diana", "age": 29, "department": "Marketing", "salary": 68000, "active": true, "email": "diana@company.com"},
  {"name": "Eve", "age": 27, "department": "Sales", "salary": 72000, "active": false, "email": "eve@company.com"},
  {"name": "Frank", "age": 33, "department": "Engineering", "salary": 90000, "active": true, "email": "frank@company.com"}
]
```

## Filtering Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `=` | Equal to | `name=John` |
| `!=` | Not equal to | `department!=Sales` |
| `>` | Greater than | `age>30` |
| `>=` | Greater than or equal | `salary>=75000` |
| `<` | Less than | `age<30` |
| `<=` | Less than or equal | `salary<=70000` |
| `~` | Contains | `name~an` |
| `!~` | Does not contain | `email!~gmail` |
| `=~` | Regex match | `email=~.*@company\.com` |
| `!=~` | Regex does not match | `name!=~^[A-C].*` |

## Basic Examples

### 1. Equality Filter

Filter employees in the Engineering department:

```bash
tablo -f demo/employees.json --where "department=Engineering"
```

### 2. Numeric Comparison

Find employees with salary greater than $75,000:

```bash
tablo -f demo/employees.json --where "salary>75000"
```

### 3. Boolean Filter

Show only active employees:

```bash
tablo -f demo/employees.json --where "active=true"
```

### 4. String Contains

Find employees whose names contain "an":

```bash
tablo -f demo/employees.json --where "name~an"
```

### 5. Not Equal

Show employees not in Sales:

```bash
tablo -f demo/employees.json --where "department!=Sales"
```

## Advanced Examples

### 6. Multiple Conditions (AND Logic)

Find active Engineering employees aged 30 or older:

```bash
tablo -f demo/employees.json --where "department=Engineering" --where "active=true" --where "age>=30"
```

### 7. Regex Pattern Matching

Find employees with company emails:

```bash
tablo -f demo/employees.json --where "email=~.*@company\.com"
```

### 8. Range Filtering

Find employees aged between 25 and 30:

```bash
tablo -f demo/employees.json --where "age>=25" --where "age<=30"
```

## Combined with Other Features

### 9. Filtering + Column Selection

Show only names and salaries for Engineering employees:

```bash
tablo -f demo/employees.json --where "department=Engineering" --select "name,salary"
```

### 10. Filtering + Output Format

Export filtered data as CSV:

```bash
tablo -f demo/employees.json --where "active=true" --where "salary>=70000" --style csv
```

### 11. Filtering + Limiting

Show top 3 highest-paid employees:

```bash
tablo -f demo/employees.json --where "salary>=75000" --limit 3
```

### 12. Filtering + Index Column

Add row numbers to filtered results:

```bash
tablo -f demo/employees.json --where "department=Marketing" --index-column
```

## Working with Nested Data

When using `--dive` for flattened nested objects, filtering works on the flattened paths:

```bash
# Sample nested data
echo '[
  {"user": {"name": "John", "profile": {"age": 30}}, "status": {"active": true}},
  {"user": {"name": "Jane", "profile": {"age": 25}}, "status": {"active": false}}
]' | tablo --dive --where "user.profile.age>25"
```

## Error Handling

Invalid filter expressions will show helpful error messages:

```bash
# Invalid operator
tablo -f demo/employees.json --where "name John"
# Error: invalid filter condition: no valid operator found

# Invalid regex
tablo -f demo/employees.json --where "email=~[invalid"
# Error: invalid filter condition: invalid regex pattern
```

## Tips and Best Practices

1. **Quote your expressions**: Always quote filter expressions to avoid shell interpretation
2. **Escape regex patterns**: Use `\.` for literal dots in regex patterns
3. **Type-aware comparisons**: Numeric comparisons work with numbers, string comparisons with text
4. **Multiple conditions**: Use multiple `--where` flags for AND logic
5. **Test incrementally**: Start with simple filters and add complexity gradually

## Real-World Use Cases

### Finding Issues in Logs

```bash
# Filter error logs from JSON log files
tablo -f logs.json --where "level=ERROR" --where "timestamp>2024-01-01" --select "message,timestamp"
```

### Analyzing API Responses

```bash
# Find failed API calls
tablo -f api_responses.json --where "status!=200" --select "url,status,response_time"
```

### Processing Configuration Data

```bash
# Find enabled services with high memory usage
tablo -f services.json --where "enabled=true" --where "memory_mb>512" --select "name,memory_mb"
```

### Data Quality Checks

```bash
# Find records with missing email addresses
tablo -f users.json --where "email=" --select "id,name,email"
```
