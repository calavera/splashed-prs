.PHONY: all build clean deps help test

help: ## Show this help.
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: deps test build ## Run tests and build the binary.

build: ## Build the Go binary
	@echo "Building binary in dist"
	@mkdir -p dist
	@GO111MODULE=on CGO_ENABLED=0 go build -o dist/splashed_prs

clean: ## Remove all artifacts.
	@rm -rf dist

deps: ## Install dependencies.
	@echo "Installing dependencies"
	@GO111MODULE=on go mod verify
	@GO111MODULE=on go mod tidy

test: deps ## Run tests.
	@GO111MODULE=on go test -v ./...
