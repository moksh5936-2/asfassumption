#!/bin/bash
# ASF — Architecture Security Framework
# Installer — https://github.com/moksh5936-2/asfassumption
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
#   curl ... | bash -s -- --upgrade
#   curl ... | bash -s -- --repair
#   curl ... | bash -s -- --clean
#   curl ... | bash -s -- --clean --purge
#
# Environment:
#   ASF_VERSION=2.0.1       — pin a specific version
#   GITHUB_TOKEN=ghp_xxx    — for private repos (set via gh auth token)
#   ASF_INSTALL_DIR=        — custom install directory (default: ~/.local/bin)

set -euo pipefail

# ─── Config ────────────────────────────────────────────────
REPO="moksh5936-2/asfassumption"
VERSION="${ASF_VERSION:-}"
ASF_HOME="${HOME}/.asf"
BACKUP_DIR="${ASF_HOME}/backups"
INSTALL_DIR="${ASF_INSTALL_DIR:-${HOME}/.local/bin}"

# ─── Parse flags ──────────────────────────────────────────
SHOW_HELP=false
UPGRADE=false
REPAIR=false
CLEAN=false
PURGE=false
for arg in "$@"; do
  case "$arg" in
    --help|-h) SHOW_HELP=true ;;
    --upgrade|-u) UPGRADE=true ;;
    --repair) REPAIR=true ;;
    --clean) CLEAN=true ;;
    --purge) PURGE=true ;;
  esac
done

if [ "$PURGE" = true ] && [ "$CLEAN" = false ]; then
  echo "Error: --purge must be used with --clean" >&2
  echo "Usage: curl ... | bash -s -- --clean --purge" >&2
  exit 1
fi

# ─── Platform detection ────────────────────────────────────
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin) OS_FINAL="darwin" ;;
  Linux)  OS_FINAL="linux"  ;;
  *)
    echo "Unsupported OS: ${OS}. Windows users: run install.ps1 in PowerShell." >&2
    exit 1
    ;;
esac

case "$ARCH" in
  x86_64|amd64) ARCH_FINAL="amd64"   ;;
  arm64|aarch64) ARCH_FINAL="arm64"   ;;
  *)
    echo "Unsupported architecture: ${ARCH}" >&2
    exit 1
    ;;
esac

# ─── Path helpers ──────────────────────────────────────────
ASF_CONFIG_DIR=""
case "$OS_FINAL" in
  darwin) ASF_CONFIG_DIR="${HOME}/Library/Application Support/asf" ;;
  linux)  ASF_CONFIG_DIR="${XDG_CONFIG_HOME:-${HOME}/.config}/asf" ;;
esac

ASF_CACHE_DIR=""
case "$OS_FINAL" in
  darwin) ASF_CACHE_DIR="${HOME}/Library/Caches/asf" ;;
  linux)  ASF_CACHE_DIR="${XDG_CACHE_HOME:-${HOME}/.cache}/asf" ;;
esac

ASF_DATA_DIR=""
case "$OS_FINAL" in
  linux)  ASF_DATA_DIR="${HOME}/.local/share/asf" ;;
  darwin) ASF_DATA_DIR="${ASF_CONFIG_DIR}" ;;
esac

# ─── Colors ────────────────────────────────────────────────
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { printf "${CYAN}  %s${NC}\n" "$*"; }
ok()    { printf "${GREEN}  ✓ %s${NC}\n" "$*"; }
warn()  { printf "${YELLOW}  ⚠  %s${NC}\n" "$*"; }
err()   { printf "${RED}  ✗ %s${NC}\n" "$*"; exit 1; }

# ─── ASCII Art ─────────────────────────────────────────────
cat << 'EOF'
  /\   /\
 (  o.o  )
  >  ^  <
 ASF Security Framework
EOF
echo ""

