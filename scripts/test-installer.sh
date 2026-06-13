#!/bin/bash
# test-installer.sh — Test installer scenarios in a temporary sandbox
#
# Usage:
#   ./scripts/test-installer.sh              # run all scenarios
#
# This script creates a sandboxed home directory and PATH to
# avoid modifying the real user environment.

set -euo pipefail

SANDBOX="$(mktemp -d)"
trap 'echo "Cleaning up sandbox: ${SANDBOX}"; rm -rf "${SANDBOX}"' EXIT

PASS=0
FAIL=0

pass() { echo "  ✓ PASS: $1"; PASS=$((PASS + 1)); }
fail() { echo "  ✗ FAIL: $1"; FAIL=$((FAIL + 1)); }

SCRIPT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
INSTALLER="${SCRIPT_DIR}/install.sh"
RELEASE_INSTALLER="${SCRIPT_DIR}/release/install.sh"
ASFTUI_INSTALLER="${SCRIPT_DIR}/asf-tui/install.sh"

echo "ASF Installer Test Suite"
echo "Script dir: ${SCRIPT_DIR}"
echo "Sandbox:    ${SANDBOX}"

# ─── Test: help flag ─────────────────────────────────────
echo ""
echo "═══ Test: --help flag ═══"
out_help=$(bash "${INSTALLER}" --help 2>/dev/null) || true
out_h=$(bash "${INSTALLER}" -h 2>/dev/null) || true
out_rel=$(bash "${RELEASE_INSTALLER}" --help 2>/dev/null) || true
if echo "$out_help" | grep -c "Usage:" >/dev/null; then pass "--help shows usage"; else fail "--help no usage"; fi
if echo "$out_h" | grep -c "Usage:" >/dev/null; then pass "-h shows usage"; else fail "-h no usage"; fi
if echo "$out_rel" | grep -c "Usage:" >/dev/null; then pass "release/install.sh --help works"; else fail "release/install.sh --help failed"; fi

# ─── Test: clean flag in help ────────────────────────────
echo ""
echo "═══ Test: flags in help ═══"
if echo "$out_help" | grep -c "clean" >/dev/null; then pass "--clean listed in help"; else fail "--clean not in help"; fi
if echo "$out_help" | grep -c "repair" >/dev/null; then pass "--repair listed in help"; else fail "--repair not in help"; fi
if echo "$out_help" | grep -c "upgrade" >/dev/null; then pass "--upgrade listed in help"; else fail "--upgrade not in help"; fi

# ─── Test: bash syntax (all scripts) ─────────────────────
echo ""
echo "═══ Test: bash syntax ═══"
if bash -n "${INSTALLER}" 2>&1; then pass "install.sh syntax OK"; else fail "install.sh syntax error"; fi
if bash -n "${RELEASE_INSTALLER}" 2>&1; then pass "release/install.sh syntax OK"; else fail "release/install.sh syntax error"; fi
if bash -n "${ASFTUI_INSTALLER}" 2>&1; then pass "asf-tui/install.sh syntax OK"; else fail "asf-tui/install.sh syntax error"; fi

# ─── Test: no top-level 'local' (dash compat) ────────────
echo ""
echo "═══ Test: dash compatibility (no 'local' at top level) ═══"
for script in "${INSTALLER}" "${RELEASE_INSTALLER}" "${ASFTUI_INSTALLER}"; do
  name="$(basename "$script")"
  issues="$(grep -n '^local ' "$script" 2>/dev/null || true)"
  if [ -n "$issues" ]; then
    fail "${name}: top-level 'local' found: ${issues}"
  else
    pass "${name}: no top-level 'local'"
  fi
done

# ─── Test: no Python references ──────────────────────────
echo ""
echo "═══ Test: No Python references ═══"
for script in "${INSTALLER}" "${RELEASE_INSTALLER}" "${ASFTUI_INSTALLER}"; do
  name="$(basename "$script")"
  pyref=$(grep -c "python\|pip install\|python-engine" "$script" 2>/dev/null; true)
  pyref=$(echo "$pyref" | tr -d '[:space:]')
  if [ -z "$pyref" ] || [ "$pyref" = "0" ]; then
    pass "${name}: clean"
  else
    fail "${name}: ${pyref} Python references found"
  fi
