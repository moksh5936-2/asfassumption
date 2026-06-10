# Error Handling Audit

## Before: Confusing Failures

**Old error message (hardcoded path):**
```
ASF engine error:
chdir /Users/moksh/Project/cybersec:
no such file or directory
```

This failure was:
- Misleading (no indication it was a path configuration issue)
- Machine-specific (referenced a developer's home directory)
- Unrecoverable (user had no way to fix it)

## After: Helpful Diagnostics

**New error messages:**

| Scenario | Error | Next Step |
|----------|-------|-----------|
| Python not found | `ASF engine error: exec: "python3": executable file not found` | Install Python 3.8+ |
| ASF package missing | `ASF engine error: exit status 1 (stderr: No module named asf.cli.main)` | `pip install -e /path/to/asf` |
| Engine configured wrong | `ASF engine error: exec: "/custom/python": stat /custom/python: no such file or directory` | Check `engine.python_path` in config |
| Config file not found | Config loaded with defaults | `asf doctor` shows config path |
| License not found | `No license key found. Place your license key in ~/Library/Application Support/asf/license.key` | Shows the actual path |
| Download fails (404) | `Download failed (HTTP 404). No release binary for ASF v1.0.0 darwin/arm64.` | Check platform support or version |

## `asf doctor` Diagnostic Command

Newly added command that reports:

```
ASF Doctor — System Diagnostic

System
  OS:               darwin
  Architecture:     arm64
  Binary:           /usr/local/bin/asf
  Version:          1.0.0

Paths
  Config directory: /Users/me/Library/Application Support/asf (directory)
  Config file:      /Users/me/Library/Application Support/asf/config.yaml (not found)
  Cache directory:  /Users/me/Library/Caches/asf (directory)
  License file:     /Users/me/Library/Application Support/asf/license.key (not found)

Permissions
  Config dir write: writable
  Cache dir write:  writable

Configuration
  Theme:            Dark
  Analysis depth:   deep
  Export format:    markdown
  ...

Python Engine
  Python binary:    /usr/bin/python3
  ASF engine:       not found

Dependencies
  tesseract:        not found
  ollama:           not found
  python3:          available
```

## Improvements Made

1. **Removed hardcoded paths from error messages** — all paths are now runtime-discovered
2. **Added `asf doctor` command** — single diagnostic entry point
3. **Improved engine error messages** — `cmd.Dir` now uses cache dir instead of hardcoded project dir
4. **Config migration** — automatic migration from old `~/.asf/` and `~/.config/asf/` paths
5. **License path** — shows actual config-path-dependent location
