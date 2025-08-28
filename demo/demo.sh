#!/bin/bash

# Demo script for tablo CLI
# Showcases ALL CLI options and features available in tablo
# Works from any directory.

set -eou pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DEMO_DIR="$PROJECT_ROOT/demo"
TABLO_BIN="$PROJECT_ROOT/bin/tablo"

# Build tablo CLI if not present
if ! [ -x "$TABLO_BIN" ]; then
    echo "tablo CLI not found, building..."
    (cd "$PROJECT_ROOT" && make build)
fi

divider() {
    printf "\n"
    printf '%*s\n' "${COLUMNS:-80}" '' | tr ' ' '='
    printf "\n"
}

section() {
    divider
    echo -e "\033[1;36m$1\033[0m"
    divider
}

demo_cmd() {
    echo -e "\033[1;33m\$ $*\033[0m"
    eval "$@"
    echo
}

echo -e "\033[1;32m🎯 TABLO CLI - COMPLETE FEATURE DEMONSTRATION\033[0m"
echo "This demo showcases ALL available CLI options and features"

# =============================================================================
# INPUT OPTIONS
# =============================================================================

section "📥 INPUT OPTIONS"

echo -e "\033[1;35m1. File Input (--file, -f)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\""

echo -e "\033[1;35m2. Raw Input String (--input, -i)\033[0m"
demo_cmd "$TABLO_BIN --input '{\"name\":\"Alice\",\"age\":30,\"city\":\"NYC\"}'"

echo -e "\033[1;35m3. Format Specification (--format, -F)\033[0m"
echo "Testing explicit format specification:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --format json"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.yaml\" --format yaml"

echo -e "\033[1;35m4. Piped Input (stdin)\033[0m"
demo_cmd "echo '{\"product\":\"laptop\",\"price\":999.99,\"available\":true}' | $TABLO_BIN"

echo -e "\033[1;35m5. JSON with Comments Support\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/comments.jsonc\""

echo -e "\033[1;35m6. YAML with Comments Support\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/comments.yaml\""

# =============================================================================
# FLATTENING/DIVING OPTIONS
# =============================================================================

section "🏗️ FLATTENING & DIVING OPTIONS"

echo -e "\033[1;35m1. Basic Diving (--dive, -d)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --dive"

echo -e "\033[1;35m2. Dive Specific Paths (--dive-path, -D)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --dive --dive-path meta --dive-path settings"

echo -e "\033[1;35m3. Maximum Depth (--max-depth, -m)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --dive --max-depth 2"

echo -e "\033[1;35m4. Flatten Simple Arrays (--flatten-simple-arrays)\033[0m"
demo_cmd "$TABLO_BIN --input '{\"tags\":[\"red\",\"blue\",\"green\"],\"numbers\":[1,2,3,4]}' --dive --flatten-simple-arrays"

# =============================================================================
# SELECTION OPTIONS
# =============================================================================

section "🎯 SELECTION & FILTERING OPTIONS"

echo -e "\033[1;35m1. Select Columns (--select, -s)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --select 'id,name,active'"

echo -e "\033[1;35m2. Select with Glob Patterns\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --dive --select 'meta.owner.*'"

echo -e "\033[1;35m3. Select from File (--select-file)\033[0m"
echo "id" > /tmp/tablo_select.txt
echo "name" >> /tmp/tablo_select.txt
echo "meta.rating" >> /tmp/tablo_select.txt
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --dive --select-file /tmp/tablo_select.txt"
rm -f /tmp/tablo_select.txt

echo -e "\033[1;35m4. Exclude Columns (--exclude, -E)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --dive --exclude 'description,items'"

echo -e "\033[1;35m5. Exclude with Glob Patterns\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --dive --exclude 'meta.*'"

echo -e "\033[1;35m6. Strict Selection (--strict-select)\033[0m"
echo "This will error if selected paths don't exist:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --select 'id,nonexistent' --strict-select" || echo "Expected error: path not found"

# =============================================================================
# OUTPUT FORMATTING OPTIONS
# =============================================================================

section "🎨 OUTPUT FORMATTING OPTIONS"

echo -e "\033[1;35m1. Table Styles (--style)\033[0m"
for style in heavy light double ascii markdown compact borderless; do
    echo -e "\033[1;34mStyle: $style\033[0m"
    "$TABLO_BIN" --file "$DEMO_DIR/list.json" --style "$style" --limit 1
    echo
done