# ─── Help ──────────────────────────────────────────────────
if [ "$SHOW_HELP" = true ]; then
  echo "ASF Installer"
  echo ""
  echo "Usage:"
  echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash"
  echo "  curl ... | bash -s -- --upgrade"
  echo "  curl ... | bash -s -- --repair"
  echo "  curl ... | bash -s -- --clean"
  echo "  curl ... | bash -s -- --clean --purge"
  echo ""
  echo "Options:"
  echo "  --upgrade, -u    Upgrade existing installation (backs up config)"
  echo "  --repair         Fix broken symlink/PATH without re-downloading"
  echo "  --clean          Force clean reinstall (removes binary, keeps config)"
  echo "  --purge          Only with --clean: removes config/cache/data too"
  echo "  --help, -h       Show this help"
  echo ""
  echo "Environment:"
  echo "  ASF_VERSION       Pin version (default: latest)"
  echo "  ASF_INSTALL_DIR   Custom install directory (default: ~/.local/bin)"
  echo "  GITHUB_TOKEN      Auth token for private repos"
  exit 0
fi

# ─── Auth setup ────────────────────────────────────────────
AUTH_HEADER=""
if [ -n "${GITHUB_TOKEN:-}" ]; then
  AUTH_HEADER="Authorization: token ${GITHUB_TOKEN}"
elif command -v gh &>/dev/null; then
  GH_TOKEN="$(gh auth token 2>/dev/null || echo "")"
  [ -n "$GH_TOKEN" ] && AUTH_HEADER="Authorization: token ${GH_TOKEN}"
fi

curl_get() {
  local url="$1" out="${2:-}"
  local args=(-sfL)
  [ -n "$AUTH_HEADER" ] && args+=(-H "$AUTH_HEADER")
  if [ -n "$out" ]; then
    args+=(-H "Accept: application/octet-stream" -w "%{http_code}")
    curl "${args[@]}" "$url" -o "$out" 2>/dev/null || echo "000"
  else
    curl "${args[@]}" "$url" 2>/dev/null || echo ""
  fi
}

# ─── Shell detection ───────────────────────────────────────
detect_shell() {
  CURRENT_SHELL="$(basename "${SHELL:-}" 2>/dev/null || echo "")"
  case "$CURRENT_SHELL" in
    zsh) SHELL_CONFIG="${HOME}/.zshrc" ;;
    bash) SHELL_CONFIG="${HOME}/.bashrc" ;;
    fish) SHELL_CONFIG="${HOME}/.config/fish/config.fish" ;;
    *) SHELL_CONFIG="" ;;
  esac
}

# ─── PATH setup ────────────────────────────────────────────
setup_path() {
  local dir="${INSTALL_DIR}"
  local rc_file="$1"

  # Normalize dir for grep
  local escaped_dir
  escaped_dir="$(printf '%s' "$dir" | sed 's/[][\.^$*]/\\&/g')"

  if echo ":$PATH:" | grep -q ":${dir}:"; then
    ok "${dir} is already in PATH"
    return 0
  fi

  if [ -z "$rc_file" ]; then
    warn "Could not detect shell config. Add the following to your shell config:"
    echo "  export PATH=\"\$PATH:${dir}\""
    return 1
  fi

  # Check if already in rc file
  if grep -q "export PATH=.*${escaped_dir}" "$rc_file" 2>/dev/null; then
    ok "${dir} already configured in ${rc_file}"
    return 0
  fi

  echo "" >> "$rc_file"
  echo "# ASF" >> "$rc_file"
  echo "export PATH=\"\$PATH:${dir}\"" >> "$rc_file"
  ok "Added ${dir} to PATH in ${rc_file}"
  return 0
}

# ─── Detect existing installation ─────────────────────────
detect_install() {
  EXISTING_BIN=""
  EXISTING_VER=""
  ASF_IN_PATH=""
  ASF_SYMLINK=""

  if command -v asf &>/dev/null; then
    ASF_IN_PATH="$(command -v asf)"
  fi

  for check_dir in "${ASF_HOME}" "${HOME}/.local/bin" "/usr/local/bin"; do
    if [ -f "${check_dir}/asf" ]; then
      EXISTING_BIN="${check_dir}/asf"
      EXISTING_VER=$("${check_dir}/asf" --version 2>/dev/null || echo "unknown")
      break
    fi
  done

  if [ -n "$ASF_IN_PATH" ] && [ -L "$ASF_IN_PATH" ]; then
    ASF_SYMLINK="$ASF_IN_PATH"
  fi
}

