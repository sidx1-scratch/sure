VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
DATE    ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -s -w -X github.com/sidx1/sure/internal/version.Version=$(VERSION) \
           -X github.com/sidx1/sure/internal/version.Commit=$(COMMIT) \
           -X github.com/sidx1/sure/internal/version.Date=$(DATE)

.PHONY: build install uninstall clean test

build:
	go build -ldflags "$(LDFLAGS)" -o sure .

install: build
	sudo cp sure /usr/local/bin/sure
	@echo "sure installed to /usr/local/bin/sure"
	@echo ""
	@echo "Add to your shell config:"
	@echo '  Bash: eval "$$(sure shell-init --shell bash)"'
	@echo '  Zsh:  eval "$$(sure shell-init --shell zsh)"'

uninstall:
	sudo rm -f /usr/local/bin/sure
	rm -rf ~/.config/sure
	rm -rf ~/.cache/sure
	@echo "sure uninstalled"

clean:
	rm -f sure
	go clean

test:
	go test ./...
