#!/usr/bin/env bash
# ASF Release Build Script
# Usage: ./scripts/build-release.sh [version]
# Requires: Go 1.24+
# Output: release/asf-{platform} archives

set -euo pipefail

VERSION="${1:-1.0.0}"
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
RELEASE_DIR="${ROOT_DIR}/release"
BUILD_DIR="${ROOT_DIR}/asf-tui"
PLATFORMS=(
  "linux/amd64"
  "linux/arm64"
  "darwin/amd64"
  "darwin/arm64"
  "windows/amd64"
)

echo "=== ASF Release Build v${VERSION} ==="

# Check Go
if ! command -v go &> /dev/null; then
  echo "ERROR: Go is not installed. Install Go 1.24+ to build."
  exit 1
fi

GO_VERSION=$(go version | grep -oP 'go\K[0-9]+\.[0-9]+')
echo "Go version: ${GO_VERSION}"

if (( $(echo "$GO_VERSION < 1.24" | bc -l) )); then
  echo "ERROR: Go 1.24+ required, found $GO_VERSION"
  exit 1
fi

# Clean
echo ""
echo "Cleaning release directory..."
rm -f "${RELEASE_DIR}"/asf-*
rm -f "${RELEASE_DIR}"/*.tar.gz
rm -f "${RELEASE_DIR}"/*.zip
rm -f "${RELEASE_DIR}"/*.exe

# Build
echo ""
echo "Building for all platforms..."
cd "${BUILD_DIR}"

for platform in "${PLATFORMS[@]}"; do
  GOOS="${platform%%/*}"
  GOARCH="${platform##*/}"

  if [ "$GOOS" = "windows" ]; then
    OUT="${RELEASE_DIR}/asf-${GOOS}-${GOARCH}.exe"
  else
    OUT="${RELEASE_DIR}/asf-${GOOS}-${GOARCH}"
  fi

  echo "  Building ${GOOS}/${GOARCH}..."
  GOOS="$GOOS" GOARCH="$GOARCH" go build -ldflags="-X 'main.version=${VERSION}'" -o "$OUT" .

  # Verify
  file "$OUT"
done

# Generate checksums
echo ""
echo "Generating checksums..."
cd "${RELEASE_DIR}"
shasum -a 256 asf-* install.sh VERSION 2>/dev/null > checksums.txt || \
sha256sum asf-* install.sh VERSION 2>/dev/null > checksums.txt
echo "  checksums.txt updated"

# Create archives
echo ""
echo "Creating archives..."
for platform in "${PLATFORMS[@]}"; do
  GOOS="${platform%%/*}"
  GOARCH="${platform##*/}"

  if [ "$GOOS" = "windows" ]; then
    BINARY="asf-${GOOS}-${GOARCH}.exe"
    ARCHIVE="${GOOS}-${GOARCH}.zip"
    if [ -f "${RELEASE_DIR}/${BINARY}" ]; then
      cd "${RELEASE_DIR}"
      zip -j "${ARCHIVE}" "${BINARY}" install.sh VERSION checksums.txt
      echo "  ${ARCHIVE} created"
    fi
  else
    BINARY="asf-${GOOS}-${GOARCH}"
    ARCHIVE="${GOOS}-${GOARCH}.tar.gz"
    if [ -f "${RELEASE_DIR}/${BINARY}" ]; then
      cd "${RELEASE_DIR}"
      tar czf "${ARCHIVE}" "${BINARY}" install.sh VERSION checksums.txt
      echo "  ${ARCHIVE} created"
    fi
  fi
done

# Update VERSION
echo "${VERSION}" > "${RELEASE_DIR}/VERSION"

# Summary
echo ""
echo "=== Build Complete ==="
echo "Version: v${VERSION}"
echo "Release directory: ${RELEASE_DIR}"
echo ""
echo "Artifacts:"
ls -lh "${RELEASE_DIR}" 2>/dev/null | grep -v "\.git"

echo ""
echo "Next steps:"
echo "  1. Test each binary: ./release/asf-darwin-arm64 --version"
echo "  2. Verify checksums: cd release && shasum -a 256 -c checksums.txt"
echo "  3. Create GitHub release"
echo "  4. Upload archives"
echo ""
