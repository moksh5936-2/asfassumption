#!/bin/bash
# test-commands.sh — Smoke test every documented ASF command
#
# Usage:
#   ./scripts/test-commands.sh <path-to-asf-binary>
#
# Default: tests the built binary at asf-tui/asf-tui

set -euo pipefail

ASF="${1:-}"
if [ -z "$ASF" ]; then
  SCRIPT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
  ASF="${SCRIPT_DIR}/asf-tui/asf-tui"
fi

if [ ! -x "$ASF" ]; then
  echo "Error: binary not found or not executable: $ASF"
  echo "Usage: $0 <path-to-asf-binary>"
  exit 1
fi

SAMPLE_DIR="$(cd "$(dirname "$0")/../sample_data" && pwd)"

PASS=0
FAIL=0

pass() { echo "  ✓ PASS: $1"; PASS=$((PASS + 1)); }
fail() { echo "  ✗ FAIL: $1 (exit $2)"; FAIL=$((FAIL + 1)); }

expect() {
  local name="$1" expected="$2"
  shift 2
  local out tmpdir
  tmpdir=$(mktemp -d)
  out="${tmpdir}/out"
  set +e
  "$@" > "$out" 2>&1
  local rc=$?
  set -e
  if [ "$rc" -eq "$expected" ]; then
    pass "$name"
  else
    fail "$name (expected exit $expected, got $rc)" "$rc"
    head -5 "$out"
  fi
  rm -rf "$tmpdir"
}

# Capture --version output
VER_OUT=$("$ASF" --version 2>/dev/null)
echo "Testing: ${VER_OUT:-ASF binary}"
echo "Sample: ${SAMPLE_DIR}"
echo ""

echo "═══ Basic commands ═══"
expect "--version" 0 "$ASF" --version
expect "-v" 0 "$ASF" -v
expect "--help" 0 "$ASF" --help
expect "-h" 0 "$ASF" -h

echo ""
echo "═══ doctor commands ═══"
expect "doctor" 0 "$ASF" doctor
expect "doctor --verbose" 0 "$ASF" doctor --verbose
expect "doctor --fix" 0 "$ASF" doctor --fix

echo ""
echo "═══ analyze commands ═══"
expect "analyze --help" 0 "$ASF" analyze --help
expect "analyze -h" 0 "$ASF" analyze -h

echo ""
echo "═══ analyze with file ═══"
TXT="${SAMPLE_DIR}/finance_policy.txt"
if [ -f "$TXT" ]; then
  expect "analyze txt file" 0 "$ASF" analyze "$TXT"
  expect "analyze txt --json" 0 "$ASF" analyze "$TXT" --json
  expect "analyze txt --graph" 0 "$ASF" analyze "$TXT" --graph
  expect "analyze txt --json --graph" 0 "$ASF" analyze "$TXT" --json --graph
else
  echo "  ⚠  sample TXT not found, skipping file tests"
fi

echo ""
echo "═══ analyze with directory ═══"
if [ -d "$SAMPLE_DIR" ]; then
  expect "analyze directory" 0 "$ASF" analyze "$SAMPLE_DIR"
else
  echo "  ⚠  sample directory not found, skipping"
fi

echo ""
echo "═══ analyze with evidence ═══"
CSV="${SAMPLE_DIR}/backup_config.csv"
if [ -f "$TXT" ] && [ -f "$CSV" ]; then
  expect "analyze -e short" 0 "$ASF" analyze "$TXT" -e "$CSV"
  expect "analyze --evidence long" 0 "$ASF" analyze "$TXT" --evidence "$CSV"
fi

echo ""
echo "═══ edge cases ═══"
expect "invalid command" 1 "$ASF" invalid-command
expect "missing file" 1 "$ASF" analyze /tmp/nonexistent-file.txt
expect "analyze no args" 1 "$ASF" analyze

echo ""
echo "═══ Results ═══"
echo "  Passed: ${PASS}"
echo "  Failed: ${FAIL}"
[ "$FAIL" -gt 0 ] && exit 1
echo "  All commands passed!"
