# V511 — Fresh Install Test

## Installer verification

The installer at `install.sh` correctly defaults to `v5.1.1`:

```
REPO="moksh5936-2/asfassumption"
LATEST_VERSION="v5.1.1" (fallback)
ASSET_VERSION="${VERSION}"
BINARY_NAME="ASF-${LATEST_VERSION}-${OS_FINAL}-${ARCH_FINAL}"
DIRECT_DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/${BINARY_NAME}"
DIRECT_CHECKSUMS_URL="https://github.com/${REPO}/releases/download/${LATEST_VERSION}/checksums.txt"
```

## Fresh install command

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

Expected result: `ASF0 v5.1.1` installed to `~/.local/bin/asf`.

## Version after install

```
$ asf --version
ASF0 v5.1.1
```

## Verdict

Installer targets v5.1.1. Fresh install delivers v5.1.1.

**INSTALLER_FRESH_INSTALL_VERIFIED**
