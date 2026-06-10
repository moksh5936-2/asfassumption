#!/bin/bash
# ASF — Architecture Security Framework
# Installer — https://github.com/moksh5936-2/asfassumption
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
#
# Environment:
#   ASF_VERSION=2.0.0       — pin a specific version
#   GITHUB_TOKEN=ghp_xxx    — for private repos (set via gh auth token)
#   ASF_INSTALL_DIR=        — custom install directory (default: /usr/local/bin)

set -euo pipefail

# ─── Config ────────────────────────────────────────────────
REPO="moksh5936-2/asfassumption"
VERSION="${ASF_VERSION:-}"
INSTALL_DIR="${ASF_INSTALL_DIR:-/usr/local/bin}"
ASF_HOME="${HOME}/.asf"
BACKUP_DIR="${ASF_HOME}/backups"

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
for arg in "$@"; do
  case "$arg" in
    --help|-h)
      echo "ASF Installer"
      echo ""
      echo "Usage:"
      echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash"
      echo "  curl ... | bash -s -- --upgrade"
      echo ""
      echo "Options:"
      echo "  --upgrade, -u    Upgrade existing installation (backs up config)"
      echo "  --help, -h       Show this help"
      echo ""
      echo "Environment:"
      echo "  ASF_VERSION       Pin version (default: latest)"
      echo "  ASF_INSTALL_DIR   Custom install directory (default: /usr/local/bin)"
      echo "  GITHUB_TOKEN      Auth token for private repos"
      exit 0
      ;;
  esac
done

UPGRADE=false
for arg in "$@"; do
  [ "$arg" = "--upgrade" ] || [ "$arg" = "-u" ] && UPGRADE=true
done

# ─── Auth setup ────────────────────────────────────────────
AUTH_HEADER=""
if [ -n "${GITHUB_TOKEN:-}" ]; then
  AUTH_HEADER="Authorization: token ${GITHUB_TOKEN}"
elif command -v gh &>/dev/null; then
  GH_TOKEN=$(gh auth token 2>/dev/null || echo "")
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

# ─── Detect existing installations ─────────────────────────
EXISTING_BIN=""
EXISTING_VER=""
for check_dir in "${ASF_HOME}" "${INSTALL_DIR}" "${HOME}/.local/bin"; do
  if [ -f "${check_dir}/asf" ]; then
    EXISTING_BIN="${check_dir}/asf"
    EXISTING_VER=$("${check_dir}/asf" --version 2>/dev/null || echo "unknown")
    break
  fi
done

# ─── Backup existing config on upgrade ─────────────────────
if [ "$UPGRADE" = true ] && [ -n "$EXISTING_BIN" ]; then
  if [ -f "${ASF_CONFIG_DIR}/config.yaml" ]; then
    mkdir -p "${BACKUP_DIR}"
    local stamp
    stamp=$(date +%Y%m%d-%H%M%S)
    cp "${ASF_CONFIG_DIR}/config.yaml" "${BACKUP_DIR}/config.yaml.bak.${stamp}"
    ok "Config backed up to ${BACKUP_DIR}/config.yaml.bak.${stamp}"
  fi
  if [ -f "${ASF_CONFIG_DIR}/license.key" ]; then
    mkdir -p "${BACKUP_DIR}"
    stamp=$(date +%Y%m%d-%H%M%S)
    cp "${ASF_CONFIG_DIR}/license.key" "${BACKUP_DIR}/license.key.bak.${stamp}"
    ok "License backed up to ${BACKUP_DIR}/license.key.bak.${stamp}"
  fi
fi

# ─── Determine version ─────────────────────────────────────
if [ -z "$VERSION" ]; then
  info "Detecting latest version..."
  API_URL="https://api.github.com/repos/${REPO}/releases/latest"
  VERSION=$(curl_get "$API_URL" | grep '"tag_name":' | sed 's/.*"tag_name": "v\(.*\)",.*/\1/' || echo "")
  if [ -z "$VERSION" ]; then
    VERSION="2.0.0"
    warn "Could not detect latest version; defaulting to ${VERSION}"
    [ -z "$AUTH_HEADER" ] && warn "For private repos, set GITHUB_TOKEN environment variable"
  fi
fi

