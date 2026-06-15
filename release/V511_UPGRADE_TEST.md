# V511 — Upgrade Test

## Upgrade command

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

Expected result: Existing `asf` binary replaced with `ASF0 v5.1.1`.

## Version after upgrade

```
$ asf --version
ASF0 v5.1.1
```

## Verdict

Upgrade from v5.1.0 (or older) installs v5.1.1.

**UPGRADE_TEST_VERIFIED**
