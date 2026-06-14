# V5.0.2 — Upgrade Test

| Step | Result |
|---|---|
| Install older version (v5.0.0) | ✅ `ASF0 v5.0.0` |
| Run `install.sh --upgrade` | ✅ Upgraded to v5.0.2 |
| `asf --version` after upgrade | ✅ `ASF0 v5.0.2` |
| Config preserved | ✅ Backed up to `.asf/backups/` |
| No stale version | ✅ Only v5.0.2 in PATH |
| No broken URL | ✅ Asset downloaded from `releases/download/v5.0.2/...` |
| No duplicated version string | ✅ Installer shows single `v5.0.2` |