echo -e "\033[1;35m1a. HTML Output (--style html)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --style html --limit 2"

echo -e "\033[1;35m1b. CSV Output (--style csv)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --style csv --limit 2"

echo -e "\033[1;35m2. ASCII Only Mode (--ascii)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --style heavy --ascii --limit 1"

echo -e "\033[1;35m3. No Header Row (--no-header)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --no-header --limit 2"

echo -e "\033[1;35m4. Header Case Options (--header-case)\033[0m"
for case in original upper lower title; do
    echo -e "\033[1;34mHeader case: $case\033[0m"
    "$TABLO_BIN" --file "$DEMO_DIR/list.json" --header-case "$case" --select 'id,name' --limit 1
    echo
done

echo -e "\033[1;35m5. Maximum Column Width (--max-col-width)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --max-col-width 20 --select 'description'"

echo -e "\033[1;35m6. Cell Wrapping Options (--wrap)\033[0m"
echo -e "\033[1;34mWrap: off\033[0m"
"$TABLO_BIN" --file "$DEMO_DIR/single.json" --max-col-width 30 --wrap off --select 'description'
echo

echo -e "\033[1;34mWrap: word\033[0m"
"$TABLO_BIN" --file "$DEMO_DIR/single.json" --max-col-width 30 --wrap word --select 'description'
echo

echo -e "\033[1;34mWrap: char\033[0m"
"$TABLO_BIN" --file "$DEMO_DIR/single.json" --max-col-width 30 --wrap char --select 'description'
echo

echo -e "\033[1;35m7. Truncate Suffix (--truncate-suffix)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --max-col-width 20 --truncate-suffix '...' --select 'description'"

echo -e "\033[1;35m8. Null String Representation (--null-str)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --null-str 'N/A' --dive --select 'misc,meta.owner.contact.phone'"

echo -e "\033[1;35m9. Boolean String Mapping (--bool-str)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --bool-str 'YES:NO' --select 'active'"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --bool-str '✓:✗' --select 'active'"

echo -e "\033[1;35m10. Float Precision (--precision)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --precision 1 --dive --select 'meta.rating'"
demo_cmd "$TABLO_BIN --input '[{\"pi\":3.14159},{\"e\":2.71828}]' --precision 2"

echo -e "\033[1;35m11. Output to File (--output, -o)\033[0m"
OUTFILE="$DEMO_DIR/demo_output.txt"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --style markdown --output \"$OUTFILE\" --limit 2"
echo "Output written to $OUTFILE:"
cat "$OUTFILE"
rm -f "$OUTFILE"

echo -e "\033[1;35m12. Index Column (--index-column)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --index-column --select 'name' --limit 3"

echo -e "\033[1;35m13. Limit Rows (--limit)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --limit 2"

echo -e "\033[1;35m14. Color Output (--color)\033[0m"
echo -e "\033[1;34mColor: auto (default)\033[0m"
"$TABLO_BIN" --file "$DEMO_DIR/list.json" --color auto --limit 1
echo

echo -e "\033[1;34mColor: always\033[0m"
"$TABLO_BIN" --file "$DEMO_DIR/list.json" --color always --limit 1
echo

echo -e "\033[1;34mColor: never\033[0m"
"$TABLO_BIN" --file "$DEMO_DIR/list.json" --color never --limit 1
echo

# =============================================================================
# GENERAL OPTIONS
# =============================================================================

section "⚙️ GENERAL OPTIONS"

echo -e "\033[1;35m1. Quiet Mode (--quiet)\033[0m"
echo "Normal output:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/nonexistent.json\" 2>&1" || echo "Error shown normally"

echo "Quiet mode (errors suppressed):"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/nonexistent.json\" --quiet 2>&1" || echo "Error suppressed in quiet mode"

echo -e "\033[1;35m2. Version Information (--version)\033[0m"
demo_cmd "$TABLO_BIN --version"

echo -e "\033[1;35m3. Help Information (--help)\033[0m"
echo "Use '$TABLO_BIN --help' to see all available options"

# =============================================================================
# COMPLEX COMBINATIONS
# =============================================================================

section "🚀 COMPLEX FEATURE COMBINATIONS"

echo -e "\033[1;35m1. Deep Diving with Selection and Formatting\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --dive --max-depth 3 --select 'meta.owner.*,settings.notifications.*' --style ascii --index-column"

