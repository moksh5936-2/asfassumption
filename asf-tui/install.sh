#!/bin/bash
# ASF — Architecture Security Framework
# Local/Dev Installer — https://github.com/moksh5936-2/asfassumption
#
# Tries local binary first, falls back to downloading from GitHub.
# For full installer with upgrade/repair/clean, use the root install.sh.
#
# Usage:
#   ./install.sh                    — install from local release/ directory
#   curl ... | bash                 — download and install latest
#   curl ... | bash -s -- --upgrade — upgrade existing installation

set -euo pipefail

ASF_VERSION="5.0.2"
ASF_REPO="moksh5936-2/asfassumption"
LATEST_VERSION="v${ASF_VERSION}"
ASSET_VERSION="${ASF_VERSION}"

# ─── Parse flags ──────────────────────────────────────────
UPGRADE=false
for arg in "$@"; do
  [ "$arg" = "--upgrade" ] || [ "$arg" = "-u" ] && UPGRADE=true
done

# ─── ASCII Art ─────────────────────────────────────────────
cat << 'EOF'
  /\   /\
 (  o.o  )
  >  ^  <
 ASF Security Framework
EOF
echo ""

# ─── Platform detection ───────────────────────────────────
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
  Darwin)
    case "$ARCH" in
      arm64) BINARY="ASF-${LATEST_VERSION}-darwin-arm64" ;;
      x86_64) BINARY="ASF-${LATEST_VERSION}-darwin-amd64" ;;
      *) echo "Unsupported arch: $ARCH"; exit 1 ;;
    esac
    ;;
  Linux)
    case "$ARCH" in
      aarch64|arm64) BINARY="ASF-${LATEST_VERSION}-linux-arm64" ;;
      x86_64) BINARY="ASF-${LATEST_VERSION}-linux-amd64" ;;
      *) echo "Unsupported arch: $ARCH"; exit 1 ;;
    esac
    ;;
  *)
    echo "Unsupported OS: $OS"
    echo "Windows users: download asf-windows-amd64.exe from https://github.com/$ASF_REPO/releases"
    exit 1
    ;;
esac

ASF_HOME="${HOME}/.asf"
INSTALL_DIR="${HOME}/.local/bin"

echo "Installing ASF v${ASF_VERSION} for ${OS}/${ARCH}..."
echo ""

# ─── Detect existing ─────────────────────────────────────
EXISTING_BIN=""
EXISTING_VER=""
ASF_IN_PATH=""
if command -v asf &>/dev/null; then
  ASF_IN_PATH="$(command -v asf)"
fi
if [ -f "${ASF_HOME}/asf" ]; then
  EXISTING_BIN="${ASF_HOME}/asf"
  EXISTING_VER="$("${ASF_HOME}/asf" --version 2>/dev/null || echo "unknown")"
fi

# ─── Upgrade or repair logic ──────────────────────────────
NEED_DOWNLOAD=true
if [ -f "${ASF_HOME}/asf" ]; then
  if echo "$EXISTING_VER" | grep -qi "v${ASF_VERSION}"; then
    if [ -n "$ASF_IN_PATH" ]; then
      echo "  ✓ ASF v${ASF_VERSION} already installed and available."
      echo ""
      echo "  Run: asf"
      exit 0
    else
      echo "  ⚠  Binary exists but 'asf' not in PATH — repairing symlink..."
      NEED_DOWNLOAD=false
    fi
  elif [ "$UPGRADE" = false ]; then
    echo "  Existing: ${EXISTING_VER}"
    echo "  Latest:   v${ASF_VERSION}"
    echo ""
    echo "  Run with --upgrade to upgrade:"
    echo "    curl -fsSL https://raw.githubusercontent.com/${ASF_REPO}/main/install.sh | bash -s -- --upgrade"
    echo ""
    if [ -z "$ASF_IN_PATH" ]; then
      echo "  ⚠  'asf' command not in PATH — repairing symlink to existing binary."
      NEED_DOWNLOAD=false
    else
      exit 0
    fi
  fi
fi

