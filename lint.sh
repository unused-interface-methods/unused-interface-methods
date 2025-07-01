#!/bin/bash

set -e

# colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # no color

# detect binary name based on OS
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "win32" || "$OS" == "Windows_NT" ]]; then
  BINARY_NAME="unusedintf.exe"
else
  BINARY_NAME="unusedintf"
fi

print_help() {
  echo -e "${CYAN}Available commands:${NC}"
  echo "  all (default)  - run tests and both linters"
  echo "  standard       - run golangci-lint only"
  echo "  interfaces     - run unused interface methods linter only"
  echo "  test           - run tests and benchmarks"
  echo "  build          - build the unusedintf linter"
  echo "  clean          - remove build artifacts"
  echo "  help           - show this help"
  echo ""
  echo "Usage: ./lint.sh [command]"
  echo "Examples:"
  echo "  ./lint.sh"
  echo "  ./lint.sh standard"
  echo "  ./lint.sh interfaces"
}

build_linter() {
  echo -e "${YELLOW}Building unusedintf linter...${NC}"
  if go build -o "$BINARY_NAME" .; then
    echo -e "${GREEN}✅ Build successful${NC}"
  else
    echo -e "${RED}❌ Build failed${NC}"
    exit 1
  fi
}

run_standard_lint() {
  echo -e "${YELLOW}Running golangci-lint...${NC}"
  if golangci-lint run .; then
    echo -e "${GREEN}✅ golangci-lint passed${NC}"
  else
    echo -e "${YELLOW}⚠️  golangci-lint found issues${NC}"
  fi
}

run_interface_lint() {
  echo -e "${YELLOW}Running unused interface methods linter...${NC}"
  build_linter
  if "./$BINARY_NAME" ./...; then
    echo -e "${GREEN}✅ No unused interface methods${NC}"
  else
    echo -e "${YELLOW}⚠️  Found unused interface methods${NC}"
  fi
}

run_tests() {
  echo -e "${YELLOW}Running tests...${NC}"
  go test -v
  echo -e "${YELLOW}Running benchmarks...${NC}"
  go test -bench=.
}

clean_artifacts() {
  echo -e "${YELLOW}Cleaning build artifacts...${NC}"
  if [[ -f "$BINARY_NAME" ]]; then
    rm -f "$BINARY_NAME"
    echo -e "${GREEN}✅ Cleaned $BINARY_NAME${NC}"
  fi
}

# main logic
case "${1:-all}" in
"build")
  build_linter
  ;;
"standard")
  run_standard_lint
  ;;
"interfaces")
  run_interface_lint
  ;;
"test")
  run_tests
  ;;
"clean")
  clean_artifacts
  ;;
"help")
  print_help
  ;;
"all")
  run_tests
  run_standard_lint
  run_interface_lint
  ;;
*)
  echo -e "${RED}Unknown command: $1${NC}"
  print_help
  exit 1
  ;;
esac
