# Detect OS and set binary extension
ifeq ($(OS),Windows_NT)
    BINARY_NAME = unused-interface-methods.exe
    RM_CMD = del /f /q
else
    BINARY_NAME = unused-interface-methods
    RM_CMD = rm -f
endif

.PHONY: lint build test bench clean try

# Run golangci-lint
lint:
	golangci-lint run .

# Build our custom unused-interface-methods linter
build:
	go build -o $(BINARY_NAME) .

# Run tests
test:
	go test -v

# Run benchmarks
bench:
	go test -bench=.

# Clean build artifacts
clean:
	$(RM_CMD) $(BINARY_NAME)

# Run our custom unused-interface-methods linter
try: build
	./$(BINARY_NAME) ./...

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build unused-interface-methods linter"
	@echo "  try             - Run unused-interface-methods linter"
	@echo "  lint            - Run golangci-lint"
	@echo "  test            - Run tests"
	@echo "  bench           - Run benchmarks"
	@echo "  clean           - Remove build artifacts"
