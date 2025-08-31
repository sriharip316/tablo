#!/bin/bash

# Selection & Column Control Demo - Tablo CLI
# Comprehensive demonstration of column selection, exclusion, and control features

set -eou pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

echo -e "\033[1;32m🎯 TABLO CLI - SELECTION & COLUMN CONTROL DEMO\033[0m"

# =============================================================================
# UNDERSTANDING COLUMN SELECTION
# =============================================================================

section "🤔 Understanding Column Selection"

echo "Column selection lets you choose exactly which data to display."
echo ""
echo "Benefits:"
echo "• Focus on relevant information"
echo "• Reduce visual clutter"
echo "• Create targeted reports"
echo "• Improve readability"
echo ""
echo "First, let's see all available columns:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --limit 2"

# =============================================================================
# BASIC COLUMN SELECTION
# =============================================================================

section "🎯 Basic Column Selection"

echo -e "\033[1;35m1. Select Specific Columns (--select, -s)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,age,department' --limit 3"

echo -e "\033[1;35m3. Single Column Selection\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name' --limit 5"

echo -e "\033[1;35m4. Reorder Columns\033[0m"
echo "Columns appear in the order you specify:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'department,name,age' --limit 3"

# =============================================================================
# COLUMN EXCLUSION
# =============================================================================

section "🚫 Column Exclusion"

echo -e "\033[1;35m1. Exclude Specific Columns (--exclude, -E)\033[0m"
echo "Hide columns you don't need:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --exclude 'email,city' --limit 2"

echo -e "\033[1;35m3. Multiple Exclusions\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --exclude 'id,email,city,score' --limit 3"

# =============================================================================
# WORKING WITH NESTED DATA
# =============================================================================

section "🏗️ Selection with Nested Data"

echo "Selection becomes powerful when combined with flattening:"
echo ""

echo -e "\033[1;35m1. View Nested Structure First\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive --limit 3"

echo -e "\033[1;35m2. Select Nested Columns\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive --select 'id,name,meta.owner.username,meta.rating'"

echo -e "\033[1;35m3. Complex Nested Selection\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/list.json\" --dive --select 'id,name,settings.theme' --limit 2"

echo -e "\033[1;35m4. Exclude Nested Columns\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive --exclude 'description,items,settings.notifications' --limit 3"

# =============================================================================
# GLOB PATTERNS
# =============================================================================

section "🌟 Glob Pattern Matching"

echo "Use wildcards for powerful column selection:"
echo ""

echo -e "\033[1;35m1. Select All Meta Columns\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive --select 'id,name,settings.notifications.*'"

echo -e "\033[1;35m2. Wildcard at End\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/list.json\" --dive --select 'id,name,settings.notifications.*' --limit 2"

echo -e "\033[1;35m3. Wildcard in Middle\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/list.json\" --dive --select 'id,name,items.*.vaue' --limit 2"

echo -e "\033[1;35m5. Exclude with Wildcards\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive --exclude 'items.*.*' --limit 3"

# =============================================================================
# FILE-BASED SELECTION
# =============================================================================

section "📄 File-Based Selection"

echo "Store column lists in files for reusable selection patterns:"
echo ""

echo -e "\033[1;35m1. Create Selection File\033[0m"
SELECTION_FILE="/tmp/tablo_columns.txt"
cat > "$SELECTION_FILE" << 'EOF'
id
name
items.*.value
settings.theme
EOF

echo "Selection file contents:"
cat "$SELECTION_FILE"
echo ""

echo -e "\033[1;35m2. Use Selection File (--select-file)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/list.json\" --select-file \"$SELECTION_FILE\" --dive"

# Clean up temp files
rm -f "$SELECTION_FILE"

# =============================================================================
# STRICT SELECTION MODE
# =============================================================================

section "🔒 Strict Selection Mode"

echo "Control what happens when selected columns don't exist:"
echo ""

echo -e "\033[1;35m1. Normal Selection (Missing Columns Ignored)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,nonexistent_column,age' --limit 2"

echo -e "\033[1;35m2. Strict Selection (--strict-select)\033[0m"
echo "This will error if any selected column doesn't exist:"
echo "\$ tablo --file users.json --select 'name,nonexistent_column' --strict-select"
echo "Expected: Error - column 'nonexistent_column' not found"
echo ""

echo -e "\033[1;35m3. Strict Selection with Valid Columns\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,age,department' --strict-select --limit 2"

# =============================================================================
# REAL-WORLD SELECTION EXAMPLES
# =============================================================================

section "💼 Real-World Selection Examples"

echo -e "\033[1;35m1. Employee Directory\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,email,department' --where 'active=true' --style borderless --header-case title"

echo -e "\033[1;35m2. Performance Review Data\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,department,score,salary' --where 'score>85' --style heavy --precision 1 --index-column"

echo -e "\033[1;35m3. Contact Information Export\033[0m"
TEMP_CONTACTS="/tmp/tablo_contacts.csv"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,email,city' --style csv --header-case title --output \"$TEMP_CONTACTS\""
echo "Contact list exported:"
head -5 "$TEMP_CONTACTS"
rm -f "$TEMP_CONTACTS"

echo -e "\033[1;35m4. Product Catalog Summary\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --select 'name,category,price,in_stock' --where 'price<200' --style markdown --bool-str 'Available:Out of Stock'"

echo -e "\033[1;35m5. System Monitoring Dashboard\033[0m"
MONITORING_DATA='[
  {"service": "web", "cpu": 45.2, "memory": 67.8, "status": "healthy", "uptime": 99.9},
  {"service": "db", "cpu": 78.1, "memory": 82.3, "status": "warning", "uptime": 99.5},
  {"service": "cache", "cpu": 23.4, "memory": 34.1, "status": "healthy", "uptime": 100.0}
]'
demo_cmd "$TABLO_BIN --input '$MONITORING_DATA' --select 'service,status,cpu,memory' --where 'cpu>50' --style ascii --precision 1"

# =============================================================================
# SELECTION PATTERNS REFERENCE
# =============================================================================

section "📋 Selection Patterns Reference"

echo -e "\033[1;35mBasic Selection:\033[0m"
echo "  --select 'col1,col2,col3'        Select specific columns"
echo "  --select 'name'                  Select single column"
echo "  --exclude 'unwanted1,unwanted2'  Exclude specific columns"
echo ""

echo -e "\033[1;35mNested Data Selection:\033[0m"
echo "  --select 'user.profile.name'     Select nested column"
echo "  --select 'meta.*'                Select all meta columns"
echo "  --select '*.name'                Select name from any parent"
echo ""

echo -e "\033[1;35mAdvanced Patterns:\033[0m"
echo "  --select-file columns.txt        Load selection from file"
echo "  --strict-select                  Error on missing columns"
echo ""


# =============================================================================
# END OF DEMO
# =============================================================================

divider
echo -e "\033[1;32m✅ SELECTION & COLUMN CONTROL DEMO COMPLETE!\033[0m"
