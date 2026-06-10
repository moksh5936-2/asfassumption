# Security & Portability Review

## Path Traversal Risk Assessment

### File Loading (User-Provided Paths)

| Location | How Path is Used | Risk | Mitigation |
|----------|-----------------|------|------------|
| `engine.go:RunAnalysis()` | Architecture file path from user | Low | User-provided, CWD-relative |
| `parser.go:ParseArchitecture()` | Reads file at user-provided path | Low | `os.Open()` on user path |
| `export.go:ExportResult()` | Writes to `outputDir` parameter | Low | Parameter comes from config or user |

### User-Controlled Paths

All file loading uses paths provided by the user through the TUI input. No normalization or validation is performed beyond checking file existence.

Risk: **Low** — The user is the one providing paths. ASF operates as a desktop application, not a network service. There's no privilege boundary to cross.

### Exports

Exports go to the configured `output.directory` (default `./reports`). The user controls this path through the Settings UI. No sanitization is applied.

Risk: **Low** — Same justification as above.

## User-Provided Path Validation

Currently, user-provided paths are:
- Not normalized (could contain `..` components)
- Not validated for type (could be a directory, symlink, etc.)
- Only checked for existence

Recommendation: Add `filepath.Clean()` normalization and a `filepath.IsAbs()` check for non-relative paths. This is LOW priority for a desktop TUI application.

## Platform-Specific Security

| Platform | Concern | Status |
|----------|---------|--------|
| Linux | `~/Library/Application Support` doesn't exist | ❌ **Not applicable** — Linux uses `~/.config/asf` |
| macOS | TCC permissions | ASF stores config in `~/Library/Application Support/asf` which may trigger TCC prompts on sandboxed builds. Non-sandboxed binaries are unaffected. |
| Windows | Admin rights | Installer uses `%LOCALAPPDATA%` — no admin required. Binary runs as user. |

## Binary Security

- Static binary with `CGO_ENABLED=0`
- Stripped with `-ldflags="-s -w"`
- No runtime code generation
- No dynamic linking

## Supply Chain

- Go dependencies from `go.sum` (checksum-verified)
- Release binaries verified via SHA-256 checksums
- Install scripts support checksum verification
- GitHub Actions workflow builds from tagged commits

## Recommendations

1. Add `filepath.Clean()` to user-provided paths (low priority)
2. Consider code signing for macOS releases (medium priority)
3. Consider Windows Authenticode signing (low priority)
