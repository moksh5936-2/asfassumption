# Contributing to ASF

## Welcome

ASF is an open-source project. We welcome contributions from security researchers, Go developers, threat modelers, and documentation writers.

## Code of Conduct

This project is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). All participants must abide by it.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/asf.git`
3. Set up the development environment:
   ```bash
   cd asf/asf-tui
   go build -o asf-tui .
   ```
4. Install the Python ASF engine:
   ```bash
   cd ..  # cybersec/
   pip install -e .
   ```

## Development Workflow

### Branch naming

- `feature/description` — New features
- `fix/description` — Bug fixes
- `docs/description` — Documentation
- `test/description` — Testing

### Commit style

Use conventional commits:
```
feat: add YAML architecture parser
fix: handle empty evidence array from Python CLI
docs: update risk model documentation
test: add edge case tests for confidence engine
```

### Before submitting

1. Run `go vet ./...` — must produce zero warnings
2. Run `go test ./...` — all tests must pass
3. Build the binary — must compile cleanly

## What to Work On

### High priority

- **Validation study** — Help design and run an expert validation study
- **CI/CD pipeline** — GitHub Actions for build, test, release
- **Code signing** — macOS notarization and code signing

### Medium priority

- **Parser improvements** — Better Draw.io/Mermaid parsing
- **Additional STRIDE rules** — Expand the 33 keyword patterns
- **Additional export formats** — AsciiDoc, DOCX
- **Multi-architecture batch analysis**
- **Results database persistence**

### Low priority

- **Additional themes**
- **Internationalization**
- **Windows TUI testing and fixes**

## Testing

```bash
# Run all tests
cd asf-tui && go test ./... -v

# Run specific test
go test -run TestRiskMatrix -v

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Code Review

All submissions require review. We use GitHub pull requests.

### PR requirements

1. Clear description of what the change does
2. Link to any related issues
3. All tests pass
4. `go vet` clean
5. No regression in existing functionality

## Documentation

Documentation is in the `docs/` directory. When adding features, update:

- `docs/TECHNICAL_REFERENCE.md` — Add new structs, interfaces, services
- `docs/USER_MANUAL.md` — Add new workflows, settings
- `docs/DEVELOPER_GUIDE.md` — Add new extension points

## Questions

Open a GitHub issue or contact the maintainers.

## License

By contributing, you agree that your contributions will be licensed under the project's license.
