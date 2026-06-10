# Runtime Engine Packaging

ASF consists of two runtime components:

| Component | Language | Purpose | File |
|-----------|----------|---------|------|
| **Go TUI** | Go | Terminal UI, CLI, analysis orchestration | `asf-tui/` |
| **Python Engine** | Python | Claim extraction, assumption analysis, evidence verification | `asf/` |

## Installation Layout

### Go TUI Binary

Installed to a platform-specific binary directory and symlinked into `PATH`:

```
~/.asf/asf                              ← canonical install location
~/.local/bin/asf → ~/.asf/asf           ← symlink (Linux)
```

### Python Engine Bundle

Packaged as `asf-python-engine-v{VERSION}.tar.gz` and installed alongside the Go TUI:

| Platform | Engine Directory |
|----------|-----------------|
| Linux    | `~/.local/share/asf/engine/` |
| macOS    | `~/Library/Application Support/asf/engine/` |
| Windows  | `%LOCALAPPDATA%/ASF/engine/` |

The engine directory contains:

```
engine/
├── asf/               ← Python package (asf.cli.main, etc.)
├── pyproject.toml     ← Build metadata
└── setup.py           ← Minimal setup
```

### How the Go TUI Finds the Engine

At runtime, `engine.go` sets `PYTHONPATH` to the engine directory before spawning the Python subprocess:

```go
cmd.Env = append(cmd.Env, "PYTHONPATH="+asfEngineDir())
```

This makes `import asf` resolve to the bundled engine without requiring `pip install`.

## Doctor Diagnostics

`asf doctor` now clearly distinguishes between the two components:

```
ASF Doctor — System Diagnostic

── System ──
  OS:               linux
  Architecture:     amd64
  Version:          1.1.0
  Go binary (invoked):  /home/user/.local/bin/asf
  Go binary (resolved): /home/user/.asf/asf           ← (if symlinked)

── ASF Go TUI Binary ──
  /home/user/.asf/asf                ASF v1.1.0 [ACTIVE]

── ASF Python Engine ──
  Python binary:    /usr/bin/python3
  Python version:   Python 3.11.0
  Engine directory: /home/user/.local/share/asf/engine ✓
  ASF module version: 0.1.0 ✓
```

### Troubleshooting

**Missing engine:**
```
── ASF Python Engine ──
  Engine directory: /home/user/.local/share/asf/engine ✗
  ASF module status: not installed
```

Fix: `asf doctor --fix` downloads and extracts the engine automatically.

**Broken import:**
```
── ASF Python Engine ──
  Engine directory: /home/user/.local/share/asf/engine ✓
  ASF module status: NOT importable (PYTHONPATH issue)
```

Fix: Check that `PYTHONPATH` includes the engine directory, or re-run `asf doctor --fix`.

## Release Artifacts

Each GitHub Release includes:

| Artifact | Description |
|----------|-------------|
| `ASF-v{VERSION}-{OS}-{ARCH}` | Go TUI binary for each platform |
| `asf-python-engine-v{VERSION}.tar.gz` | Portable Python engine bundle (platform-agnostic) |
| `checksums.txt` | SHA-256 checksums for all assets |

The Python engine is packaged by `scripts/package-python-engine.sh` in CI:

```bash
./scripts/package-python-engine.sh 1.1.0
# → release/asf-python-engine-v1.1.0.tar.gz
```

The tarball includes only runtime-essential files:
- `asf/` package directory
- `pyproject.toml`
- `setup.py`

Build artifacts (`__pycache__`, `.pyc`, `*.egg-info`) are stripped.

## Installer Flow

`install.sh` (and `install.ps1` on Windows):

1. Detect OS/architecture
2. Download `ASF-v{VERSION}-{OS}-{ARCH}` binary
3. Verify SHA-256 checksum
4. Install binary to `~/.asf/asf` with `~/.local/bin/asf` symlink
5. Download `asf-python-engine-v{VERSION}.tar.gz`
6. Extract to data directory (e.g., `~/.local/share/asf/engine/`)
7. Run `asf doctor` post-install verification
8. Print install summary

## doctor --fix

`asf doctor --fix` performs two operations:

1. **Duplicate binary cleanup** — removes stale ASF binaries from `PATH`, keeping only the active one
2. **Python engine repair** — detects missing/broken engine and downloads/extracts it automatically

If the engine download fails (no internet, old release), doctor suggests running `install.sh` again.
