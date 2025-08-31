#!/bin/bash

# Common environment and helper functions for Tablo demo scripts

# Environment setup
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEMO_DIR="$(cd "$SCRIPT_DIR" && pwd)"
PROJECT_ROOT="$(cd "$DEMO_DIR/.." && pwd)"
TABLO_BIN="$PROJECT_ROOT/bin/tablo"
COLUMNS=$(tput cols 2>/dev/null || echo 80)

# Build tablo CLI if not present
if ! [ -x "$TABLO_BIN" ]; then
    echo "Building tablo CLI..."
    (cd "$PROJECT_ROOT" && make build)
fi

# Helper functions

divider() {
    printf '%*s\n' "${COLUMNS}" '' | tr ' ' '-'
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
