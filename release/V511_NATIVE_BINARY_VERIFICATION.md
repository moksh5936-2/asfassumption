# V511 — Native Binary Verification

## Platform

darwin-arm64 (native)

## Version

```
$ dist/ASF-v5.1.1-darwin-arm64 --version
ASF0 v5.1.1
```

## Doctor

```
$ dist/ASF-v5.1.1-darwin-arm64 doctor
ASF Doctor — System Diagnostic
- OS: darwin
- Architecture: arm64
- Version: 5.1.1
```

## TUI Verification

Binary launches with expected ASF0 TUI including:
- Result split-pane layout ✓
- Selection-follow viewport behavior ✓
- No dropdown/inline expansion ✓
- SDRI navigation with selection always visible ✓
- Trust navigation with selection always visible ✓
- Detail pane focus on Enter ✓
- Esc returns to list focus ✓
- Resize preserves selection ✓

**NATIVE_BINARY_VERIFIED**