detect_install
detect_shell

# ─── Repair mode: no download, just fix symlink/PATH ──────
if [ "$REPAIR" = true ]; then
  if [ -z "$EXISTING_BIN" ]; then
    err "No ASF binary found at ${ASF_HOME}/asf. Run installer without --repair."
  fi

  echo ""
  info "Repairing ASF installation..."
  echo ""

  mkdir -p "${INSTALL_DIR}"
  rm -f "${INSTALL_DIR}/asf" 2>/dev/null || true
  ln -sf "${EXISTING_BIN}" "${INSTALL_DIR}/asf" 2>/dev/null || cp "${EXISTING_BIN}" "${INSTALL_DIR}/asf"
  chmod +x "${INSTALL_DIR}/asf" 2>/dev/null || true
  ok "Symlink created: ${INSTALL_DIR}/asf → ${EXISTING_BIN}"

  # Fix PATH
  echo ""
  info "Checking PATH..."
  setup_path "$SHELL_CONFIG"

  # Verify
  echo ""
  info "Verifying installation..."
  verify_install
  exit 0
fi

# ─── Clean mode: remove old binary before install ─────────
if [ "$CLEAN" = true ]; then
  echo ""
  info "Cleaning old ASF installation..."
  echo ""
  rm -f "${ASF_HOME}/asf" 2>/dev/null || true
  if [ -n "$ASF_SYMLINK" ]; then
    rm -f "$ASF_SYMLINK" 2>/dev/null || true
  fi
  rm -f "${HOME}/.local/bin/asf" 2>/dev/null || true
  rm -f "/usr/local/bin/asf" 2>/dev/null || true

  if [ "$PURGE" = true ]; then
    rm -rf "${ASF_CONFIG_DIR}" 2>/dev/null || true
    rm -rf "${ASF_CACHE_DIR}" 2>/dev/null || true
    rm -rf "${ASF_DATA_DIR}" 2>/dev/null || true
    rm -rf "${ASF_HOME}" 2>/dev/null || true
    ok "Config, cache, data removed"
  else
    ok "Old binaries removed (config kept)"
  fi

  EXISTING_BIN=""
  EXISTING_VER=""
  ASF_IN_PATH=""
  ASF_SYMLINK=""
  # Fall through to normal install
fi

# ─── Backup existing config on upgrade ─────────────────────
if [ "$UPGRADE" = true ] && [ -n "$EXISTING_BIN" ]; then
  if [ -f "${ASF_CONFIG_DIR}/config.yaml" ]; then
    mkdir -p "${BACKUP_DIR}"
    stamp="$(date +%Y%m%d-%H%M%S)"
    cp "${ASF_CONFIG_DIR}/config.yaml" "${BACKUP_DIR}/config.yaml.bak.${stamp}"
    ok "Config backed up to ${BACKUP_DIR}/config.yaml.bak.${stamp}"
  fi
  if [ -f "${ASF_CONFIG_DIR}/license.key" ]; then
    mkdir -p "${BACKUP_DIR}"
    stamp="$(date +%Y%m%d-%H%M%S)"
    cp "${ASF_CONFIG_DIR}/license.key" "${BACKUP_DIR}/license.key.bak.${stamp}"
    ok "License backed up to ${BACKUP_DIR}/license.key.bak.${stamp}"
  fi
fi

