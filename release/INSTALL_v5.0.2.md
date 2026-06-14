# ASF0 v5.0.2 — Install Guide

## macOS / Linux

### Quick install
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

### Upgrade from older version
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

### Manual install (from release assets)
1. Download the appropriate binary from the GitHub release:
   - macOS Apple Silicon: `ASF-v5.0.2-darwin-arm64`
   - macOS Intel: `ASF-v5.0.2-darwin-amd64`
   - Linux amd64: `ASF-v5.0.2-linux-amd64`
   - Linux arm64: `ASF-v5.0.2-linux-arm64`
2. `chmod +x ASF-v5.0.2-<arch>`
3. `mkdir -p ~/.local/bin && mv ASF-v5.0.2-<arch> ~/.local/bin/asf`
4. Ensure `~/.local/bin` is in your PATH.

## Windows

1. Download `ASF-v5.0.2-windows-amd64.exe` from the GitHub release.
2. Rename to `asf.exe`.
3. Place in a directory in your PATH (e.g., `C:\Windows\System32\` or create `%USERPROFILE%\bin\` and add to PATH).
4. Open a new Command Prompt or PowerShell and run `asf`.

## Verify installation
```bash
asf --version
# Expected: ASF0 v5.0.2
```

## Troubleshooting

### Wrong version shown
- Run `asf doctor --fix` to remove stale binaries.
- Ensure `which asf` points to the v5.0.2 binary location.

### "asf: command not found"
- Add `~/.local/bin` to your PATH:
  ```bash
  echo 'export PATH="$PATH:$HOME/.local/bin"' >> ~/.zshrc
  source ~/.zshrc
  ```

### Stale binary from old install
- Run `asf doctor --fix` or manually remove old copies:
  ```bash
  rm -f ~/.asf/asf
  ```
  Then reinstall.

### Broken URL / download fails
- Check GitHub release exists at: https://github.com/moksh5936-2/asfassumption/releases/tag/v5.0.2
- Check your internet connection.
- Try manual download and placement.
