#!/bin/bash

# Output Options Demo - Tablo CLI
# Comprehensive demonstration of output formats, data presentation, and file operations

set -eou pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

echo -e "\033[1;32m📤 TABLO CLI - OUTPUT OPTIONS DEMO\033[0m"

# =============================================================================
# OUTPUT DESTINATIONS
# =============================================================================

section "🎯 Output Destinations"

echo "Tablo can output to terminal, files, or other processes."
echo ""

echo -e "\033[1;35m1. Terminal Output (Default)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,department' --limit 2"

echo -e "\033[1;35m2. File Output (--output, -o)\033[0m"
TEMP_OUTPUT="/tmp/tablo_demo_output.txt"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,department,salary' --output \"$TEMP_OUTPUT\" --limit 3"
echo "Output saved to file:"
cat "$TEMP_OUTPUT"
rm -f "$TEMP_OUTPUT"

echo -e "\033[1;35m4. Pipe Output\033[0m"
echo "Output can be piped to other commands:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name' --limit 3 --style borderless --no-header | wc -l"

# =============================================================================
# CSV OUTPUT
# =============================================================================

section "📊 CSV Output"

echo -e "\033[1;35m1. Basic CSV Output\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style csv --select 'name,age,department' --limit 3"

echo -e "\033[1;35m2. CSV with Custom Data Formatting\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style csv --select 'name,active,salary,score' --bool-str '1:0' --precision 2 --limit 3"

echo -e "\033[1;35m3. CSV without Headers\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style csv --no-header --select 'name,age' --limit 3"

# =============================================================================
# HTML OUTPUT
# =============================================================================

section "🌐 HTML Output"

echo -e "\033[1;35m1. Basic HTML Table\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style html --select 'name,department,active' --limit 2"

echo -e "\033[1;35m2. HTML with Custom Formatting\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --style html --select 'name,price,in_stock,rating' --bool-str 'In Stock:Out of Stock' --precision 2 --limit 2"

# =============================================================================
# MARKDOWN OUTPUT
# =============================================================================

section "📝 Markdown Output"

echo -e "\033[1;35m1. Basic Markdown Table\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style markdown --select 'name,department,score' --limit 3"

echo -e "\033[1;35m2. Markdown with Custom Formatting\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --style markdown --select 'name,category,price,in_stock' --bool-str 'Available:Out of Stock' --precision 2 --limit 3"

# =============================================================================
# OUTPUT OPTIONS REFERENCE
# =============================================================================

section "📋 Output Options Reference"

echo -e "\033[1;35mOutput Destinations:\033[0m"
echo "  --output FILE, -o FILE    Write to file instead of terminal"
echo "  (no option)               Write to stdout (default)"
echo ""

# =============================================================================
# END OF DEMO
# =============================================================================

divider
echo -e "\033[1;32m✅ OUTPUT OPTIONS DEMO COMPLETE!\033[0m"
