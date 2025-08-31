#!/bin/bash

# Formatting & Styling Demo - Tablo CLI
# Comprehensive demonstration of table styles, colors, and appearance options

set -eou pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

echo -e "\033[1;32m🎨 TABLO CLI - FORMATTING & STYLING DEMO\033[0m"

# =============================================================================
# HEADER FORMATTING
# =============================================================================

section "📝 Header Formatting"

echo -e "\033[1;35m1. Original Case (Default)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --header-case original --select 'name,age,department' --limit 2"

echo -e "\033[1;35m2. Uppercase Headers\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --header-case upper --select 'name,age,department' --limit 2"

echo -e "\033[1;35m3. Lowercase Headers\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --header-case lower --select 'name,age,department' --limit 2"

echo -e "\033[1;35m4. Title Case Headers\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --header-case title --select 'name,age,department' --limit 2"

echo -e "\033[1;35m5. No Headers\033[0m"
echo "Useful for data-only output:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --no-header --select 'name,age' --limit 2"

# =============================================================================
# COLUMN WIDTH & WRAPPING
# =============================================================================

section "📏 Column Width & Text Wrapping"

echo -e "\033[1;35m1. Default Width (No Limit)\033[0m"
LONG_DATA='[{"name": "Alice", "description": "Senior Software Engineer with expertise in distributed systems, microservices, and cloud architecture"}]'
demo_cmd "$TABLO_BIN --input '$LONG_DATA'"

echo -e "\033[1;35m2. Maximum Column Width\033[0m"
demo_cmd "$TABLO_BIN --input '$LONG_DATA' --max-col-width 30"

echo -e "\033[1;35m3. Text Wrapping: Off (Default)\033[0m"
demo_cmd "$TABLO_BIN --input '$LONG_DATA' --max-col-width 30 --wrap off"

echo -e "\033[1;35m4. Text Wrapping: Word\033[0m"
echo "Wrap at word boundaries:"
demo_cmd "$TABLO_BIN --input '$LONG_DATA' --max-col-width 30 --wrap word"

echo -e "\033[1;35m5. Text Wrapping: Character\033[0m"
echo "Wrap at any character:"
demo_cmd "$TABLO_BIN --input '$LONG_DATA' --max-col-width 30 --wrap char"

echo -e "\033[1;35m6. Custom Truncation Suffix\033[0m"
demo_cmd "$TABLO_BIN --input '$LONG_DATA' --max-col-width 25 --truncate-suffix '... (more)'"

# =============================================================================
# DATA REPRESENTATION
# =============================================================================

section "🔢 Data Representation"

echo -e "\033[1;35m1. Boolean Formatting - Default\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,active' --limit 3"

echo -e "\033[1;35m2. Custom Boolean Strings\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,active' --bool-str 'YES:NO' --limit 3"

echo -e "\033[1;35m3. Boolean with Emojis\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,active' --bool-str '✅:❌' --limit 3"

echo -e "\033[1;35m4. Null Value Representation\033[0m"
NULL_DATA='[{"name": "Alice", "phone": null}, {"name": "Bob", "phone": "555-1234"}]'
echo "Default null representation:"
demo_cmd "$TABLO_BIN --input '$NULL_DATA'"

echo "Custom null string:"
demo_cmd "$TABLO_BIN --input '$NULL_DATA' --null-str 'N/A'"

echo -e "\033[1;35m5. Number Precision\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,score' --precision 0 --limit 3"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,score' --precision 2 --limit 3"

# =============================================================================
# STYLE COMBINATIONS
# =============================================================================

section "🎯 Style Combinations"

echo -e "\033[1;35m1. Professional Report Style\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style heavy --header-case title --bool-str 'Active:Inactive' --precision 1 --index-column --limit 3"

echo -e "\033[1;35m2. Modern Minimal Style\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style borderless --header-case lower --bool-str '✓:✗' --null-str '—' --limit 3"

echo -e "\033[1;35m3. Technical Documentation Style\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style markdown --header-case upper --bool-str 'TRUE:FALSE' --limit 3"

echo -e "\033[1;35m5. Executive Summary Style\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style double --header-case title --bool-str 'Yes:No' --precision 0 --select 'name,department,salary' --limit 3"

# =============================================================================
# EXPORT-ORIENTED FORMATTING
# =============================================================================

section "📤 Export-Oriented Formatting"

echo -e "\033[1;35m1. CSV Export with Custom Formatting\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style csv --bool-str '1:0' --precision 2 --limit 3"

echo -e "\033[1;35m2. HTML with Custom Styling\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style html --bool-str 'Active:Inactive' --null-str 'Not Available' --limit 2"

echo -e "\033[1;35m3. Markdown for Documentation\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --style markdown --header-case title --bool-str 'In Stock:Out of Stock' --precision 2 --select 'name,price,in_stock' --limit 3"

