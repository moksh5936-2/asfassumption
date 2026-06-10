# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 1.0.x   | ✅ Active          |
| < 1.0   | ❌ Not supported   |

## Reporting a Vulnerability

ASF is a security tool, and we take vulnerabilities seriously.

**Do not open public issues for security vulnerabilities.**

Instead, report them to: security@asfsecurity.com

We will:
1. Acknowledge receipt within 48 hours
2. Provide an initial assessment within 5 business days
3. Release a fix based on severity
4. Credit the reporter (if desired)

## Scope

The following are in scope for security reports:

- The Go TUI application (`asf-tui/`)
- The Python ASF engine (`asf/`)
- The build and release pipeline

The following are out of scope:

- Ollama vulnerabilities (report to Ollama)
- Tesseract vulnerabilities (report to Tesseract)
- Third-party Go dependency vulnerabilities

## Security Features

- **No telemetry**: ASF does not phone home or collect usage data
- **No cloud dependency**: All processing is local
- **Offline license validation**: HMAC-based, no network required
- **Deterministic output**: Same input always produces same output — no hidden behavior

## Known Security Considerations

1. License HMAC key (`asf-enterprise-secret-2024`) is embedded in the binary. This is a symmetric key and could be extracted. Enterprise deployments should use asymmetric signing or online validation.
2. Python CLI bridge executes a subprocess. Ensure the Python path is not attacker-controlled.
3. AI enhancement calls `localhost:11434` (Ollama API). No external network calls.
4. OCR runs `tesseract` as a subprocess. Ensure the tesseract path is not attacker-controlled.
