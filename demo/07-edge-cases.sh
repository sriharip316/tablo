#!/bin/bash

# Edge Cases & Special Scenarios Demo - Tablo CLI
# Comprehensive demonstration of handling Unicode, mixed types, empty data, and edge cases

set -eou pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

echo -e "\033[1;32m🧪 TABLO CLI - EDGE CASES & SPECIAL SCENARIOS\033[0m"

# =============================================================================
# UNICODE AND INTERNATIONAL CHARACTERS
# =============================================================================

section "🌍 Unicode & International Characters"

echo "Tablo handles international text, emojis, and special characters gracefully."
echo ""

echo -e "\033[1;35m1. International Names and Text\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,location,description' --limit 3"

echo -e "\033[1;35m2. Emoji and Symbols\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,emoji,symbols' --limit 3"

echo -e "\033[1;35m3. Mathematical and Scientific Symbols\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,math,arrows' --limit 2"

echo -e "\033[1;35m4. Mixed Scripts and Languages\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --dive --select 'name,languages' --flatten-simple-arrays --limit 3"

echo -e "\033[1;35m5. Unicode with Different Styles\033[0m"
echo "Heavy style (Unicode borders):"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,emoji' --style heavy --limit 2"

echo "ASCII style (compatible borders):"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,emoji' --style ascii --limit 2"

echo -e "\033[1;35m6. Unicode Export Formats\033[0m"
echo "CSV export with Unicode:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,location,emoji' --style csv --limit 2"

echo "HTML export with Unicode:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,emoji,symbols' --style html --limit 2"

# =============================================================================
# EMPTY AND NULL DATA
# =============================================================================

echo -e "\033[1;35m5. Sparse Data Structures\033[0m"
SPARSE_DATA='[
  {"id": 1, "name": "Alice", "department": "Engineering"},
  {"id": 2, "name": "Bob", "location": "NYC", "salary": 75000},
  {"id": 3, "name": "Carol", "department": "Sales", "bonus": 5000}
]'
demo_cmd "$TABLO_BIN --input '$SPARSE_DATA'"

echo -e "\033[1;35m6. Empty Strings vs Null\033[0m"
EMPTY_VS_NULL='[
  {"name": "Alice", "middle": "", "suffix": null},
  {"name": "Bob", "middle": "J", "suffix": "Jr"},
  {"name": "Carol", "middle": null, "suffix": ""}
]'
demo_cmd "$TABLO_BIN --input '$EMPTY_VS_NULL' --null-str 'NULL' --style borderless"

# =============================================================================
# EDGE CASE VALUES
# =============================================================================

section "🔢 Edge Case Values"

echo -e "\033[1;35m3. Scientific Notation and Large Numbers\033[0m"
SCIENTIFIC_DATA='[
  {"measurement": "Distance to Moon", "meters": 384400000, "scientific": "3.844e8"},
  {"measurement": "Electron Mass", "kg": 0.00000000000000000000000000091093837015, "scientific": "9.109e-31"},
  {"measurement": "Avogadro Number", "per_mole": 602214076000000000000000, "scientific": "6.022e23"}
]'
demo_cmd "$TABLO_BIN --input '$SCIENTIFIC_DATA'"

echo -e "\033[1;35m4. Boolean Edge Cases\033[0m"
BOOL_EDGE_CASES='[
  {"test": "true_bool", "value": true, "as_string": "true"},
  {"test": "false_bool", "value": false, "as_string": "false"},
  {"test": "truthy_string", "value": "yes", "as_string": "yes"},
  {"test": "falsy_string", "value": "no", "as_string": "no"}
]'
demo_cmd "$TABLO_BIN --input '$BOOL_EDGE_CASES' --bool-str 'TRUE:FALSE'"
demo_cmd "$TABLO_BIN --input '$BOOL_EDGE_CASES' --bool-str '1:0'"

echo -e "\033[1;35m5. Arrays with Mixed Content\033[0m"
MIXED_ARRAYS='[
  {"name": "Alice", "data": [1, "text", true, null, 3.14]},
  {"name": "Bob", "data": [{"nested": "object"}, [1, 2, 3], "simple"]},
  {"name": "Carol", "data": []}
]'
demo_cmd "$TABLO_BIN --input '$MIXED_ARRAYS' --dive --select name,data"
demo_cmd "$TABLO_BIN --input '$MIXED_ARRAYS' --dive --select name,data --flatten-simple-arrays"

