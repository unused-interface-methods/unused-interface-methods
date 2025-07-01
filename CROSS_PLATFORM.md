# Cross-Platform Usage Guide

This project now supports cross-platform development on Windows, Linux, and macOS.

## 🚀 Quick Start

### Option 1: Use Platform-Specific Scripts (Recommended)

**Windows:**
```powershell
./lint.ps1          # run all linters and tests
./lint.ps1 build    # build only
./lint.ps1 help     # show help
```

**Linux/macOS/Unix:**
```bash
./lint.sh           # run all linters and tests  
./lint.sh build     # build only
./lint.sh help      # show help
```

### Option 2: Use Makefile (Universal)

```bash
make build          # build cross-platform binary
make lint           # run both linters
make test           # run tests  
make clean          # remove build artifacts
make help           # show available targets
```

## 🔧 Binary Names

The build system automatically detects your platform:

- **Windows**: `unusedintf.exe`
- **Linux/macOS/Unix**: `unusedintf`

## 🎯 Available Commands

All scripts support these commands:

| Command      | Description                          |
| ------------ | ------------------------------------ |
| `all`        | Run tests and both linters (default) |
| `standard`   | Run golangci-lint only               |
| `interfaces` | Run unused interface methods linter  |
| `test`       | Run tests and benchmarks             |
| `build`      | Build the unusedintf linter          |
| `clean`      | Remove build artifacts               |
| `help`       | Show help message                    |

## 🔍 Manual Usage

If you prefer to run commands manually:

**Build:**
```bash
# Windows
go build -o unusedintf.exe .

# Linux/macOS
go build -o unusedintf .
```

**Run:**
```bash
# Windows  
./unusedintf.exe ./...

# Linux/macOS
./unusedintf ./...
```

## 🚀 CI/CD Integration

### GitHub Actions

```yaml
name: Cross-Platform Lint

on: [push, pull_request]

jobs:
  lint:
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    runs-on: ${{ matrix.os }}
    
    steps:
    - uses: actions/checkout@v4
    
    - name: Setup Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.24'
        
    - name: Install golangci-lint
      run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
      
    - name: Run linting (Unix)
      if: runner.os != 'Windows'
      run: |
        chmod +x lint.sh
        ./lint.sh
        
    - name: Run linting (Windows)  
      if: runner.os == 'Windows'
      run: ./lint.ps1
```

### GitLab CI

```yaml
stages:
  - lint

variables:
  GO_VERSION: "1.24"

.lint_template: &lint_template
  image: golang:${GO_VERSION}
  before_script:
    - go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
  script:
    - make build
    - make lint

lint:linux:
  <<: *lint_template
  tags:
    - linux

lint:windows:
  <<: *lint_template  
  tags:
    - windows
  script:
    - ./lint.ps1
```

## 📝 Notes

- PowerShell is available on all platforms (Windows, Linux, macOS)
- Shell script requires bash/sh (available on most Unix systems)
- Makefile works on all platforms with make installed
- Go code itself is fully cross-platform 