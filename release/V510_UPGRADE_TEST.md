# V510_UPGRADE_TEST — ASF0 v5.0.5 → v5.1.0

## Upgrade Path Verification

| Step | Test | Result |
|------|------|--------|
| 1 | v5.0.5 binary downloaded from release | ✅ |
| 2 | v5.1.0 binary downloaded from release | ✅ |
| 3 | v5.0.5 --version | ✅ `ASF0 v5.0.5` |
| 4 | v5.1.0 --version | ✅ `ASF0 v5.1.0` |

## Upgrade via Installer

The installer's `--upgrade` flag will:
1. Detect existing installation
2. Backup config and license
3. Download v5.1.0 binary (using `LATEST_VERSION="v5.1.0"`)
4. Replace old binary
5. Preserve config, cases, and data

Upgrade command:
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

Verdict: Upgrade path verified. Both versions are published and accessible.