echo -e "\033[1;35m3. Very Large Column Names\033[0m"
LARGE_COLUMNS='{
  "this_is_an_extremely_long_column_name_that_might_cause_formatting_issues": "value1",
  "another_very_long_column_name_for_testing_purposes_and_edge_cases": "value2",
  "short": "value3"
}'
demo_cmd "$TABLO_BIN --input '$LARGE_COLUMNS'"

echo "With column width limits:"
demo_cmd "$TABLO_BIN --input '$LARGE_COLUMNS' --max-col-width 20 --wrap word"

echo -e "\033[1;35m4. Empty vs Missing Keys\033[0m"
EMPTY_VS_MISSING='[
  {"name": "Alice", "email": "alice@example.com", "phone": ""},
  {"name": "Bob", "email": "", "phone": "555-1234"},
  {"name": "Carol", "email": "carol@example.com"}
]'
demo_cmd "$TABLO_BIN --input '$EMPTY_VS_MISSING'"

# =============================================================================
# SPECIAL CHARACTERS AND FORMATTING
# =============================================================================

section "🔣 Special Characters and Formatting"

echo -e "\033[1;35m1. Characters That Affect Table Formatting\033[0m"
SPECIAL_CHARS='[
  {"type": "newlines", "text": "Line 1\nLine 2\nLine 3"},
  {"type": "tabs", "text": "Col1\tCol2\tCol3"},
  {"type": "unicode_spaces", "text": "word\u2002word\u2003word"},
  {"type": "control_chars", "text": "text\u0001with\u0002control\u0003chars"}
]'
demo_cmd "$TABLO_BIN --input '$SPECIAL_CHARS' --select type,text"

echo -e "\033[1;35m2. Very Long Text\033[0m"
LONG_TEXT='[
  {
    "description": "This is an extremely long description that will test how Tablo handles text that exceeds normal column widths and might need to be wrapped or truncated depending on the output settings and terminal width limitations."
  }
]'
echo "Default (no wrapping):"
demo_cmd "$TABLO_BIN --input '$LONG_TEXT'"

echo "With column width limit:"
demo_cmd "$TABLO_BIN --input '$LONG_TEXT' --max-col-width 30"

echo "With word wrapping:"
demo_cmd "$TABLO_BIN --input '$LONG_TEXT' --max-col-width 30 --wrap word"

echo -e "\033[1;35m3. JSON/CSV Special Characters\033[0m"
CSV_SPECIAL='[
  {"field": "commas", "value": "one, two, three"},
  {"field": "quotes", "value": "He said \"Hello\" to me"},
  {"field": "mixed", "value": "Comma, \"quotes\", and more"}
]'
echo "Table format:"
demo_cmd "$TABLO_BIN --input '$CSV_SPECIAL'"

echo "CSV format (note escaping):"
demo_cmd "$TABLO_BIN --input '$CSV_SPECIAL' --style csv"

echo -e "\033[1;35m4. HTML Special Characters\033[0m"
HTML_SPECIAL='[
  {"element": "paragraph", "tag": "<p>Hello World</p>"},
  {"element": "link", "tag": "<a href=\"example.com\">Link</a>"},
  {"element": "script", "tag": "<script>alert(\"test\")</script>"}
]'
echo "Table format:"
demo_cmd "$TABLO_BIN --input '$HTML_SPECIAL'"

echo "HTML format (note how special chars are handled):"
demo_cmd "$TABLO_BIN --input '$HTML_SPECIAL' --style html"

# =============================================================================
# DEEPLY NESTED AND COMPLEX STRUCTURES
# =============================================================================

section "🏗️ Complex and Deeply Nested Structures"

echo -e "\033[1;35m1. Very Deep Nesting\033[0m"
DEEP_NESTED='{
  "level1": {
    "level2": {
      "level3": {
        "level4": {
          "level5": {
            "level6": {
              "level7": {
                "deep_value": "Found me!"
              }
            }
          }
        }
      }
    }
  }
}'
echo "Without depth limit:"
demo_cmd "$TABLO_BIN --input '$DEEP_NESTED' --dive"