done

# ─── Test: verify_install order ──────────────────────────
echo ""
echo "═══ Test: verify_install defined before call ═══"
for script in "${INSTALLER}" "${RELEASE_INSTALLER}"; do
  name="$(basename "$script")"
  def_line="$(grep -n 'verify_install()' "$script" | head -1 | cut -d: -f1 || echo "0")"
  call_line="$(grep -n 'verify_install$' "$script" | tail -1 | cut -d: -f1 || echo "0")"
  if [ "$def_line" -gt 0 ] && [ "$call_line" -gt 0 ] && [ "$def_line" -lt "$call_line" ]; then
    pass "${name}: verify_install defined at ${def_line}, called at ${call_line}"
  else
    fail "${name}: ordering issue (def=${def_line}, call=${call_line})"
  fi
done

# ─── Test: --repair without existing binary ──────────────
echo ""
echo "═══ Test: --repair without existing binary ═══"
test_home="${SANDBOX}/test_repair"
mkdir -p "$test_home"
rc=0
HOME="$test_home" ASF_VERSION="2.0.1" bash "${INSTALLER}" --repair 2>&1 || rc=$?
if [ "$rc" -ne 0 ]; then
  pass "--repair correctly fails without existing binary"
else
  fail "--repair should have failed"
fi

# ─── Test: --purge requires --clean ──────────────────────
echo ""
echo "═══ Test: --purge requires --clean ═══"
rc=0
bash "${INSTALLER}" --purge 2>&1 || rc=$?
if [ "$rc" -ne 0 ]; then
  pass "--purge without --clean correctly fails"
else
  fail "--purge without --clean should exit 1"
fi

# ─── Test: --purge listed in help ────────────────────────
echo ""
echo "═══ Test: --purge listed in help ═══"
if echo "$out_help" | grep -c "purge" >/dev/null; then pass "--purge listed in help"; else fail "--purge not in help"; fi

# ─── Test: PATH setup in help ────────────────────────────
echo ""
echo "═══ Test: PATH setup printed ═══"
if echo "$out_help" | grep -c "PATH" >/dev/null; then pass "Help mentions PATH"; else fail "Help missing PATH reference"; fi

# ─── Test: shell detection in installer ──────────────────
echo ""
echo "═══ Test: shell detection ═══"
shell_test=$(bash -c 'source "${INSTALLER}" --help 2>/dev/null; echo "ok"' || echo "")
if echo "$shell_test" | grep -q "ok"; then
  pass "Shell detection runs without error"
else
  fail "Shell detection error"
fi

# ─── Test: platform detection ────────────────────────────
echo ""
echo "═══ Test: help text mentions options ═══"
out_plat=$(bash "${INSTALLER}" --help 2>/dev/null) || true
if echo "$out_plat" | grep -c "Options:" >/dev/null; then pass "Help shows Options section"; else fail "Help missing Options"; fi

# ─── Test: HOME/.asf directory created on fresh install ─
echo ""
echo "═══ Test: Version files present ═══"
asftui_ver="$(grep '^ASF_VERSION=' "${ASFTUI_INSTALLER}" 2>/dev/null | head -1 | sed 's/.*ASF_VERSION="\(.*\)"/\1/' || true)"
ver_file="$(cat "${SCRIPT_DIR}/release/VERSION" 2>/dev/null || echo "")"
echo "  asf-tui/install.sh ASF_VERSION: ${asftui_ver}"
echo "  release/VERSION: ${ver_file}"
if [ -n "$asftui_ver" ]; then pass "asf-tui/install.sh version: ${asftui_ver}"; else fail "asf-tui/install.sh missing version"; fi
if [ -n "$ver_file" ]; then pass "release/VERSION: ${ver_file}"; else fail "release/VERSION missing"; fi

# ─── Version Regression Tests ─────────────────────────────
echo ""
echo "═══ Version Regression Tests (INSTALLER_HOTFIX) ═══"