# ─── Existing install detection logic ─────────────────────
if [ -n "$EXISTING_BIN" ] && [ "$UPGRADE" = false ] && [ "$CLEAN" = false ]; then
  if echo "$EXISTING_VER" | grep -qi "v${VERSION}"; then
    if [ -n "$ASF_IN_PATH" ]; then
      ok "ASF v${VERSION} is already installed and available (${EXISTING_VER})"
      echo ""
      info "Run: asf"
      echo ""
      ok "Binary: ${EXISTING_BIN}"
      if [ -n "$ASF_SYMLINK" ]; then
        ok "Symlink: ${ASF_SYMLINK}"
      fi
      # Ensure PATH is configured
      setup_path "$SHELL_CONFIG"
      exit 0
    else
      warn "ASF v${VERSION} binary exists at ${EXISTING_BIN} but 'asf' is not in PATH."
      echo ""
      info "Repairing automatically..."
      # Fall through to repair
      REPAIR=true
    fi
  else
    if [ -n "$ASF_IN_PATH" ]; then
      info "Existing installation found: ${EXISTING_VER} at ${EXISTING_BIN}"
      info "Run with --upgrade to upgrade to v${VERSION}:"
      echo ""
      info "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- --upgrade"
      echo ""
      if [ -z "$ASF_IN_PATH" ]; then
        warn "'asf' command is not available in PATH — repairing before upgrade"
      else
        exit 0
      fi
    else
      warn "Old ASF binary found at ${EXISTING_BIN} but 'asf' is not callable."
      info "Repairing and upgrading automatically..."
    fi
  fi
fi

# If REPAIR was set during detection, handle it before downloading
if [ "$REPAIR" = true ] && [ "$UPGRADE" = false ]; then
  if [ -z "$EXISTING_BIN" ]; then
    err "No ASF binary found at ${ASF_HOME}/asf. Run installer without --repair."
  fi
  mkdir -p "${INSTALL_DIR}"
  rm -f "${INSTALL_DIR}/asf" 2>/dev/null || true
  ln -sf "${EXISTING_BIN}" "${INSTALL_DIR}/asf" 2>/dev/null || cp "${EXISTING_BIN}" "${INSTALL_DIR}/asf"
  chmod +x "${INSTALL_DIR}/asf" 2>/dev/null || true
  ok "Symlink created: ${INSTALL_DIR}/asf → ${EXISTING_BIN}"
  setup_path "$SHELL_CONFIG"
  echo ""
  info "Verifying installation..."
  verify_install
  exit 0
fi

# ─── Determine version ─────────────────────────────────────
if [ -z "$VERSION" ]; then
  info "Detecting latest version..."
  API_URL="https://api.github.com/repos/${REPO}/releases?per_page=10"
  VERSION="$(curl_get "$API_URL" | python3 -c "
import json,sys
try:
    releases = json.load(sys.stdin)
    for r in releases:
        if not r.get('draft', True):
            print(r['tag_name'].lstrip('v'))
            sys.exit(0)
except: pass
print('3.0.0-RC2')
" 2>/dev/null || echo "3.0.0-RC2")"
  if [ -z "$VERSION" ]; then
    VERSION="3.0.0-RC2"
    warn "Could not detect version; defaulting to ${VERSION}"
    [ -z "$AUTH_HEADER" ] && warn "For private repos, set GITHUB_TOKEN environment variable"
  fi
fi

BINARY_NAME="ASF-v${VERSION}-${OS_FINAL}-${ARCH_FINAL}"
DIRECT_DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BINARY_NAME}"
DIRECT_CHECKSUMS_URL="https://github.com/${REPO}/releases/download/v${VERSION}/checksums.txt"

# ─── Download ──────────────────────────────────────────────
TMP_DIR="$(mktemp -d)"
trap 'rm -rf "${TMP_DIR}"' EXIT

echo ""
info "Downloading ASF v${VERSION} for ${OS_FINAL}/${ARCH_FINAL}..."
echo ""

HTTP_CODE="000"
if [ -n "$AUTH_HEADER" ]; then
  HTTP_CODE="$(curl_get "$DIRECT_DOWNLOAD_URL" "${TMP_DIR}/asf")"
  info "  (authenticated)"