# =============================================================================
# WORKING WITH UNICODE DATA
# =============================================================================

section "🌍 Unicode & International Data"

echo -e "\033[1;35m1. Unicode Data with Different Styles\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,emoji,location' --style heavy --limit 2"

echo -e "\033[1;35m2. Unicode with ASCII Mode\033[0m"
echo "ASCII mode for compatibility with unicode content:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,emoji,location' --style heavy --ascii --limit 2"

echo -e "\033[1;35m3. Unicode in Different Export Formats\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/unicode.json\" --select 'name,symbols' --style csv --limit 2"

# =============================================================================
# REAL-WORLD STYLING EXAMPLES
# =============================================================================

section "💼 Real-World Styling Examples"

echo -e "\033[1;35m1. Financial Report\033[0m"
FINANCIAL_DATA='[
  {"quarter": "Q1 2023", "revenue": 1250000.50, "profit": 187500.75, "growth": true},
  {"quarter": "Q2 2023", "revenue": 1375000.25, "profit": 206250.04, "growth": true},
  {"quarter": "Q3 2023", "revenue": 1180000.00, "profit": 177000.00, "growth": false}
]'
demo_cmd "$TABLO_BIN --input '$FINANCIAL_DATA' --style double --header-case title --precision 2 --bool-str 'Positive:Negative' --index-column"

echo -e "\033[1;35m2. System Status Dashboard\033[0m"
STATUS_DATA='[
  {"service": "web-server", "status": "running", "uptime": 99.95, "healthy": true},
  {"service": "database", "status": "degraded", "uptime": 97.80, "healthy": false},
  {"service": "cache", "status": "running", "uptime": 99.99, "healthy": true}
]'
demo_cmd "$TABLO_BIN --input '$STATUS_DATA' --style borderless --bool-str '🟢:🔴' --precision 2 --header-case upper"

echo -e "\033[1;35m3. Product Inventory\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --style light --header-case title --bool-str 'Available:Sold Out' --precision 2 --select 'name,category,price,in_stock,rating' --limit 4"

echo -e "\033[1;35m4. Employee Directory\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --style heavy --header-case title --bool-str 'Active:Inactive' --null-str 'Not Set' --index-column --limit 4"

# =============================================================================
# ADVANCED FORMATTING TECHNIQUES
# =============================================================================

section "🚀 Advanced Formatting Techniques"

echo -e "\033[1;35m1. Conditional Formatting with Selection\033[0m"
echo "Combine formatting with smart column selection:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --dive --select 'name,age,department,active' --style markdown --bool-str '✅ Active:❌ Inactive' --header-case title --limit 3"

echo -e "\033[1;35m2. Multi-Format Output Workflow\033[0m"
echo "Generate different formats for different audiences:"
echo ""
echo "For technical teams (detailed):"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --style ascii --precision 3 --header-case lower --limit 2"

echo "For executives (summary):"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --style double --header-case title --select 'name,category,price' --precision 0 --limit 2"

echo -e "\033[1;35m3. Responsive Formatting\033[0m"
echo "Adjust based on available space:"
echo ""
echo "Full width (detailed view):"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --style borderless --limit 2"

echo "Narrow width (mobile/compact view):"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --style compact --max-col-width 15 --wrap word --select 'name,price,in_stock' --limit 2"

# =============================================================================
# STYLE QUICK REFERENCE
# =============================================================================

section "📋 Style Quick Reference"

echo -e "\033[1;35mTable Styles:\033[0m"
echo "  --style heavy        Bold, professional (default)"
echo "  --style light        Clean, minimal"
echo "  --style double       Elegant double lines"
echo "  --style ascii        ASCII-only compatibility"
echo "  --style markdown     Documentation format"
echo "  --style compact      Space-efficient"
echo "  --style borderless   Modern, minimal"
echo "  --style html         Web integration"
echo "  --style csv          Data export"
echo ""

echo -e "\033[1;35mAppearance Options:\033[0m"
echo "  --header-case original|upper|lower|title"
echo "  --no-header          Omit header row"
echo "  --color auto|always|never"
echo "  --ascii              Force ASCII borders"
echo "  --index-column       Add row numbers"
echo ""

echo -e "\033[1;35mText Formatting:\033[0m"
echo "  --max-col-width N    Limit column width"
echo "  --wrap off|word|char Text wrapping mode"
echo "  --truncate-suffix STR Truncation indicator"
echo ""

echo -e "\033[1;35mData Formatting:\033[0m"
echo "  --bool-str 'true:false'  Boolean representation"
echo "  --null-str 'null'        Null value string"
echo "  --precision N            Decimal places"

# =============================================================================
# END OF DEMO
# =============================================================================

divider
echo -e "\033[1;32m✅ FORMATTING & STYLING DEMO COMPLETE!\033[0m"
