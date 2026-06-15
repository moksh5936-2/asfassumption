.PHONY: all build install clean test lint run doctor docs completion-bash completion-zsh completion-fish

APP     := asf
VERSION ?= 5.1.0
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ 2>/dev/null || echo "unknown")
LDFLAGS := -ldflags "-X main.ASFVersion=$(VERSION) -X main.buildCommit=$(COMMIT) -X main.buildDate=$(DATE)"

all: build

build:
	go build $(LDFLAGS) -o $(APP) .

install:
	go install $(LDFLAGS) ./...

clean:
	rm -f $(APP)
	rm -rf dist/

test:
	go test ./... -v -count=1

test-short:
	go test ./... -short -count=1

test-coverage:
	go test ./... -coverprofile=coverage.out -count=1
	go tool cover -html=coverage.out -o coverage.html

lint:
	golangci-lint run ./... 2>/dev/null || staticcheck ./... 2>/dev/null || go vet ./...

run:
	./$(APP)

doctor:
	./$(APP) doctor

docs:
	go doc -all ./...

release: build
	mkdir -p dist
	cp $(APP) dist/$(APP)-$(VERSION)
	cd dist && shasum -a 256 $(APP)-$(VERSION) > $(APP)-$(VERSION).sha256
	echo "Release $(VERSION) built in dist/"

completion-bash:
	./$(APP) completion bash > dist/completion.bash 2>/dev/null || \
	echo '# eval "$(asf completion bash)"' > dist/completion.bash
	@echo "Saved to dist/completion.bash"

completion-zsh:
	./$(APP) completion zsh > dist/completion.zsh 2>/dev/null || \
	echo '# asf completion zsh > /usr/local/share/zsh/site-functions/_asf' > dist/completion.zsh
	@echo "Saved to dist/completion.zsh"

completion-fish:
	./$(APP) completion fish > dist/completion.fish 2>/dev/null || \
	echo '# asf completion fish > ~/.config/fish/completions/asf.fish' > dist/completion.fish
	@echo "Saved to dist/completion.fish"

dist: | dist/completion.bash dist/completion.zsh dist/completion.fish

help:
	@echo "ASF0 v$(VERSION) — Makefile targets:"
	@echo "  build               Build the binary"
	@echo "  install             Install to GOPATH/bin"
	@echo "  clean               Remove build artifacts"
	@echo "  test                Run all tests"
	@echo "  test-short          Run short tests only"
	@echo "  test-coverage       Generate coverage report"
	@echo "  lint                Run linters"
	@echo "  run                 Run the TUI"
	@echo "  doctor              Run diagnostics"
	@echo "  release             Build and prepare release"
	@echo "  completion-bash     Generate bash completion"
	@echo "  completion-zsh      Generate zsh completion"
	@echo "  completion-fish     Generate fish completion"
