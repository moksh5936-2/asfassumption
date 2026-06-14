# V504 — TUI Acceptance

| Check | Status |
|---|---|
| Binary launches without crash | ✅ (error: no TTY — expected in non-interactive env) |
| `--version` shows ASF0 v5.0.4 | ✅ |
| `--help` shows usage | ✅ |
| `doctor` runs clean | ✅ |
| Startup splash screen | ⚠️ Requires real terminal |
| Sidebar renders | ⚠️ Requires real terminal |
| Case workspace opens | ⚠️ Requires real terminal |
| Contradictions tab | ⚠️ Requires real terminal |

## Limitation

Full TUI acceptance requires a real terminal session with a working `/dev/tty`. This environment does not provide one. The binary built and runs — confirmed by `--version`, `--help`, and `doctor`. Recommend performing a full TUI walkthrough on a developer workstation before final distribution.
