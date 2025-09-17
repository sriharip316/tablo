#!/bin/bash

# Basic Usage Demo - Tablo CLI
# Shows fundamental input/output operations and core concepts

set -eou pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

echo -e "\033[1;32m📚 TABLO CLI - BASIC USAGE DEMO\033[0m"

# =============================================================================
# WHAT IS TABLO?
# =============================================================================

section "🤔 What is Tablo?"

echo "Tablo converts JSON and YAML data into formatted tables for easy reading."
echo ""
echo "Basic concept:"
echo "  JSON/YAML → Tablo → Beautiful Table"
echo ""
echo "Perfect for:"
echo "  • API responses"
echo "  • Configuration files"
echo "  • Data analysis"
echo "  • Reports and exports"
echo ""

# =============================================================================
# SIMPLEST EXAMPLES
# =============================================================================

section "🎯 Simplest Examples"

echo -e "\033[1;35m1. Single Object → Table Row\033[0m"
demo_cmd "$TABLO_BIN --input '{\"name\":\"Alice\",\"age\":30,\"city\":\"NYC\"}'"

echo -e "\033[1;35m2. Array of Objects → Table Rows\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"name\":\"Alice\",\"age\":30},{\"name\":\"Bob\",\"age\":25}]'"

echo -e "\033[1;35m3. From File (Most Common)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\""

# =============================================================================
# INPUT METHODS
# =============================================================================

section "📥 Input Methods"

echo -e "\033[1;35m1. Direct String Input (--input)\033[0m"
echo "Good for: Quick tests, simple data"
demo_cmd "$TABLO_BIN --input '{\"product\":\"laptop\",\"price\":999.99}'"

echo -e "\033[1;35m2. File Input (--file)\033[0m"
echo "Good for: Existing files, larger datasets"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --limit 2"

echo -e "\033[1;35m3. Pipe Input (stdin)\033[0m"
echo "Good for: Command pipelines, API responses"
demo_cmd "echo '{\"status\":\"success\",\"items\":42}' | $TABLO_BIN"

echo -e "\033[1;35m4. JSON Lines Input\033[0m"
echo "Good for: Streaming data, large datasets, log processing"
demo_cmd "echo -e '{\"event\":\"login\",\"user\":\"alice\"}\n{\"event\":\"logout\",\"user\":\"bob\"}' | $TABLO_BIN"

# =============================================================================
# OUTPUT STYLES
# =============================================================================

section "🎨 Output Styles"

echo "Tablo supports many table styles. Here are the most popular:"
echo ""

echo -e "\033[1;35m1. Heavy (Default) - Bold and Professional\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"style\":\"heavy\",\"use\":\"default\"}]' --style heavy"

echo -e "\033[1;35m2. Light - Clean and Minimal\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"style\":\"light\",\"use\":\"clean\"}]' --style light"

echo -e "\033[1;35m5. Compact Style - Space Efficient\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"style\":\"compact\",\"use\":\"space\"}]' --style compact"

echo -e "\033[1;35m3. Double Style - Elegant & Distinctive\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"style\":\"double\",\"use\":\"elegant\"}]' --style double"

echo -e "\033[1;35m4. ASCII - Universal Compatibility\033[0m (--ascii for quick access)"
demo_cmd "$TABLO_BIN --input '[{\"style\":\"ascii\",\"use\":\"compatibility\"}]' --style ascii"

echo -e "\033[1;35m5. Borderless - Minimal and Modern\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"style\":\"borderless\",\"use\":\"modern\"}]' --style borderless"

echo -e "\033[1;35m3. Markdown - Perfect for Documentation\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"style\":\"markdown\",\"use\":\"docs\"}]' --style markdown"

echo -e "\033[1;35m2. HTML Style - Web Integration\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"style\":\"html\",\"use\":\"web\"}]' --style html"

echo -e "\033[1;35m3. CSV Style - Data Export\033[0m"
demo_cmd "$TABLO_BIN --input '[{\"style\":\"csv\",\"use\":\"export\"}]' --style csv"

# =============================================================================
# BASIC DATA CONTROL
# =============================================================================

section "🎛️ Basic Data Control"

echo -e "\033[1;35m1. Select Specific Columns\033[0m"
echo "Only show the columns you care about:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,age,department'"

echo -e "\033[1;35m2. Limit Number of Rows\033[0m"
echo "Don't overwhelm with too much data:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --limit 2"

echo -e "\033[1;35m3. Add Index Column\033[0m"
echo "Helpful for counting and referencing:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --index-column"

# =============================================================================
# FORMAT SPECIFICATION
# =============================================================================

section "🏷️ Format Specification"

echo -e "\033[1;35m1. Auto-Detection (Default)\033[0m"
echo "Tablo automatically detects format from file extension and content:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\""
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/user.yaml\""

echo -e "\033[1;35m2. Explicit JSON Format (--format json)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --format json"

echo -e "\033[1;35m3. Explicit YAML Format (--format yaml)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.yaml\" --format yaml"

echo -e "\033[1;35m4. JSON Lines (JSONL) Format\033[0m"
echo "One JSON value per line - perfect for streaming data:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.jsonl\" --select 'name,department'"

echo -e "\033[1;35m5. JSONL with Arrays (Auto-flattening)\033[0m"
echo "Arrays in JSONL are automatically flattened into individual rows:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users-array.jsonl\" --select 'name,active'"

echo -e "\033[1;35m6. Explicit JSONL Format (--format jsonl)\033[0m"
demo_cmd "echo -e '{\"name\":\"Alice\",\"age\":30}\n{\"name\":\"Bob\",\"age\":25}' | $TABLO_BIN --format jsonl"

echo -e "\033[1;35m7. Force Format on Ambiguous Input\033[0m"
echo "Useful when piping data without clear format indicators:"
demo_cmd "echo -e 'name: Alice\nage: 30\ndepartment: Engineering' | $TABLO_BIN --format yaml"

echo -e "\033[1;35m8. Format Override\033[0m"
echo "Force YAML parsing on a .json file (if it contains YAML):"
# This would normally fail, but demonstrates the concept
echo "# $TABLO_BIN --file some-yaml-data.json --format yaml"

# =============================================================================
# SPECIAL INPUT CASES
# =============================================================================

section "⚡ Special Input Cases"

echo -e "\033[1;35m1. JSON with Comments (JSONC)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/comments.jsonc\""

echo -e "\033[1;35m2. YAML with Comments\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/comments.yaml\""

echo -e "\033[1;35m3. Empty Input Handling\033[0m"
demo_cmd "echo '{}' | $TABLO_BIN"
demo_cmd "echo '[]' | $TABLO_BIN"

echo -e "\033[1;35m4. Single Value Input\033[0m"
demo_cmd "echo '\"simple string\"' | $TABLO_BIN"
demo_cmd "echo '42' | $TABLO_BIN"
demo_cmd "echo 'true' | $TABLO_BIN"

# =============================================================================
# END OF DEMO
# =============================================================================

divider
echo -e "\033[1;32m✅ BASIC USAGE DEMO COMPLETE!\033[0m"
