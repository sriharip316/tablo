#!/bin/bash

# Flattening & Diving Demo - Tablo CLI
# Comprehensive demonstration of nested data handling and flattening features

set -eou pipefail

source "$(dirname "${BASH_SOURCE[0]}")/_common.sh"

echo -e "\033[1;32m🏗️ TABLO CLI - FLATTENING & DIVING DEMO\033[0m"

# =============================================================================
# WHAT IS FLATTENING?
# =============================================================================

section "🤔 What is Flattening?"

echo "Flattening converts nested JSON/YAML into flat table rows."
echo ""
echo "Without flattening:"
echo "  {\"user\": {\"name\": \"Alice\", \"age\": 30}} → [user] column"
echo ""
echo "With flattening (--dive):"
echo "  {\"user\": {\"name\": \"Alice\", \"age\": 30}} → [user.name] [user.age] columns"
echo ""

# =============================================================================
# BASIC FLATTENING
# =============================================================================

section "🎯 Basic Flattening Examples"

echo -e "\033[1;35m1. Simple Nested Object (Without --dive)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\""

echo -e "\033[1;35m2. Same Data With Flattening (-d|--dive)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive"

echo -e "\033[1;35m4. Array of Objects With Nested Data\033[0m"
echo "Without flattening - nested objects shown as complex values:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/list.json\""

echo "With flattening - nested objects become separate columns:"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/list.json\" --dive"

# =============================================================================
# DEPTH CONTROL
# =============================================================================

section "📏 Controlling Flattening Depth"

echo "Complex nested data can create too many columns. Use -m|--max-depth to control this."
echo ""

echo -e "\033[1;35m1. Deep Nested Data - No Depth Limit\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/nested.json\" --dive --max-col-width 60 --wrap word"

echo -e "\033[1;35m2. Limit to Depth 1 (--max-depth 1)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/nested.json\" --dive --max-col-width 60 --wrap word --max-depth 1"

echo -e "\033[1;35m3. Limit to Depth 2 (--max-depth 2)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/nested.json\" --dive --max-col-width 60 --wrap word --max-depth 2"

echo -e "\033[1;35m4. Limit to Depth 3 (--max-depth 3)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/nested.json\" --dive --max-col-width 60 --wrap word --max-depth 3"

# =============================================================================
# SELECTIVE DIVING
# =============================================================================

section "🎯 Selective Diving (-D|--dive-path)"

echo "Sometimes you only want to flatten specific parts of your data."
echo "Use --dive-path to specify which top-level paths to dive into."
echo ""

echo -e "\033[1;35m1. Show Available Top-Level Paths\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive --max-depth 1"

echo -e "\033[1;35m2. Dive Only Into 'meta' Path\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive --dive-path meta"

echo -e "\033[1;35m3. Dive Into Multiple Specific Paths\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive --dive-path meta --dive-path settings"

# =============================================================================
# ARRAY FLATTENING
# =============================================================================

section "📋 Array Flattening"

echo "Arrays can be handled in different ways depending on their content."
echo ""

echo -e "\033[1;35m1. Arrays of Objects (Default Behavior)\033[0m"
echo "Arrays of objects become separate rows:"
demo_cmd "$TABLO_BIN --input '[{\"name\":\"Alice\",\"skills\":[{\"name\":\"Go\",\"level\":\"expert\"},{\"name\":\"SQL\",\"level\":\"intermediate\"}]}]' --dive"

echo -e "\033[1;35m2. Arrays of Simple Values\033[0m"
echo "Without --flatten-simple-arrays:"
demo_cmd "$TABLO_BIN --input '{\"user\":\"Alice\",\"tags\":[\"engineer\",\"golang\",\"sql\"]}' --dive"

echo "With --flatten-simple-arrays:"
demo_cmd "$TABLO_BIN --input '{\"user\":\"Alice\",\"tags\":[\"engineer\",\"golang\",\"sql\"]}' --dive --flatten-simple-arrays"

echo -e "\033[1;35m4. Real-world Example: Product Tags\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/products.json\" --dive --flatten-simple-arrays --select 'name,tags' --limit 3"

# =============================================================================
# WORKING WITH COMPLEX NESTED DATA
# =============================================================================

section "🏢 Complex Nested Data Examples"

echo "Let's explore a realistic complex nested structure:"
echo ""

echo -e "\033[1;35m1. Company Structure Overview (Depth 1)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/nested.json\" --max-col-width 60 --wrap word --dive --max-depth 1"

echo -e "\033[1;35m2. Department Information (Depth 2)\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/nested.json\" --max-col-width 60 --wrap word --dive --dive-path company --max-depth 2 --select '*.departments.*' --limit 5"

echo -e "\033[1;35m3. Financial Data Analysis\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/nested.json\" --max-col-width 60 --wrap word --dive --dive-path company --max-depth 3 --select '*.financial.*.*' --limit 10"

echo -e "\033[1;35m4. Technology Stack Overview\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/nested.json\" --max-col-width 60 --wrap word --dive --dive-path company --max-depth 4 --select '*.technology.*.*.*' --limit 8"

