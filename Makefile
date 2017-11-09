SRC = $(wildcard *.go)
BIN = ws-tcp-relay
DIST_DIR = _dist

$(BIN): $(SRC) fmt ## Build binary for this platform
	go build -o "$(BIN)" $(SRC)

.PHONY: fmt dist clean help

dist: $(SRC) fmt ## Build distribution binaries for all platforms using gox
	@command -v gox 2> /dev/null || go get github.com/mitchellh/gox
	gox -output "$(DIST_DIR)/$(BIN)_{{.OS}}_{{.Arch}}"

fmt: $(SRC) ## Lint with gofmt
	gofmt -w $(SRC)

clean: ## Clean up build files
	rm -rf $(DIST_DIR) $(BIN)

help:
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