echo -e "\033[1;35m2. List Processing with Multiple Options\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --dive --select 'id,name,meta.rating,meta.owner.username' --style markdown --header-case title --bool-str 'Y:N' --precision 1 --limit 2"

echo -e "\033[1;35m3. Advanced Formatting with Custom Strings\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --dive --exclude 'description,items' --style double --null-str 'EMPTY' --bool-str '✅:❌' --max-col-width 15 --wrap word --truncate-suffix '...' --header-case upper"

echo -e "\033[1;35m4. Complex Selection with Exclusion\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --dive --flatten-simple-arrays --select 'id,name,meta.*' --exclude 'meta.tags,meta.created' --style borderless --index-column"

echo -e "\033[1;35m5. Pipeline Processing Example\033[0m"
demo_cmd "echo '[{\"name\":\"John\",\"scores\":[95,87,92]},{\"name\":\"Jane\",\"scores\":[88,94,90]}]' | $TABLO_BIN --dive --flatten-simple-arrays --style ascii --header-case title --index-column"

echo -e "\033[1;35m6. HTML Output with Formatting\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"name\":\"Alice\",\"score\":95.123,\"active\":true},{\"name\":\"Bob\",\"score\":87.456,\"active\":false}]' --style html --precision 1 --bool-str 'Yes:No' --index-column"

echo -e "\033[1;35m7. CSV Export with Precision and Boolean Formatting\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"product\":\"Laptop\",\"price\":999.99,\"in_stock\":true},{\"product\":\"Mouse\",\"price\":19.95,\"in_stock\":false}]' --style csv --precision 2 --bool-str 'Available:Out of Stock'"

echo -e "\033[1;35m8. CSV Output with Flattened Nested Data\033[0m"
demo_cmd "$TABLO_BIN --input '{\"company\":{\"name\":\"TechCorp\",\"employees\":[{\"name\":\"Alice\",\"role\":\"Engineer\"},{\"name\":\"Bob\",\"role\":\"Manager\"}]}}' --dive --style csv"

echo -e "\033[1;35m9. HTML Output with Custom Null Strings\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --dive --style html --null-str 'N/A' --select 'id,name,misc' --limit 2"

echo -e "\033[1;35m10. CSV Output without Headers\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"city\":\"New York\",\"population\":8400000},{\"city\":\"Los Angeles\",\"population\":3900000}]' --style csv --no-header"

# =============================================================================
# EDGE CASES AND ERROR HANDLING
# =============================================================================

section "🧪 EDGE CASES & ERROR HANDLING"

echo -e "\033[1;35m1. Empty Input\033[0m"
demo_cmd "echo '{}' | $TABLO_BIN"
demo_cmd "echo '[]' | $TABLO_BIN"

echo -e "\033[1;35m2. Mixed Data Types\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"a\":1,\"b\":\"text\"},{\"a\":true,\"b\":null},{\"a\":3.14,\"b\":[1,2,3]}]' --dive --flatten-simple-arrays"

echo -e "\033[1;35m3. Unicode and Special Characters\033[0m"
demo_cmd "$TABLO_BIN --input '{\"emoji\":\"🌟⭐✨\",\"unicode\":\"ಅ भा ௶ ¥ тест 测\",\"special\":\"!@#$%^&*()\"}'"

echo -e "\033[1;35m4. Large Numbers and Precision\033[0m"
demo_cmd "$TABLO_BIN --input '{\"big\":123456789.123456789,\"small\":0.000123456}' --precision 3"

divider
echo -e "\033[1;32m🎉 ALL TABLO CLI FEATURES DEMONSTRATED!\033[0m"
echo -e "\033[1;36mSummary of demonstrated features:\033[0m"
echo "📥 Input: --file, --input, --format, stdin, JSON/YAML with comments"
echo "🏗️ Flattening: --dive, --dive-path, --max-depth, --flatten-simple-arrays"
echo "🎯 Selection: --select, --select-file, --exclude, --strict-select"
echo "🎨 Formatting: --style (9 styles including html/csv), --ascii, --no-header, --header-case"
echo "📏 Layout: --max-col-width, --wrap, --truncate-suffix"
echo "🔧 Data: --null-str, --bool-str, --precision"
echo "📤 Output: --output, --index-column, --limit, --color"
echo "⚙️ General: --quiet, --version, --help"
echo "✨ New: HTML and CSV output formats for web and spreadsheet integration"
echo ""
echo "For more information, run: $TABLO_BIN --help"
divider
