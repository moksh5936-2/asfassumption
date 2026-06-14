# V505_UPGRADE_TEST — ASF0 v5.0.5

## Procedure
1. Install v5.0.4 binary to `~/.asf/asf`
2. Symlink to `~/.local/bin/asf`
3. Run `curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade`

## Results
- Config backed up successfully
- Detected latest version: v5.0.5
- Downloaded correct darwin-arm64 binary
- Checksum verified (SHA-256)
- Binary reported: ASF0 v5.0.5
- All install checks passed

## Final State
```
~/.asf/asf --version → ASF0 v5.0.5
```

Upgrade from v5.0.4 to v5.0.5: ✅ Passed