# =============================================================================
# COMBINING FLATTENING WITH OTHER OPTIONS
# =============================================================================

section "🔧 Combining Flattening with Other Features"

echo -e "\033[1;35m1. Flattening + Selection\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/single.json\" --dive --select 'id,name,meta.owner.*'"

echo -e "\033[1;35m2. Flattening + Formatting\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/list.json\" --dive --style markdown --limit 2"

echo -e "\033[1;35m3. Flattening + Export\033[0m"
TEMP_FILE="/tmp/tablo_flattened.csv"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/list.json\" --dive --style csv --output \"$TEMP_FILE\" --limit 2"
echo "Flattened data exported to CSV:"
head -5 "$TEMP_FILE"
rm -f "$TEMP_FILE"

echo -e "\033[1;35m4. Flattening + Index Column\033[0m"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/data/list.json\" --dive --index-column --limit 2"

# =============================================================================
# PRACTICAL FLATTENING PATTERNS
# =============================================================================

section "💼 Practical Flattening Patterns"

echo -e "\033[1;35m1. API Response Analysis\033[0m"
echo "Simulate analyzing a complex API response:"
demo_cmd "echo '{\"data\":{\"users\":[{\"id\":1,\"profile\":{\"name\":\"Alice\",\"settings\":{\"theme\":\"dark\",\"notifications\":true}}}]}}' | $TABLO_BIN --dive --max-depth 3"

echo -e "\033[1;35m2. Configuration File Exploration\033[0m"
echo "Explore nested configuration structures:"
CONFIG_DATA='{
  "app": {
    "name": "MyApp",
    "version": "1.0.0",
    "database": {
      "host": "localhost",
      "port": 5432,
      "credentials": {
        "username": "admin",
        "password": "***"
      }
    },
    "features": {
      "auth": {"enabled": true, "provider": "oauth2"},
      "cache": {"enabled": true, "ttl": 3600}
    }
  }
}'
demo_cmd "$TABLO_BIN --input '$CONFIG_DATA' --dive --max-depth 3"

echo -e "\033[1;35m3. Log Analysis\033[0m"
echo "Flatten structured log entries:"
LOG_DATA='[
  {
    "timestamp": "2023-10-15T10:30:00Z",
    "level": "INFO",
    "message": "User logged in",
    "context": {
      "user_id": 123,
      "ip": "192.168.1.1",
      "user_agent": "Chrome/118.0.0.0"
    }
  },
  {
    "timestamp": "2023-10-15T10:31:00Z",
    "level": "ERROR",
    "message": "Database connection failed",
    "context": {
      "database": "users",
      "error": "connection timeout",
      "retry_count": 3
    }
  }
]'
demo_cmd "$TABLO_BIN --input '$LOG_DATA' --dive --style borderless"

echo -e "\033[1;35m4. Monitoring Data\033[0m"
echo "Analyze nested monitoring metrics:"
METRICS_DATA='{
  "services": [
    {
      "name": "web-server",
      "status": "healthy",
      "metrics": {
        "cpu": {"usage": 45.2, "limit": 80},
        "memory": {"usage": 1024, "limit": 2048},
        "requests": {"total": 15000, "errors": 12}
      }
    },
    {
      "name": "database",
      "status": "degraded",
      "metrics": {
        "cpu": {"usage": 78.9, "limit": 80},
        "memory": {"usage": 3900, "limit": 4096},
        "connections": {"active": 45, "max": 50}
      }
    }
  ]
}'
demo_cmd "$TABLO_BIN --input '$METRICS_DATA' --dive --max-depth 2 --precision 1"

# =============================================================================
# QUICK REFERENCE
# =============================================================================

section "📋 Quick Reference"

echo -e "\033[1;35mFlattening Options:\033[0m"
echo "  --dive, -d                    Enable flattening"
echo "  --max-depth N, -m N           Limit flattening depth"
echo "  --dive-path PATH, -D PATH     Flatten specific paths only"
echo "  --flatten-simple-arrays       Convert primitive arrays to comma-separated strings"
echo ""

echo -e "\033[1;35mCommon Patterns:\033[0m"
echo "  tablo --file data.json --dive                    # Basic flattening"
echo "  tablo --file data.json --dive --max-depth 2      # Limited depth"
echo "  tablo --file data.json --dive --dive-path users  # Specific section"
echo "  tablo --file data.json --dive --flatten-simple-arrays  # Handle primitive arrays"
echo ""

echo -e "\033[1;35mTroubleshooting:\033[0m"
echo "  • Too many columns? Use --max-depth or --dive-path"
echo "  • Missing data? Check depth limits and paths"
echo "  • Performance issues? Use --limit and reduce depth"
echo "  • Array problems? Use --flatten-simple-arrays for primitives"

# =============================================================================
# END OF DEMO
# =============================================================================

divider
echo -e "\033[1;32m✅ FLATTENING & DIVING DEMO COMPLETE!\033[0m"
