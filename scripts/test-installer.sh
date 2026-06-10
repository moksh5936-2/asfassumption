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
HOME="$test_home" ASF_VERSION="2.0.0" bash "${INSTALLER}" --repair 2>&1 || rc=$?
if [ "$rc" -ne 0 ]; then
  pass "--repair correctly fails without existing binary"
else
  fail "--repair should have failed"
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

# ─── Summary ─────────────────────────────────────────────
echo ""
echo "═══ Results ═══"
echo "  Passed: ${PASS}"
echo "  Failed: ${FAIL}"
[ "$FAIL" -gt 0 ] && exit 1
echo "  All tests passed!"
