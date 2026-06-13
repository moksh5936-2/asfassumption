# Installer Hotfix — Root Cause Analysis

## Symptom

```
$ curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash

Downloading ASF v3.0.0-RC2
3.0.0-RC2 for darwin/arm64...

Generated URL:
https://github.com/moksh5936-2/asfassumption/releases/download/v3.0.0-RC2
3.0.0-RC2/ASF-v3.0.0-RC2
3.0.0-RC2-darwin-arm64

Download failed (HTTP 000)
```

The download URL contains embedded newlines, splitting across 3 lines. The version value
`v3.0.0-RC2` becomes `v3.0.0-RC2\n3.0.0-RC2` (duplicated with newline), and this corrupted
string propagates into the URL, the binary name, and every user-facing message.

## Root Cause

### Bug: Python bare `except:` catches `SystemExit`

File: `install.sh` (and `release/install.sh`)
Line: 369 (original)

```python
try:
    releases = json.load(sys.stdin)
    for r in releases:
        if not r.get('draft', True):
            print(r['tag_name'].lstrip('v'))
            sys.exit(0)
except: pass                          # <--- BARE except
print('3.0.0-RC2')
```

The `except:` bare clause catches **all** exceptions, including `SystemExit(0)` raised by
`sys.exit(0)`. After catching `SystemExit`, execution continues to `print('3.0.0-RC2')`,
which emits a **second** version value on stdout.

The `$()` command substitution in bash strips **trailing** newlines but preserves **embedded**
newlines. The two `print()` calls produce:

```
3.0.0-RC2\n3.0.0-RC2\n
```

After `$()` strips the final `\n`:

```
3.0.0-RC2\n3.0.0-RC2
```

This `VERSION` variable (containing `3.0.0-RC2\n3.0.0-RC2`) is then used in URL construction:

```bash
DIRECT_DOWNLOAD_URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BINARY_NAME}"
```

Which expands to:

```
https://github.com/moksh5936-2/asfassumption/releases/download/v3.0.0-RC2
3.0.0-RC2/ASF-v3.0.0-RC2
3.0.0-RC2-darwin-arm64
```

A 3-line URL containing two copies of the version with an embedded newline. The download
fails because the URL is not a valid HTTP URL.

### Reproduction

```bash
# Demonstrate the bug:
python3 -c "
import sys
try:
    print('v3.0.0-RC2')
    sys.exit(0)
except: pass
print('v3.0.0-RC2')
" | xxd
# Output:
# 00000000: 7633 2e30 2e30 2d52 4332 0a76 332e 302e  v3.0.0-RC2.v3.0.0-
# 00000010: 302d 5243 320a                            0-RC2.
# The 0a bytes are embedded newlines — TWO lines emitted.
```

```bash
# With the fix (except Exception):
python3 -c "
import sys
try:
    print('v3.0.0-RC2')
    sys.exit(0)
except Exception: pass
print('v3.0.0-RC2')
" | xxd
# Output:
# 00000000: 7633 2e30 2e30 2d52 4332 0a              v3.0.0-RC2.
# Only ONE line emitted. SystemExit(0) is not caught.
```

### Contributing Factors

1. **No version normalization** — Even without the bare `except` bug, there was no
   `tr -d '\r\n'` or whitespace stripping to catch stray carriage returns or other
   whitespace that might appear from API responses or environment variables.

2. **Single version variable** — The same variable was used for both the release tag
   (with `v` prefix) and the asset version (without `v` prefix), relying on inline
   `v${VERSION}` conversion throughout the script.

3. **No URL validation** — Before attempting a download, the installer did not verify
   that the constructed URL was free of whitespace or malformed.

## Fix Summary

| Issue | Fix |
|-------|-----|
| `except: pass` catches `SystemExit` | Changed to `except Exception: pass` |
| No version normalization | Added `tr -d '\r\n' \| xargs` after version detection |
| Single version variable ambiguity | `VERSION` (no v), `LATEST_VERSION` (with v), `ASSET_VERSION` (no v) |
| No URL validation | Added whitespace rejection and binary name prefix check before download |
