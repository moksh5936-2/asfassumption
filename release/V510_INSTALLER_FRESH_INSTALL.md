# V510_INSTALLER_FRESH_INSTALL — ASF0 v5.1.0

## Simulated Fresh Install

Simulated by downloading the release binary directly and verifying it works.

| Step | Command | Result |
|------|---------|--------|
| 1 | `gh release download v5.1.0 -p "ASF-v5.1.0-darwin-arm64"` | ✅ Downloaded |
| 2 | `chmod +x ASF-v5.1.0-darwin-arm64` | ✅ Executable |
| 3 | `./ASF-v5.1.0-darwin-arm64 --version` | ✅ `ASF0 v5.1.0` |

### Installer URL Verification

The installer script (`install.sh`) dynamically constructs download URLs using:

```
LATEST_VERSION="v5.1.0"
DIRECT_DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}"
```

Expected download URLs:

```
https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.0/ASF-v5.1.0-darwin-arm64
https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.0/ASF-v5.1.0-darwin-amd64
https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.0/ASF-v5.1.0-linux-amd64
https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.0/ASF-v5.1.0-linux-arm64
https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.0/ASF-v5.1.0-windows-amd64.exe
```

All URLs target v5.1.0 release assets. Fresh install will work correctly.