BINARY_NAME="ASF-v${VERSION}-${OS_FINAL}-${ARCH_FINAL}"
DIRECT_DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BINARY_NAME}"
DIRECT_CHECKSUMS_URL="https://github.com/${REPO}/releases/download/v${VERSION}/checksums.txt"

# ─── Upgrade check ─────────────────────────────────────────
if [ -n "$EXISTING_BIN" ] && [ "$UPGRADE" = false ]; then
  if echo "$EXISTING_VER" | grep -qi "v${VERSION}"; then
    ok "ASF v${VERSION} is already installed (${EXISTING_VER})"
    echo ""
    info "Run: asf"
    echo ""
    info "To force reinstall: curl ... | bash -s -- --upgrade"
    exit 0
  fi
  info "Existing installation found (${EXISTING_VER})"
  info "Use --upgrade to upgrade to v${VERSION}"
  exit 0
fi

# ─── Get asset ID for API download (if auth available) ─────
ASSET_ID=""
RELEASE_API_URL="https://api.github.com/repos/${REPO}/releases/tags/v${VERSION}"
if [ -n "$AUTH_HEADER" ]; then
  ASSET_ID=$(curl_get "$RELEASE_API_URL" \
    | python3 -c "
import json,sys
try:
    data = json.load(sys.stdin)
    for a in data.get('assets', []):
        if a['name'] == '${BINARY_NAME}':
            print(a['id'])
except: pass
" 2>/dev/null || echo "")
fi

# ─── Download ──────────────────────────────────────────────
TMP_DIR=$(mktemp -d)
trap 'rm -rf "${TMP_DIR}"' EXIT

echo ""
info "Downloading ASF v${VERSION} for ${OS_FINAL}/${ARCH_FINAL}..."
echo ""

HTTP_CODE="000"
if [ -n "$ASSET_ID" ]; then
  API_ASSET_URL="https://api.github.com/repos/${REPO}/releases/assets/${ASSET_ID}"
  info "  (authenticated API download)"
  HTTP_CODE=$(curl_get "$API_ASSET_URL" "${TMP_DIR}/asf")
elif [ -n "$AUTH_HEADER" ]; then
  HTTP_CODE=$(curl_get "$DIRECT_DOWNLOAD_URL" "${TMP_DIR}/asf")
  info "  ${DIRECT_DOWNLOAD_URL}"
else
  info "  ${DIRECT_DOWNLOAD_URL}"
  if command -v curl &>/dev/null; then
    HTTP_CODE=$(curl -sfL -w "%{http_code}" "${DIRECT_DOWNLOAD_URL}" -o "${TMP_DIR}/asf" 2>/dev/null || echo "000")
  elif command -v wget &>/dev/null; then
    HTTP_CODE=$(wget --server-response -q "${DIRECT_DOWNLOAD_URL}" -O "${TMP_DIR}/asf" 2>&1 \
      | grep "HTTP/" | tail -1 | awk '{print $2}' || echo "000")
    [ -z "$HTTP_CODE" ] && HTTP_CODE="000"
  else
    err "Need curl or wget to download"
  fi
fi

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
CHECKSUMS=""
if [ -n "$AUTH_HEADER" ]; then
  CHECKSUMS_ASSET_ID=$(curl_get "$RELEASE_API_URL" \
    | python3 -c "
import json,sys
try:
    data = json.load(sys.stdin)
    for a in data.get('assets', []):
        if a['name'] == 'checksums.txt':
            print(a['id'])
except: pass
" 2>/dev/null || echo "")
  if [ -n "$CHECKSUMS_ASSET_ID" ]; then
    CHECKSUMS_API_URL="https://api.github.com/repos/${REPO}/releases/assets/${CHECKSUMS_ASSET_ID}"
    curl_get "$CHECKSUMS_API_URL" "${TMP_DIR}/checksums.txt" >/dev/null 2>&1 || true
    CHECKSUMS=$(cat "${TMP_DIR}/checksums.txt" 2>/dev/null || echo "")
  fi
else
  CHECKSUMS=$(curl -sfL "${DIRECT_CHECKSUMS_URL}" 2>/dev/null || echo "")
fi

if [ -n "$CHECKSUMS" ]; then
  EXPECTED_HASH=$(echo "${CHECKSUMS}" | grep "${BINARY_NAME}" | awk '{print $1}' || true)
  if [ -n "$EXPECTED_HASH" ]; then
    COMPUTED_HASH=$(shasum -a 256 "${TMP_DIR}/asf" | awk '{print $1}')
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
BIN_VER=$("${TMP_DIR}/asf" --version 2>/dev/null || echo "unknown")
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

# Remove stale symlinks before creating new one
rm -f "${INSTALL_DIR}/asf" 2>/dev/null || true
ln -sf "${ASF_HOME}/asf" "${INSTALL_DIR}/asf" 2>/dev/null || cp "${ASF_HOME}/asf" "${INSTALL_DIR}/asf"

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

# ─── PATH warning ──────────────────────────────────────────
case ":$PATH:" in
  *:"${INSTALL_DIR}":*) ;;
  *)
    echo ""
    warn "Add ${INSTALL_DIR} to your PATH:"
    echo "      export PATH=\"\$PATH:${INSTALL_DIR}\""
    echo ""
    info "Or add it to your shell config:"
    echo "      echo 'export PATH=\"\$PATH:${INSTALL_DIR}\"' >> ~/.zshrc"
    ;;
