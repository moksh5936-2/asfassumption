# V504 — Installer Audit

| Installer | Version | URL Pattern | Status |
|---|---|---|---|
| `asf-tui/install.sh` | `5.0.4` | Dynamic from `LATEST_VERSION` | ✅ |
| `install.sh` (root) | `v5.0.4` (fallback) | Dynamic from `LATEST_VERSION` | ✅ |
| `release/install.sh` | `v5.0.4` (fallback) | Dynamic from `LATEST_VERSION` | ✅ |

All installers construct download URLs dynamically:

```
https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.4/ASF-v5.0.4-{OS}-{ARCH}
https://github.com/moksh5936-2/asfassumption/releases/download/v5.0.4/checksums.txt
```

No stale version strings. No duplicated versions. No broken URL patterns.
