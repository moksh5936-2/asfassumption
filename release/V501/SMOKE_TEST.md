# v5.0.1 Smoke Test

## Binary Tested
`dist/ASF-v5.0.1-darwin-arm64`

## --version
```
$ ./dist/ASF-v5.0.1-darwin-arm64 --version
ASF0 v5.0.1
```
PASS

## --help
Displays usage, configuration paths, documentation URL. PASS.

## doctor
```
ASF Doctor — System Diagnostic
  OS:               darwin
  Architecture:     arm64
  Version:          5.0.1
  Theme:            ASF0
  AI enabled:       false
  Ollama running:   yes
```
PASS — all diagnostics pass, version confirmed.

## Result
PASS — binary executes correctly, version reports v5.0.1.