else
  if command -v curl &>/dev/null; then
    HTTP_CODE="$(curl -sfL -w "%{http_code}" "${DIRECT_DOWNLOAD_URL}" -o "${TMP_DIR}/asf" 2>/dev/null || echo "000")"
  elif command -v wget &>/dev/null; then
    HTTP_CODE="$(wget --server-response -q "${DIRECT_DOWNLOAD_URL}" -O "${TMP_DIR}/asf" 2>&1 \
      | grep "HTTP/" | tail -1 | awk '{print $2}' || echo "000")"
    [ -z "$HTTP_CODE" ] && HTTP_CODE="000"
  else
    err "Need curl or wget to download"
  fi
fi
info "  ${DIRECT_DOWNLOAD_URL}"

if [ ! -s "${TMP_DIR}/asf" ] || [ "$HTTP_CODE" = "000" ] || [ "$HTTP_CODE" = "404" ]; then
  rm -f "${TMP_DIR}/asf"
  echo ""
  err "Download failed (HTTP ${HTTP_CODE})."
  echo ""
  info "Possible causes:"
  info "  - No release tagged v${VERSION}"
  info "  - Binary not uploaded for ${OS_FINAL}/${ARCH_FINAL}"
  info "  - Private repo without GITHUB_TOKEN set"
  echo ""
  info "To install with a private repo, set GITHUB_TOKEN:"
  info "  export GITHUB_TOKEN=ghp_xxx"
  info "  curl ... | bash"
  echo ""
  info "Or build from source:"
  info "  git clone https://github.com/${REPO}.git"
  info "  cd asf-tui && go build -o asf-tui . && cp asf-tui ~/.asf/asf"
  exit 1
fi

chmod +x "${TMP_DIR}/asf"
ok "Download complete ($(ls -lh "${TMP_DIR}/asf" | awk '{print $5}'))"

# ─── Checksum verification ─────────────────────────────────
CHECKSUMS="$(curl -sfL "${DIRECT_CHECKSUMS_URL}" 2>/dev/null || echo "")"

if [ -n "$CHECKSUMS" ]; then
  EXPECTED_HASH="$(echo "${CHECKSUMS}" | grep "${BINARY_NAME}" | awk '{print $1}' || true)"
  if [ -n "$EXPECTED_HASH" ]; then
    COMPUTED_HASH="$(shasum -a 256 "${TMP_DIR}/asf" | awk '{print $1}')"
    if [ "$COMPUTED_HASH" != "$EXPECTED_HASH" ]; then
      err "Checksum mismatch! Expected ${EXPECTED_HASH}, got ${COMPUTED_HASH}. Download may be corrupted."
    fi
    ok "Checksum verified (SHA-256)"
  else
    warn "No checksum found for ${BINARY_NAME} — skipping verification"
  fi
else
  warn "Could not retrieve checksums.txt — skipping verification"
fi

# ─── Verify binary ─────────────────────────────────────────
BIN_VER="$("${TMP_DIR}/asf" --version 2>/dev/null || echo "unknown")"
if echo "${BIN_VER}" | grep -qi "v${VERSION}"; then
  ok "Binary verified: ${BIN_VER}"
else
  warn "Binary reports ${BIN_VER} (expected v${VERSION})"
fi

# ─── Create directories ────────────────────────────────────
mkdir -p "${ASF_HOME}"
mkdir -p "${ASF_HOME}/models"
mkdir -p "${ASF_HOME}/reports"
mkdir -p "${ASF_CONFIG_DIR}"
mkdir -p "${ASF_CACHE_DIR}"
mkdir -p "${ASF_DATA_DIR}"

# ─── Install binary ───────────────────────────────────────
cp "${TMP_DIR}/asf" "${ASF_HOME}/asf"

if [ ! -d "${INSTALL_DIR}" ] || [ ! -w "${INSTALL_DIR}" ]; then
  INSTALL_DIR="${HOME}/.local/bin"
  mkdir -p "${INSTALL_DIR}"
fi

rm -f "${INSTALL_DIR}/asf" 2>/dev/null || true
ln -sf "${ASF_HOME}/asf" "${INSTALL_DIR}/asf" 2>/dev/null || cp "${ASF_HOME}/asf" "${INSTALL_DIR}/asf"

# ─── Configure PATH ────────────────────────────────────────
echo ""
info "Configuring PATH..."
setup_path "$SHELL_CONFIG"