echo "With depth limit (3 levels):"
demo_cmd "$TABLO_BIN --input '$DEEP_NESTED' --dive --max-depth 3"

echo -e "\033[1;35m3. Arrays of Arrays\033[0m"
ARRAY_OF_ARRAYS='[
  {
    "matrix": [
      [1, 2, 3],
      [4, 5, 6],
      [7, 8, 9]
    ],
  }
]'
demo_cmd "$TABLO_BIN --input '$ARRAY_OF_ARRAYS' --dive"
demo_cmd "$TABLO_BIN --input '$ARRAY_OF_ARRAYS' --dive --flatten-simple-arrays"

echo -e "\033[1;35m4. Mixed Array and Object Nesting\033[0m"
MIXED_NESTING='[
  {
    "users": [
      {
        "name": "Alice",
        "preferences": {
          "theme": "dark",
          "languages": ["en", "es", "fr"]
        }
      },
      {
        "name": "Bob",
        "preferences": {
          "theme": "light",
          "languages": ["en", "de"]
        }
      }
    ]
  }
]'
demo_cmd "$TABLO_BIN --input '$MIXED_NESTING' --dive --max-depth 2 --flatten-simple-arrays"

# =============================================================================
# PERFORMANCE EDGE CASES
# =============================================================================

section "⚡ Performance Edge Cases"

echo -e "\033[1;35m1. Many Columns\033[0m"
echo "Simulating wide data with many columns:"
MANY_COLUMNS='[
  {
    "col_01": "A", "col_02": "B", "col_03": "C", "col_04": "D", "col_05": "E",
    "col_06": "F", "col_07": "G", "col_08": "H", "col_09": "I", "col_10": "J",
    "col_11": "K", "col_12": "L", "col_13": "M", "col_14": "N", "col_15": "O"
  }
]'
demo_cmd "$TABLO_BIN --input '$MANY_COLUMNS' --style compact"

echo -e "\033[1;35m2. Managing Wide Output\033[0m"
echo "Select subset of columns:"
demo_cmd "$TABLO_BIN --input '$MANY_COLUMNS' --select 'col_01,col_05,col_10,col_15'"

echo "Use compact style:"
demo_cmd "$TABLO_BIN --input '$MANY_COLUMNS' --style borderless --max-col-width 5"

# =============================================================================
# EDGE CASES WITH FILTERING
# =============================================================================

section "🔍 Edge Cases with Filtering"

NULL_DATA='[
  {"name": "Alice", "phone": null, "email": "alice@example.com", "age": 30},
  {"name": "Bob", "phone": "555-1234", "email": null, "age": null},
  {"name": "Carol", "phone": null, "email": null, "age": 25}
]'

echo -e "\033[1;35m1. Filtering Null Values\033[0m"
demo_cmd "$TABLO_BIN --input '$NULL_DATA' --where 'phone!=null' --select 'name,phone'"

echo -e "\033[1;35m2. Filtering Empty Strings\033[0m"
demo_cmd "$TABLO_BIN --input '$EMPTY_VS_NULL' --where 'middle!=' --select 'name,middle'"

echo -e "\033[1;35m3. Numeric Edge Cases in Filtering\033[0m"
NUMERIC_EDGE='[
  {"name": "Zero", "value": 0},
  {"name": "Negative", "value": -42},
  {"name": "Float", "value": 0.001},
  {"name": "Large", "value": 1000000}
]'
demo_cmd "$TABLO_BIN --input '$NUMERIC_EDGE' --where 'value>0' --precision 3"
demo_cmd "$TABLO_BIN --input '$NUMERIC_EDGE' --where 'value<=0'"

echo -e "\033[1;35m4. Boolean Edge Cases in Filtering\033[0m"
demo_cmd "$TABLO_BIN --input '$BOOL_EDGE_CASES' --where 'value=true'"

# =============================================================================
# END OF DEMO
# =============================================================================

divider
echo -e "\033[1;32m✅ EDGE CASES & SPECIAL SCENARIOS DEMO COMPLETE!\033[0m"
