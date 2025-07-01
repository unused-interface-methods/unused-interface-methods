# Detect OS and set binary extension
ifeq ($(OS),Windows_NT)
    BINARY_NAME = unusedintf.exe
    RM_CMD = del /f /q
else
    BINARY_NAME = unusedintf
    RM_CMD = rm -f
endif

.PHONY: lint build test clean all

# Build the custom linter
build:
	go build -o $(BINARY_NAME) .

# Run golangci-lint
lint-standard:
	golangci-lint run .

# Run our custom unused interface method linter
lint-interfaces: build
	./$(BINARY_NAME) ./...

# Run both linters
lint: lint-standard lint-interfaces

# Run tests
test:
	go test -v
	go test -bench=.

# Clean build artifacts
clean:
	$(RM_CMD) $(BINARY_NAME)

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