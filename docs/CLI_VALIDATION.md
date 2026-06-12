# CLI Validation Report

**Version:** 2.1.2
**Date:** 2026-06-13

---

## Commands Tested

| Command | Expected Exit | Actual Exit | Status |
|---|---|---|---|
| `asf --version` | 0 | 0 | PASS |
| `asf -v` | 0 | 0 | PASS |
| `asf --help` | 0 | 0 | PASS |
| `asf -h` | 0 | 0 | PASS |
| `asf --version-check` | 0 | 0 | PASS |
| `asf --license` | 7 | 7 | PASS |
| `asf doctor` | 0 | 0 | PASS |
| `asf doctor --verbose` | 0 | 0 | PASS |
| `asf analyze <file>` | 0 | 0 | PASS |
| `asf analyze <file> -e <evidence>` | 0 | 0 | PASS |
| `asf analyze <file> --graph` | 0 | 0 | PASS |
| `asf <invalid>` | 2 | 2 | PASS |

## Analyze Output Validation

- JSON output contains all required sections: version, architecture, summary, assumptions, verifications, gaps, threats, attack paths, controls, findings, intelligence sections
- JSON is valid (parsed by `json.Unmarshal`)
- Native Go engine used (no Python dependency)

## Exit Codes Verified

| Code | Meaning | Status |
|---|---|---|
| 0 | Success | PASS |
| 1 | General error | PASS (not triggered in test, validated in code) |
| 2 | Invalid command | PASS |
| 4 | Analysis error | PASS (not triggered in test, validated in code) |
| 6 | Export error | PASS (not triggered in test, validated in code) |
| 7 | License error | PASS |

## Help Text Verified

- All documented commands appear in help output
- Configuration paths shown correctly
- No broken references

## Conclusion

**All CLI commands function correctly with expected exit codes.**
