.PHONY: lint build test clean all

# Build the custom linter
build:
	go build -o unusedintf.exe .

# Run golangci-lint
lint-standard:
	golangci-lint run .

# Run our custom unused interface method linter
lint-interfaces: build
	./unusedintf.exe ./...

# Run both linters
lint: lint-standard lint-interfaces

# Run tests
test:
	go test -v
	go test -bench=.

# Clean build artifacts
clean:
	rm -f unusedintf.exe

# Run everything
all: test lint

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the unusedintf linter"
	@echo "  lint-standard   - Run golangci-lint"
	@echo "  lint-interfaces - Run unused interface methods linter"
	@echo "  lint            - Run both linters"
	@echo "  test            - Run tests and benchmarks"
	@echo "  clean           - Remove build artifacts"
	@echo "  all             - Run tests and linting" 