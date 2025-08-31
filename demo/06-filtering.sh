#!/bin/bash

# Filtering Demo - Tablo CLI
# Comprehensive demonstration of the --where filtering feature

set -eou pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

echo -e "\033[1;32m🔍 TABLO CLI - FILTERING WITH --where\033[0m"

# =============================================================================
# WHAT IS FILTERING?
# =============================================================================

section "🤔 What is Filtering?"

echo "The --where option lets you filter table rows based on conditions."
echo ""
echo "Basic syntax:"
echo "  tablo --file data.json --where 'column=value'"
echo "  tablo --file data.json --where 'age>30' --where 'active=true'"
echo ""
echo "Multiple --where conditions use AND logic (all must be true)."
echo ""
echo "First, let's see our sample data:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --select 'name,age,department,active,salary' --limit 5"

# =============================================================================
# EQUALITY FILTERING
# =============================================================================

section "⚖️ Equality Filtering"

echo -e "\033[1;35m1. Exact String Match (=)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'department=Engineering' --select 'name,department,salary'"

echo -e "\033[1;35m2. Boolean Filtering\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'active=true' --select 'name,active,department'"

echo -e "\033[1;35m3. Exact Number Match\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'age=28' --select 'name,age,city'"

echo -e "\033[1;35m4. Not Equal (!=)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'department!=Engineering' --select 'name,department'"

# =============================================================================
# NUMERIC COMPARISONS
# =============================================================================

section "🔢 Numeric Comparisons"

echo -e "\033[1;35m1. Greater Than (>)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'age>30' --select 'name,age,department'"

echo -e "\033[1;35m2. Greater Than or Equal (>=)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'salary>=75000' --select 'name,salary,department'"

echo -e "\033[1;35m3. Less Than (<)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'age<30' --select 'name,age,city'"

echo -e "\033[1;35m4. Less Than or Equal (<=)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'score<=90' --select 'name,score,department'"

echo -e "\033[1;35m5. Decimal Precision\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'score>90.0' --select 'name,score' --precision 1"

# =============================================================================
# STRING OPERATIONS
# =============================================================================

section "📝 String Operations"

echo -e "\033[1;35m1. Contains (~)\033[0m"
echo "Find users whose names contain 'son':"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'name~son' --select 'name,department'"

echo -e "\033[1;35m2. Department Contains\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'department~Eng' --select 'name,department'"

echo -e "\033[1;35m3. Email Domain Check\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'email~example.com' --select 'name,email'"

echo -e "\033[1;35m4. Does Not Contain (!~)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'city!~San' --select 'name,city'"

echo -e "\033[1;35m5. Case Sensitivity\033[0m"
echo "String matching is case-sensitive:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'department~eng' --select 'name,department'" || echo "No matches (case-sensitive)"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'department~Eng' --select 'name,department'"

# =============================================================================
# REGEX PATTERN MATCHING
# =============================================================================

section "🎯 Regular Expression Matching"

echo -e "\033[1;35m1. Pattern Match (=~)\033[0m"
echo "Find emails matching a pattern:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'email=~.*@example\.com' --select 'name,email'"

echo -e "\033[1;35m2. Name Patterns\033[0m"
echo "Names starting with vowels:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'name=~^[AEIOU]' --select 'name'"

echo -e "\033[1;35m3. Complex Patterns\033[0m"
echo "Names with exactly two words:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'name=~^[A-Z][a-z]+ [A-Z][a-z]+$' --select 'name'"

echo -e "\033[1;35m4. Negative Pattern Match (!=~)\033[0m"
echo "Names not starting with 'A':"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'name!=~^A' --select 'name'"

echo -e "\033[1;35m5. Number Patterns\033[0m"
echo "IDs that are single digits:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'id=~^[1-9]$' --select 'id,name'"

# =============================================================================
# MULTIPLE CONDITIONS (AND LOGIC)
# =============================================================================

section "🔗 Multiple Conditions (AND Logic)"

echo -e "\033[1;35m1. Two Conditions\033[0m"
echo "Active Engineering employees:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'department=Engineering' --where 'active=true' --select 'name,department,active'"

echo -e "\033[1;35m2. Salary Range\033[0m"
echo " Filter by salary range (Salary > 1000 AND salary < 70000):"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'salary>1000' --where 'salary<70000' --select 'name,age,salary,department'"

echo -e "\033[1;35m2. Age and Salary Range\033[0m"
echo "Young high earners (age < 35 AND salary > 70000):"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/users.json\" --where 'age<35' --where 'salary>70000' --select 'name,age,salary,department'"

# =============================================================================
# QUICK REFERENCE
# =============================================================================

section "📋 Quick Reference"

echo -e "\033[1;35mOperators:\033[0m"
echo "  =      Equal to"
echo "  !=     Not equal to"
echo "  >      Greater than"
echo "  >=     Greater than or equal"
echo "  <      Less than"
echo "  <=     Less than or equal"
echo "  ~      Contains (substring)"
echo "  !~     Does not contain"
echo "  =~     Matches regex pattern"
echo "  !=~    Does not match regex pattern"
echo ""

echo -e "\033[1;35mExamples:\033[0m"
echo "  --where 'name=Alice'              # Exact match"
echo "  --where 'age>30'                  # Numeric comparison"
echo "  --where 'active=true'             # Boolean"
echo "  --where 'email~@gmail.com'        # Contains"
echo "  --where 'name=~^[A-M]'            # Regex (names A-M)"
echo "  --where 'dept=Engineering' --where 'age<35'  # Multiple (AND)"
echo ""

echo -e "\033[1;35mCombination Examples:\033[0m"
echo "  # Active engineers over 30"
echo "  tablo --file users.json --where 'department=Engineering' --where 'active=true' --where 'age>30'"
echo ""
echo "  # High-value products in stock"
echo "  tablo --file products.json --where 'price>200' --where 'in_stock=true' --style csv"
echo ""
echo "  # Export filtered data"
echo "  tablo --file data.json --where 'category=premium' --output filtered.csv --style csv"

# =============================================================================
# END OF DEMO
# =============================================================================

divider
echo -e "\033[1;32m✅ FILTERING DEMO COMPLETE!\033[0m"