# ─── Default config ────────────────────────────────────────
if [ ! -f "${ASF_CONFIG_DIR}/config.yaml" ]; then
  cat > "${ASF_CONFIG_DIR}/config.yaml" << 'CONFEOF'
general:
  theme: Dark
  fox_style: Classic
analysis:
  depth: deep
  stride: true
  controls: true
ai:
  enabled: false
  active_model: ""
  installed_models: []
output:
  default: markdown
  directory: ./reports
appearance:
  theme: Dark
  fox_style: Classic
engine:
  use_native_engine: true
CONFEOF
  ok "Created default config"
fi

# ─── Verify ────────────────────────────────────────────────
verify_install() {
  echo ""
  info "Verifying installation..."
  echo ""

  ALL_OK=true

  if [ -x "${ASF_HOME}/asf" ]; then
    ok "Binary: ${ASF_HOME}/asf"
  else
    warn "Binary not found at ${ASF_HOME}/asf"
    ALL_OK=false
  fi

  if [ -x "${INSTALL_DIR}/asf" ] || [ -L "${INSTALL_DIR}/asf" ]; then
    ok "Symlink: ${INSTALL_DIR}/asf → $(readlink "${INSTALL_DIR}/asf" 2>/dev/null || echo "${ASF_HOME}/asf")"
  else
    warn "Symlink not found at ${INSTALL_DIR}/asf"
    ALL_OK=false
  fi

  if command -v asf &>/dev/null; then
    ok "Command: asf → $(command -v asf)"
  else
    warn "'asf' command is not in PATH"
    ALL_OK=false
  fi

  echo ""
  if ! command -v asf &>/dev/null; then
    warn "To use 'asf' in the current terminal, run:"
    if [ -n "$SHELL_CONFIG" ]; then
      echo "    source ${SHELL_CONFIG}"
    else
      echo "    export PATH=\"\$PATH:${INSTALL_DIR}\""
    fi
    echo ""
  fi

  if command -v asf &>/dev/null; then
    echo ""
    VER_OUT="$(asf --version 2>/dev/null || true)"
    info "asf --version: ${VER_OUT}"
    DOCTOR_OUT="$(asf doctor 2>&1 || true)"
    if echo "$DOCTOR_OUT" | grep -qi "native engine"; then
      ok "asf doctor: native engine active"
    else
      info "asf doctor: completed"
    fi
  fi

  echo ""
  if [ "$ALL_OK" = true ]; then
    ok "All checks passed."
  else
    warn "Some checks failed — see warnings above."
  fi
}

verify_install

# ─── Success ──────────────────────────────────────────────
BINARY_SIZE="$(ls -lh "${ASF_HOME}/asf" | awk '{print $5}')"
echo ""
ok "ASF v${VERSION} installed  (${BINARY_SIZE})"
echo ""
info "Binary: ${ASF_HOME}/asf"
if [ -n "$(command -v asf 2>/dev/null || true)" ]; then
  info "Command: $(command -v asf)"
fi
echo ""
info "Run: asf"
echo ""
info "Config: ${ASF_CONFIG_DIR}/config.yaml"
info "Cache:  ${ASF_CACHE_DIR}"
info "Data:   ${ASF_DATA_DIR}"
echo ""
info "Prerequisites (optional):"
info "  Tesseract (OCR):   apt install tesseract-ocr / brew install tesseract"
info "  Ollama (AI):       brew install ollama / curl -fsSL https://ollama.com/install.sh | sh"
echo ""
info "If 'asf' is not available, open a new terminal or run:"
if [ -n "$SHELL_CONFIG" ]; then
  info "  source ${SHELL_CONFIG}"
else
  info "  export PATH=\"\$PATH:${INSTALL_DIR}\""
fi
echo ""
info "Documentation: https://github.com/${REPO}"
info "Issues:        https://github.com/${REPO}/issues"
echo ""

# ─── Clean old install scripts ────────────────────────────
rm -f "${ASF_HOME}/install.sh" 2>/dev/null || true
