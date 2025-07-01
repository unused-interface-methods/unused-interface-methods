#!/bin/bash

# Wrapper script for VS Code integration (Linux/Mac)
set -e

BINARY_NAME="unusedintf"
PATH_ARG="${1:-./...}"
STANDARD_ONLY=false
INTERFACES_ONLY=false

# Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
  --standard-only)
    STANDARD_ONLY=true
    shift
    ;;
  --interfaces-only)
    INTERFACES_ONLY=true
    shift
    ;;
  *)
    PATH_ARG="$1"
    shift
    ;;
  esac
done

run_standard_lint() {
  echo "[LINT] Running golangci-lint..."
  golangci-lint run "$PATH_ARG"
}

run_interface_lint() {
  echo "[LINT] Running unusedintf..."

  # Build if necessary
  if [[ ! -f "$BINARY_NAME" ]]; then
    echo "[BUILD] Building unusedintf..."
    go build -o "$BINARY_NAME" .
  fi

  # Run interface linter
  "./$BINARY_NAME" "$PATH_ARG"
}

# Main execution
EXIT_CODE=0

if [[ "$STANDARD_ONLY" == true ]]; then
  run_standard_lint
  EXIT_CODE=$?
elif [[ "$INTERFACES_ONLY" == true ]]; then
  run_interface_lint
  EXIT_CODE=$?
else
  # Run both linters
  echo "[LINT] Running all linters..."

  STANDARD_EXIT=0
  INTERFACE_EXIT=0

  run_standard_lint || STANDARD_EXIT=$?
  run_interface_lint || INTERFACE_EXIT=$?

  # Return non-zero if any linter failed
  EXIT_CODE=$((STANDARD_EXIT > INTERFACE_EXIT ? STANDARD_EXIT : INTERFACE_EXIT))
fi

echo ""
if [[ $EXIT_CODE -eq 0 ]]; then
  echo "[SUCCESS] All linters passed!"
else
  echo "[ERROR] Linting failed with exit code: $EXIT_CODE"
fi

exit $EXIT_CODE