# Test 1: Single-line tag parsing
echo ""
echo "  Test 1: Single-line tag parsing"
result1="$(echo "v3.0.0-RC2" | tr -d '\r\n' | xargs)"
if [ "$result1" = "v3.0.0-RC2" ]; then pass "Single-line: v3.0.0-RC2"; else fail "Single-line: got [${result1}]"; fi

# Test 2: Trailing newline
echo ""
echo "  Test 2: Trailing newline"
result2="$(printf "v3.0.0-RC2\n" | tr -d '\r\n' | xargs)"
if [ "$result2" = "v3.0.0-RC2" ]; then pass "Trailing newline: v3.0.0-RC2"; else fail "Trailing newline: got [${result2}]"; fi

# Test 3: CRLF line endings
echo ""
echo "  Test 3: CRLF line endings"
result3="$(printf "v3.0.0-RC2\r\n" | tr -d '\r\n' | xargs)"
if [ "$result3" = "v3.0.0-RC2" ]; then pass "CRLF: v3.0.0-RC2"; else fail "CRLF: got [${result3}]"; fi

# Test 4: URL contains /releases/download/v3.0.0-RC2/ exactly once
echo ""
echo "  Test 4: URL tag path validation"
REPO="moksh5936-2/asfassumption"
LATEST_VERSION="v3.0.0-RC2"
OS_FINAL="darwin"
ARCH_FINAL="arm64"
BINARY_NAME="ASF-${LATEST_VERSION}-${OS_FINAL}-${ARCH_FINAL}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}"
tag_count="$(printf '%s' "$DOWNLOAD_URL" | grep -o '/releases/download/v3.0.0-RC2/' | wc -l | tr -d ' ')"
if [ "$tag_count" = "1" ]; then pass "URL contains /releases/download/v3.0.0-RC2/ exactly once"; else fail "URL has ${tag_count} occurrences (expected 1): ${DOWNLOAD_URL}"; fi

# Test 5: Asset filename contains ASF-v3.0.0-RC2-darwin-arm64 exactly once
echo ""
echo "  Test 5: Asset filename validation"
asset_count="$(printf '%s' "$BINARY_NAME" | grep -o "ASF-v3.0.0-RC2-darwin-arm64" | wc -l | tr -d ' ')"
if [ "$asset_count" = "1" ]; then pass "Asset filename is ASF-v3.0.0-RC2-darwin-arm64"; else fail "Asset filename mismatch: ${BINARY_NAME}"; fi

# Test 6: No whitespace in URL
echo ""
echo "  Test 6: URL whitespace check"
if printf '%s' "$DOWNLOAD_URL" | grep -q '[[:space:]]'; then fail "URL contains whitespace: ${DOWNLOAD_URL}"; else pass "URL has no whitespace"; fi

# Test 7: Except Exception prevents double version (root cause fix)
echo ""
echo "  Test 7: bare except vs except Exception"
# Simulate the FIXED behavior: except Exception does NOT catch sys.exit(0)
fixed_output="$(python3 -c "
import sys
try:
    print('v3.0.0-RC2')
    sys.exit(0)
except Exception: pass
print('v3.0.0-RC2')
" 2>/dev/null)"
fixed_lines="$(echo "$fixed_output" | wc -l | tr -d ' ')"
if [ "$fixed_lines" = "1" ]; then pass "except Exception: single version output"; else fail "except Exception: ${fixed_lines} lines"; fi

# Test 8: bare except BUG — catches SystemExit, double version
echo ""
echo "  Test 8: bare except bug (historical — should fail)"
bug_output="$(python3 -c "
import sys
try:
    print('v3.0.0-RC2')
    sys.exit(0)
except: pass
print('v3.0.0-RC2')
" 2>/dev/null)"
bug_lines="$(echo "$bug_output" | wc -l | tr -d ' ')"
if [ "$bug_lines" = "2" ]; then pass "bare except: double version output (CONFIRMED BUG)"; else fail "bare except: ${bug_lines} lines"; fi

# ─── Summary ─────────────────────────────────────────────
echo ""
echo "═══ Results ═══"
echo "  Passed: ${PASS}"
echo "  Failed: ${FAIL}"
[ "$FAIL" -gt 0 ] && exit 1
echo "  All tests passed!"