# ─── Download or use local ────────────────────────────────
if [ "$NEED_DOWNLOAD" = true ]; then
  mkdir -p "${ASF_HOME}"

  SCRIPT_DIR="$(cd "$(dirname "$0")" 2>/dev/null && pwd || echo "")"
  LOCAL_BINARY=""
  if [ -n "$SCRIPT_DIR" ]; then
    LOCAL_BINARY="${SCRIPT_DIR}/${BINARY}"
    [ ! -f "$LOCAL_BINARY" ] && LOCAL_BINARY="${SCRIPT_DIR}/release/${BINARY}"
    [ ! -f "$LOCAL_BINARY" ] && LOCAL_BINARY="${SCRIPT_DIR}/../release/${BINARY}"
    [ ! -f "$LOCAL_BINARY" ] && LOCAL_BINARY=""
  fi

  if [ -n "$LOCAL_BINARY" ] && [ -s "$LOCAL_BINARY" ]; then
    cp "$LOCAL_BINARY" "${ASF_HOME}/asf"
    echo "  ✓ Using local binary: ${LOCAL_BINARY}"
  else
    DOWNLOAD_URL="https://github.com/${ASF_REPO}/releases/download/${LATEST_VERSION}/${BINARY}"
    echo "  Downloading from: ${DOWNLOAD_URL}"
    echo ""

    # Validate URL before download
    if printf '%s' "$DOWNLOAD_URL" | grep -q '[[:space:]]'; then
      echo "  ✗ Download URL contains whitespace or newlines: ${DOWNLOAD_URL}"
      exit 1
    fi
    if ! printf '%s' "$BINARY" | grep -q "^ASF-"; then
      echo "  ✗ Invalid binary name: ${BINARY}"
      exit 1
    fi

    HTTP_CODE="000"
    if command -v curl &>/dev/null; then
      HTTP_CODE="$(curl -sfL -w "%{http_code}" "${DOWNLOAD_URL}" -o "${ASF_HOME}/asf" 2>/dev/null || echo "000")"
    elif command -v wget &>/dev/null; then
      HTTP_CODE="$(wget --server-response -q "${DOWNLOAD_URL}" -O "${ASF_HOME}/asf" 2>&1 \
        | grep "HTTP/" | tail -1 | awk '{print $2}' || echo "000")"
      [ -z "$HTTP_CODE" ] && HTTP_CODE="000"
    else
      echo "Error: need curl or wget"
      exit 1
    fi

    if [ ! -s "${ASF_HOME}/asf" ] || [ "$HTTP_CODE" = "000" ] || [ "$HTTP_CODE" = "404" ]; then
      rm -f "${ASF_HOME}/asf"
      echo ""
      echo "  ✗ Download failed (HTTP ${HTTP_CODE})."
      echo ""
      echo "    No release binary found for v${ASF_VERSION}."
      echo "    Possible causes:"
      echo "      - No GitHub release tagged v${ASF_VERSION}"
      echo "      - Binary not uploaded for ${OS}/${ARCH}"
      echo ""
      echo "    To build from source instead:"
      echo "      git clone https://github.com/${ASF_REPO}.git"
      echo "      cd asf-tui && go build -o asf-tui . && cp asf-tui ~/.asf/asf"
      echo ""
      echo "    If you already built the binary, copy it to ~/.asf/asf manually."
      echo ""
      exit 1
    fi
    echo "  ✓ Download complete"
  fi

  chmod +x "${ASF_HOME}/asf"
fi

# ─── Create directories ────────────────────────────────────
mkdir -p "${ASF_HOME}/models"
mkdir -p "${ASF_HOME}/reports"

# ─── Default config ────────────────────────────────────────
if [ ! -f "${ASF_HOME}/config.yaml" ]; then
  cat > "${ASF_HOME}/config.yaml" << 'CONFEOF'
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
  echo "  ✓ Created default config"
fi

# ─── Symlink to PATH ──────────────────────────────────────
mkdir -p "${INSTALL_DIR}"
rm -f "${INSTALL_DIR}/asf" 2>/dev/null || true
ln -sf "${ASF_HOME}/asf" "${INSTALL_DIR}/asf" 2>/dev/null || cp "${ASF_HOME}/asf" "${INSTALL_DIR}/asf"

# ─── Verify ────────────────────────────────────────────────
echo ""
echo "  Verifying installation..."
echo ""

if command -v asf &>/dev/null; then
  echo "  ✓ Command: asf → $(command -v asf)"
  VER_OUT="$(asf --version 2>/dev/null || true)"
  echo "  ✓ ${VER_OUT}"
else
  echo "  ⚠  'asf' not in PATH."
  echo ""
  echo "  Add ${INSTALL_DIR} to your PATH:"
  echo ""
  echo "    For zsh:"
  echo "      echo 'export PATH=\"\$PATH:${INSTALL_DIR}\"' >> ~/.zshrc"
  echo "      source ~/.zshrc"
  echo ""
  echo "    For bash:"
  echo "      echo 'export PATH=\"\$PATH:${INSTALL_DIR}\"' >> ~/.bashrc"
  echo "      source ~/.bashrc"
  echo ""
fi

# ─── Success ──────────────────────────────────────────────
BINARY_SIZE="$(ls -lh "${ASF_HOME}/asf" | awk '{print $5}')"
echo ""
echo " ✓ ASF v${ASF_VERSION} installed  (${BINARY_SIZE})"
echo ""
echo "   Run: asf"
echo ""
echo "   Prerequisites (optional):"
echo "     - Ollama (AI): brew install ollama"
echo "     - Tesseract (OCR): brew install tesseract"
echo ""
echo "   License (enterprise):"
echo "     echo 'ASF-XXXX-XXXX-XXXX-XXXX' > ~/.asf/license.key"
