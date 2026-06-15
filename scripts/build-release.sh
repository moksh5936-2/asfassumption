#!/usr/bin/env bash
# ASF Release Build Script
# Usage: ./scripts/build-release.sh [version]
# Requires: Go 1.24+
# Output: release/ASF-v{VERSION}-{OS}-{ARCH}

set -euo pipefail

VERSION="${1:-5.1.0}"
ROOT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
RELEASE_DIR="${ROOT_DIR}/release"
BUILD_DIR="${ROOT_DIR}/asf-tui"

echo "=== ASF Release Build v${VERSION} ==="

# Check Go
if ! command -v go &> /dev/null; then
  echo "ERROR: Go is not installed. Install Go 1.24+ to build."
  exit 1
fi

echo "Go: $(go version)"

# Clean old binaries
echo ""
echo "Cleaning release directory..."
rm -f "${RELEASE_DIR}"/ASF-v*
rm -f "${RELEASE_DIR}"/asf-*
rm -f "${RELEASE_DIR}"/*.tar.gz
rm -f "${RELEASE_DIR}"/*.zip
rm -f "${RELEASE_DIR}"/*.exe

# Build
echo ""
echo "Building for all platforms..."

build_platform() {
  local GOOS="$1" GOARCH="$2"
  local OUT

  if [ "$GOOS" = "windows" ]; then
    OUT="${RELEASE_DIR}/ASF-v${VERSION}-${GOOS}-${GOARCH}.exe"
  else
    OUT="${RELEASE_DIR}/ASF-v${VERSION}-${GOOS}-${GOARCH}"
  fi

  echo "  Building ${GOOS}/${GOARCH}..."
  cd "${BUILD_DIR}"
  GOOS="$GOOS" GOARCH="$GOARCH" CGO_ENABLED=0 go build -ldflags="-s -w" -o "$OUT" .
  echo "  -> $(ls -lh "$OUT" | awk '{print $5}')"
}

build_platform linux   amd64
build_platform linux   arm64
build_platform darwin  amd64
build_platform darwin  arm64
build_platform windows amd64

# Verify
echo ""
echo "Verifying binaries..."
for f in "${RELEASE_DIR}"/ASF-v*; do
  if file "$f" | grep -q "executable"; then
    echo "  ✓ $(basename "$f")"
  else
    echo "  ✗ $(basename "$f") — not an executable"
  fi
done

# Generate checksums
echo ""
echo "Generating checksums.txt..."
cd "${RELEASE_DIR}"
shasum -a 256 ASF-v* > checksums.txt 2>/dev/null || sha256sum ASF-v* > checksums.txt
echo "  checksums.txt updated"

# Update VERSION
echo "${VERSION}" > "${RELEASE_DIR}/VERSION"

# Summary
echo ""
echo "=== Build Complete ==="
echo "Version: v${VERSION}"
echo ""
echo "Artifacts:"
ls -lh "${RELEASE_DIR}" | grep -E "(ASF-v|checksums|VERSION)" | awk '{print "  " $NF " (" $5 ")"}'
echo ""
echo "Next steps:"
echo "  git tag v${VERSION} && git push origin v${VERSION}"
echo "  # GitHub Actions will build and publish automatically"
echo "  # Or upload release/ files manually to GitHub Releases"
echo ""
