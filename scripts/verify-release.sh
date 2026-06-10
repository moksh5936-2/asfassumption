#!/bin/bash
# ASF Release Verification Script
# Run after build to verify the binary is functional.
# Usage: ./scripts/verify-release.sh <path-to-binary>
#
# Exit codes:
#   0 — all checks pass
#   1 — any check fails

set -euo pipefail

BINARY="${1:-./asf-tui/asf-tui}"

# If binary is a bare filename (no path), prepend ./
case "$BINARY" in
  */*) ;;
  *) BINARY="./$BINARY" ;;
esac

if [ ! -f "$BINARY" ]; then
  echo "✗ Binary not found: $BINARY"
  echo "  Checked from: $(pwd)"
  ls -la . 2>/dev/null | head -5
  exit 1
fi

if [ ! -x "$BINARY" ]; then
  chmod +x "$BINARY" 2>/dev/null || true
fi

echo "=== ASF Release Verification ==="
echo "Binary: $BINARY"
echo "Size:   $(ls -lh "$BINARY" | awk '{print $5}')"
echo ""

FAILED=0

check() {
  local name="$1"
  shift
  echo -n "  ✓ $name ... "
  local tmpout tmpdir
  tmpdir=$(mktemp -d)
  tmpout="${tmpdir}/output"
  if "$@" >"$tmpout" 2>&1; then
    echo "PASS"
  else
    echo "FAIL"
    cat "$tmpout"
    FAILED=1
  fi
  rm -rf "$tmpdir"
}

check "binary launches" "$BINARY" --version
check "version flag" "$BINARY" -v
check "help flag" "$BINARY" --help
check "doctor runs" "$BINARY" doctor
check "doctor verbose" "$BINARY" doctor --verbose

echo ""

if [ "$FAILED" -eq 0 ]; then
  VER=$("$BINARY" --version 2>/dev/null)
  echo "  ✓ All checks passed — $VER ready for release."
else
  echo "  ✗ Some checks failed."
  exit 1
fi
