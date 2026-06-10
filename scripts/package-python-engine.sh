#!/bin/bash
# Package the Python ASF engine into a portable tarball for release.
# Usage: ./scripts/package-python-engine.sh <version>
# Output: asf-python-engine-v<VERSION>.tar.gz

set -euo pipefail

VERSION="${1:-}"
if [ -z "$VERSION" ]; then
  echo "Usage: $0 <version>" >&2
  exit 1
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
OUTPUT_DIR="${REPO_ROOT}/release"
OUTPUT_FILE="${OUTPUT_DIR}/asf-python-engine-v${VERSION}.tar.gz"

mkdir -p "$OUTPUT_DIR"

# Create a staging directory with only what's needed at runtime
STAGING_DIR=$(mktemp -d)
trap 'rm -rf "${STAGING_DIR}"' EXIT

# Copy the Python package, pyproject.toml, and setup.py
cp -r "${REPO_ROOT}/asf" "${STAGING_DIR}/asf"
cp "${REPO_ROOT}/pyproject.toml" "${STAGING_DIR}/"
cp "${REPO_ROOT}/setup.py" "${STAGING_DIR}/"

# Remove any __pycache__, .pyc, .pyo, .egg-info
find "${STAGING_DIR}" \( -name "__pycache__" -o -name "*.pyc" -o -name "*.pyo" -o -name "*.egg-info" \) -exec rm -rf {} + 2>/dev/null || true

# Create the tarball
cd "${STAGING_DIR}"
tar -czf "$OUTPUT_FILE" .

SIZE=$(ls -lh "$OUTPUT_FILE" | awk '{print $5}')
echo "Packaged: ${OUTPUT_FILE} (${SIZE})"
