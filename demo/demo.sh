#!/bin/bash

# Demo script for tablo CLI
# Showcases input types, formats, output styles, selection, exclusion, flattening, formatting, and more.
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

section "Demo: Basic file input (JSON)"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\""

section "Demo: Basic file input (YAML)"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.yaml\""

section "Demo: List input (JSON array)"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\""

section "Demo: List input (YAML array)"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.yaml\""

section "Demo: Raw input string"
demo_cmd "$TABLO_BIN --input '{\"a\":1,\"b\":2,\"c\":{\"d\":3}}'"

section "Demo: Flatten nested objects"
demo_cmd "$TABLO_BIN --input '{\"a\":{\"b\":1},\"tags\":[\"x\",\"y\",3]}' --dive --flatten-simple-arrays"

section "Demo: Select columns (name, age) and index column"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.yaml\" --select 'name,age' --index-column --style ascii"

section "Demo: Exclude columns (meta, misc)"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --exclude 'meta,misc' --style markdown"

section "Demo: Limit rows and markdown style"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --limit 2 --style markdown"

section "Demo: Custom boolean, null, and float formatting"
demo_cmd "$TABLO_BIN --input '[{\"a\":1.2345,\"b\":true},{\"b\":false}]' --style ascii --precision 2 --bool-str 'Y:N' --index-column"

section "Demo: Dive only into specific paths"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --dive-path meta --dive-path settings --dive --max-depth 2"

section "Demo: Output styles (heavy, light, double, ascii, markdown, compact, borderless)"
for style in heavy light double ascii markdown compact borderless; do
    echo -e "\033[1;35mStyle: $style\033[0m"
    "$TABLO_BIN" --file "$DEMO_DIR/list.json" --style "$style" --limit 1
    echo
done

section "Demo: Select columns using glob pattern"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --select 'meta.owner.*'"

section "Demo: Exclude columns using glob pattern"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --exclude 'meta.*'"

section "Demo: No header row"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --no-header"

section "Demo: Header case (upper, lower, title)"
for case in upper lower title; do
    echo -e "\033[1;35mHeader case: $case\033[0m"
    "$TABLO_BIN" --file "$DEMO_DIR/list.json" --header-case "$case" --limit 1
    echo
done

section "Demo: Max column width and cell wrapping"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/single.json\" --max-col-width 20 --wrap word"

section "Demo: Output to file"
OUTFILE="$DEMO_DIR/demo_output.txt"
demo_cmd "$TABLO_BIN --file \"$DEMO_DIR/list.json\" --style ascii --output \"$OUTFILE\""
echo "Output written to $OUTFILE"
cat "$OUTFILE"
rm -f "$OUTFILE"
divider

section "Demo: Piped input"
demo_cmd "echo '{\"name\":\"Alice\",\"age\":30,\"city\":\"Wonderland\"}' | $TABLO_BIN --style ascii"

divider
echo -e "\033[1;32mAll demos completed!\033[0m"
divider
