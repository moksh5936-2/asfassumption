# V5.0.2 — Fresh Install Test

| Step | Result |
|---|---|
| `rm -rf ~/.asf ~/.local/bin/asf` | ✅ Clean slate |
| `curl ... install.sh | bash` | ✅ Installed v5.0.2 |
| `asf --version` | ✅ `ASF0 v5.0.2` |
| Download from GitHub release | ✅ 12M binary, checksum verified |
| Config created | ✅ `~/.asf/config.yaml` |
| Symlink created | ✅ `~/.local/bin/asf → ~/.asf/asf` |
