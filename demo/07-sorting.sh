#!/bin/bash

# Sorting Demo - Tablo CLI
# Shows how to sort data based on one or more columns.

set -eou pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

echo -e "\033[1;32m📚 TABLO CLI - SORTING DEMO\033[0m"

# =============================================================================
# SORTING
# =============================================================================

section "🔀 Sorting"

echo "Tablo allows you to sort the output table based on one or more columns."
echo ""

echo -e "\033[1;35m1. Sort by a Single Column (age)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --sort age"

echo -e "\033[1;35m2. Sort in Reverse Order (age, descending)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --sort '-age'"

echo -e "\033[1;35m3. Sort by Multiple Columns (department, then age)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --sort 'department,age'"

echo -e "\033[1;35m4. Sort by Multiple Columns with Mixed Order\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --sort 'department,-age'"

# =============================================================================
# END OF DEMO
# =============================================================================

divider
echo -e "\033[1;32m✅ SORTING DEMO COMPLETE!\033[0m"