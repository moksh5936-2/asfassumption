# Legacy Python Engine Reference

**Status: ARCHIVED — NOT PART OF PRODUCTION RUNTIME**

## Purpose

ASF v2.x production runtime is the Go-native single binary (`asf-tui/`).
The Python REST API at `asf/` is retained only as historical reference for
the original engine design and migration traceability.

## What This Covers

- **Python REST API** (`asf/api/`) — FastAPI endpoints for document analysis
- **Python CLI** (`asf/cli/`) — Legacy command-line interface
- **Python engine modules** — All subpackages under `asf/` (models, extraction,
  assumption, evidence, verification, confidence, gaps, graph, ingestion, analyzer)

## Config Separation

| System | Config Path | Format |
|--------|------------|--------|
| Python API (legacy) | `asf.config.yaml` in CWD | YAML with db_path, evidence_schema, llm settings |
| Go TUI (production) | `~/.asf/config.yaml` | YAML with theme, analysis, ai, output settings |

These are **separate config systems** and do not share state. Changes made
via the Go TUI have no effect on the Python API, and vice versa.

## How to Use This Reference

The Python source is preserved so that:

1. The original engine architecture can be studied
2. Migration decisions can be traced back to their Python equivalents
3. Any future re-implementation can reference the original design

**Do not** attempt to run the Python API as part of ASF v2.x. Use the
Go binary (`asf`) for all production analysis.

## Key Files

- `asf/config.py` — Python config loader for `asf.config.yaml`
- `asf/settings.py` — Python settings (separate from Go config)
- `asf/api/app.py` — FastAPI application (historical)
- `asf/cli/main.py` — Legacy CLI entry point
