#!/bin/bash
set -euo pipefail

ASF_VERSION="1.0.0"
ASF_REPO="moksh5936-2/asfassumption"

echo "  /\   /\ "
echo " (  o.o  )"
echo "  >  ^  < "
echo " ASF Security Framework"
echo ""

# ─── Detect platform ───────────────────────────────────
OS="$(uname -s)"
ARCH="$(uname -m)"

case "$OS" in
    Darwin)
        case "$ARCH" in
            arm64) BINARY="asf-darwin-arm64" ;;
            x86_64) BINARY="asf-darwin-amd64" ;;
            *) echo "Unsupported arch: $ARCH"; exit 1 ;;
        esac
        ;;
    Linux)
        case "$ARCH" in
            aarch64|arm64) BINARY="asf-linux-arm64" ;;
            x86_64) BINARY="asf-linux-amd64" ;;
            *) echo "Unsupported arch: $ARCH"; exit 1 ;;
        esac
        ;;
    *)
        echo "Unsupported OS: $OS"
        echo "Windows users: download asf-windows-amd64.exe from https://github.com/$ASF_REPO/releases"
        exit 1
        ;;
esac

echo "Installing ASF v${ASF_VERSION} for ${OS}/${ARCH}..."

mkdir -p "$HOME/.asf"

# ─── Try local binary first, fall back to download ────
SCRIPT_DIR="$(cd "$(dirname "$0")" 2>/dev/null && pwd || echo "")"
LOCAL_BINARY=""
if [ -n "$SCRIPT_DIR" ]; then
    LOCAL_BINARY="${SCRIPT_DIR}/${BINARY}"
    [ ! -f "$LOCAL_BINARY" ] && LOCAL_BINARY="${SCRIPT_DIR}/release/${BINARY}"
    [ ! -f "$LOCAL_BINARY" ] && LOCAL_BINARY="${SCRIPT_DIR}/../release/${BINARY}"
    [ ! -f "$LOCAL_BINARY" ] && LOCAL_BINARY=""
fi

if [ -n "$LOCAL_BINARY" ] && [ -s "$LOCAL_BINARY" ]; then
    cp "$LOCAL_BINARY" "$HOME/.asf/asf"
    echo "  ✓ Using local binary: ${LOCAL_BINARY}"
else
    DOWNLOAD_URL="https://github.com/${ASF_REPO}/releases/download/v${ASF_VERSION}/${BINARY}"
    echo "  Downloading from: $DOWNLOAD_URL"
    echo ""

    if command -v curl &>/dev/null; then
        HTTP_CODE=$(curl -sfL -w "%{http_code}" "$DOWNLOAD_URL" -o "$HOME/.asf/asf" 2>/dev/null || echo "000")
    elif command -v wget &>/dev/null; then
        HTTP_CODE=$(wget --server-response -q "$DOWNLOAD_URL" -O "$HOME/.asf/asf" 2>&1 | grep "HTTP/" | tail -1 | awk '{print $2}' || echo "000")
        [ -z "$HTTP_CODE" ] && HTTP_CODE="000"
    else
        echo "Error: need curl or wget"
        exit 1
    fi

    if [ ! -s "$HOME/.asf/asf" ] || [ "$HTTP_CODE" = "000" ] || [ "$HTTP_CODE" = "404" ]; then
        rm -f "$HOME/.asf/asf"
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

chmod +x "$HOME/.asf/asf"

# ─── Create directories ────────────────────────────────
mkdir -p "$HOME/.asf/models"
mkdir -p "$HOME/.asf/reports"

# ─── Default config ────────────────────────────────────
if [ ! -f "$HOME/.asf/config.yaml" ]; then
    cat > "$HOME/.asf/config.yaml" << 'CONFEOF'
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
CONFEOF
    echo "  ✓ Created default config"
fi

# ─── Symlink to PATH ──────────────────────────────────
INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

ln -sf "$HOME/.asf/asf" "$INSTALL_DIR/asf" 2>/dev/null || cp "$HOME/.asf/asf" "$INSTALL_DIR/asf"

# ─── Warn if PATH missing ──────────────────────────────
case ":$PATH:" in
    *:"$INSTALL_DIR":*) ;;
    *)
        echo ""
        echo "  ⚠  Add $INSTALL_DIR to your PATH:"
        echo "      export PATH=\"\$PATH:$INSTALL_DIR\""
        echo ""
        echo "     Or add it to your shell config (~/.zshrc, ~/.bashrc):"
        echo "      echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.zshrc"
        ;;
esac

echo ""
BINARY_SIZE=$(ls -lh "$HOME/.asf/asf" | awk '{print $5}')
echo " ✓ ASF v${ASF_VERSION} installed  (${BINARY_SIZE})"
echo ""
echo "   Run: asf"
echo ""
echo "   Prerequisites (full functionality):"
echo "     - Python ASF engine: cd /path/to/asf && pip install -e ."
echo "     - Ollama (AI): brew install ollama"
echo "     - Tesseract (OCR): brew install tesseract"
echo ""
echo "   License (enterprise):"
echo "     echo 'ASF-XXXX-XXXX-XXXX-XXXX' > ~/.asf/license.key"