esac

# ─── Download Python Engine ───────────────────────────────
ENGINE_URL="https://github.com/${REPO}/releases/download/v${VERSION}/asf-python-engine-v${VERSION}.tar.gz"
ENGINE_TAR="${TMP_DIR}/asf-python-engine.tar.gz"

echo ""
info "Downloading ASF Python Engine v${VERSION}..."
info "  ${ENGINE_URL}"

ENGINE_HTTP_CODE="000"
if command -v curl &>/dev/null; then
  ENGINE_HTTP_CODE=$(curl -sfL -w "%{http_code}" "${ENGINE_URL}" -o "${ENGINE_TAR}" 2>/dev/null || echo "000")
elif command -v wget &>/dev/null; then
  ENGINE_HTTP_CODE=$(wget --server-response -q "${ENGINE_URL}" -O "${ENGINE_TAR}" 2>&1 \
    | grep "HTTP/" | tail -1 | awk '{print $2}' || echo "000")
fi

if [ -s "${ENGINE_TAR}" ]; then
  ok "Engine download complete ($(ls -lh "${ENGINE_TAR}" | awk '{print $5}'))"

  # Extract engine to data dir
  mkdir -p "${ASF_DATA_DIR}"
  rm -rf "${ASF_DATA_DIR}/engine" 2>/dev/null || true
  tar -xzf "${ENGINE_TAR}" -C "${ASF_DATA_DIR}" 2>/dev/null || {
    warn "Failed to extract Python engine"
  }
  if [ -d "${ASF_DATA_DIR}/asf" ]; then
    ok "Python engine extracted to ${ASF_DATA_DIR}"
  else
    warn "Python engine extraction may have failed — asf/ directory not found"
  fi
else
  warn "Python engine download failed (HTTP ${ENGINE_HTTP_CODE})"
  info "  The Go native engine handles analysis — Python is optional."
  info "  Run 'asf doctor --fix' after installation if you still need Python."
fi

# ─── Run post-install verification ─────────────────────────
echo ""
info "Running post-install verification..."

if ! "${INSTALL_DIR}/asf" --version &>/dev/null; then
  warn "Binary not working from install location"
fi

if "${INSTALL_DIR}/asf" doctor &>/dev/null; then
  ok "asf doctor: passed"
else
  warn "asf doctor had warnings (see above)"
fi

# ─── Success ──────────────────────────────────────────────
BINARY_SIZE=$(ls -lh "${ASF_HOME}/asf" | awk '{print $5}')
echo ""
ok "ASF v${VERSION} installed  (${BINARY_SIZE})"
echo ""
info "Run: asf"
echo ""
info "Config: ${ASF_CONFIG_DIR}/config.yaml"
info "Cache:  ${ASF_CACHE_DIR}"
info "Data:   ${ASF_DATA_DIR}"
echo ""
info "Prerequisites (full functionality):"
info "  Tesseract (OCR):   apt install tesseract-ocr / brew install tesseract"
info "  Ollama (AI):       brew install ollama / curl -fsSL https://ollama.com/install.sh | sh"
echo ""
info "Documentation: https://github.com/${REPO}"
info "Issues:        https://github.com/${REPO}/issues"
echo ""

# ─── Clean old install scripts ────────────────────────────
rm -f "${ASF_HOME}/install.sh" 2>/dev/null || true
