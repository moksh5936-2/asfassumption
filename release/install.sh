#!/bin/bash
set -euo pipefail

ASF_VERSION="1.0.0"
ASF_REPO="moksh5936-2/asfassumption"

echo "  /\   /\ "
echo " (  o.o  )"
echo "  >  ^  < "
echo " ASF Security Framework"
echo ""

if [ -f "$HOME/.asf/config.yaml" ]; then
    echo "ASF appears to be installed already."
    echo "Run 'asf' to start."
    echo ""
    echo "To reinstall, remove ~/.asf first."
    exit 0
fi

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
        echo "Windows users: download from https://github.com/$ASF_REPO/releases"
        exit 1
        ;;
esac

echo "Installing ASF v${ASF_VERSION} for ${OS}/${ARCH}..."
echo ""

mkdir -p "$HOME/.asf"

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
    echo "  ✗ Download failed (HTTP $HTTP_CODE or empty file)."
    echo ""
    echo "    This usually means the release binary doesn't exist yet."
    echo "    Possible causes:"
    echo "      - No GitHub release tagged v${ASF_VERSION}"
    echo "      - Binary not uploaded for your platform (${OS}/${ARCH})"
    echo "      - Release URL is incorrect"
    echo ""
    echo "    To build from source instead:"
    echo "      git clone https://github.com/${ASF_REPO}.git"
    echo "      cd asf-tui && go build -o asf-tui ."
    echo ""
    exit 1
fi

chmod +x "$HOME/.asf/asf"
echo "  ✓ Download complete ($(ls -lh "$HOME/.asf/asf" | awk '{print $5}'))"

mkdir -p "$HOME/.asf/models"
mkdir -p "$HOME/.asf/reports"

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
    echo "Created default config."
fi

INSTALL_DIR="/usr/local/bin"
if [ ! -w "$INSTALL_DIR" ]; then
    INSTALL_DIR="$HOME/.local/bin"
    mkdir -p "$INSTALL_DIR"
fi

ln -sf "$HOME/.asf/asf" "$INSTALL_DIR/asf" 2>/dev/null || cp "$HOME/.asf/asf" "$INSTALL_DIR/asf"

echo ""
echo " ✓ ASF v${ASF_VERSION} installed"
echo ""
echo "   Run: asf"
echo ""
echo "   To set up a license:"
echo "     echo 'ASF-XXXX-XXXX-XXXX-XXXX' > ~/.asf/license.key"
echo ""
echo "   Requirements for full functionality:"
echo "     - Ollama (for AI features): brew install ollama"
echo "     - Tesseract (for image OCR): brew install tesseract"
