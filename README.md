# 🔍 unusedintf – Go Unused Interface Methods Analyzer

[![Go Version](https://img.shields.io/github/go-mod/go-version/Headcrab/lint)](https://golang.org/dl/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/Headcrab/lint)](https://goreportcard.com/report/github.com/Headcrab/lint)
[![CI](https://github.com/Headcrab/lint/workflows/CI/badge.svg)](https://github.com/Headcrab/lint/actions)

> 🚀 **Lightning-fast static analyzer** that hunts down unused interface methods in your Go codebase

## 🎯 Overview

`unusedintf` is a **powerful static analysis tool** for Go that detects interface methods that are **declared but never used** anywhere in your codebase. Built on top of [golang.org/x/tools/go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis), it seamlessly integrates with your development workflow.

### 💡 Why You Need This

- 🧹 **Clean APIs**: Dead interface methods confuse users and bloat your public API
- ⚡ **Faster Builds**: Removing unused code makes compilation faster
- 🔧 **Easier Refactoring**: Less surface area = simpler maintenance
- 🚦 **CI-Ready**: Non-zero exit status when issues are found

## ✨ Features

- 🎯 **Smart Detection**: Finds unused methods on ordinary and **generic** interfaces (Go 1.18+)
- 🧠 **Context-Aware**: Understands complex usage patterns:
  - 📎 Method values & function pointers
  - 🔄 Type assertions & type switches  
  - 📦 Embedded interfaces (bidirectional)
  - 🖨️ `fmt` package implicit `String()` calls
- 📊 **Clean Output**: Sorted by file path and line numbers
- ⚙️ **Configurable**: Optional `-skipGenerics` flag for legacy codebases
- 🔌 **Editor Integration**: Works with `go vet`, `gopls`, and your favorite IDE
- 🌍 **Cross-Platform**: Full support for Windows, Linux, and macOS

## 🚀 Quick Start

### Installation

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Clone and build
git clone https://github.com/Headcrab/lint.git
cd lint
make build
```

## 🔧 Binary Names

The build system automatically detects your platform:

- **Windows**: `unusedintf.exe`
- **Linux/macOS/Unix**: `unusedintf`

## 🎮 Usage Options

### Option 1: Task (Recommended)

Install [Task](https://taskfile.dev/) for the modern build system:

```bash
# Install Task
go install github.com/go-task/task/v3/cmd/task@latest

# Available tasks
task --list

# Development workflow
task build           # build binary
task lint-interfaces # run unused interface checker
task lint            # run golangci-lint
task lint-all        # run all linters
task test            # run tests
task dev             # build and run on test files
task clean           # clean build artifacts
task ci              # full CI pipeline
```

**VSCode Integration** (keyboard shortcuts):
- `Ctrl+Shift+L` - All linters
- `Ctrl+Shift+I` - Unused interfaces only
- `Ctrl+Shift+G` - Golangci-lint only
- `Ctrl+Shift+B` - Build
- `Ctrl+Shift+T` - Tests
- `Ctrl+Shift+D` - Dev mode
- `Ctrl+Shift+K` - Full check

## 🔧 VS Code Integration

This linter integrates seamlessly with **"Go: Lint Workspace"** functionality:

### Method 1: Integrated Wrapper (Recommended)

Use the provided wrapper scripts for automatic execution:

```bash
# Windows
./lint-wrapper.ps1                    # run both linters
./lint-wrapper.ps1 -InterfacesOnly    # run only unusedintf
./lint-wrapper.ps1 -StandardOnly      # run only golangci-lint

# Linux/Mac  
./lint-wrapper.sh                     # run both linters
./lint-wrapper.sh --interfaces-only   # run only unusedintf
./lint-wrapper.sh --standard-only     # run only golangci-lint
```

**VS Code Integration:**
1. `Ctrl+Shift+P` → `Tasks: Run Task` → `Go: Lint Workspace (Integrated)`
2. Or run directly: `task lint-wrapper`

### Method 2: Manual VS Code Tasks

Available VS Code tasks (Command Palette):
- `Go: Lint Interfaces` - Run only interface linter  
- `Go: Lint Standard` - Run only golangci-lint
- `Go: Lint All` - Run both linters sequentially
- `Go: Lint Workspace (Integrated)` - Cross-platform wrapper

### Method 3: Editor Configuration  

Configure in `.vscode/settings.json`:

```json
{
    "go.lintTool": "staticcheck",
    "go.lintOnSave": "workspace"
}
```

### ✨ Features

The VS Code integration provides:
- ✅ **Real-time highlighting** of unused interface methods
- ✅ **Problems panel** integration with clickable errors
- ✅ **File explorer markers** showing files with issues
- ✅ **Cross-platform compatibility** (Windows/Linux/Mac)
- ✅ **Automatic problem matching** for both linters

### Option 2: Makefile (Universal)

```bash
make build           # build cross-platform binary
make lint            # run both linters
make lint-interfaces # run unused interface methods linter
make test            # run tests  
make clean           # remove build artifacts
make help            # show available targets
```

### Option 3: Manual Commands

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

## 📋 Sample Output

```
D:\work\go\lint\test\interfaces.go:38:2: method "OnError" of interface "EventHandler" is declared but not used
D:\work\go\lint\test\interfaces.go:39:2: method "Subscribe" of interface "EventHandler" is declared but not used
```

> 💡 **Pro Tip**: Output format is identical to `go vet` - your editor will highlight issues automatically!

## 🔧 Integration

### With your own analyzer

```go
import "github.com/Headcrab/lint"

// Add to your multichecker
analyzers := []*analysis.Analyzer{
    unusedintf.Analyzer,
    // ... your other analyzers
}
```

### CI/CD Pipeline

#### GitHub Actions

**Using Task (Recommended):**
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
        
    - name: Install Task
      run: go install github.com/go-task/task/v3/cmd/task@latest
      
    - name: Full CI pipeline
      run: task ci
```

**Using Makefile:**
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
      
    - name: Run build and lint
      run: |
        make build
        make lint
```

#### GitLab CI

```yaml
stages:
  - lint

variables:
  GO_VERSION: "1.24"

.lint_template: &lint_template
  image: golang:${GO_VERSION}
  before_script:
    - go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    - go install github.com/go-task/task/v3/cmd/task@latest
  script:
    - task ci

lint:
  <<: *lint_template
```

## ⚠️ Known Limitations

- 🪞 **Reflection**: Cannot track `reflect.Value.Call()` usage
- 🤖 **Code Generation**: Dynamic/generated code is not analyzed
- 🔌 **Plugins**: Runtime plugin loading is not tracked
- 🧪 **Generics**: Best-effort matching; edge cases may slip through

## 🤝 Contributing

We ❤️ contributions! Please include:

1. 🐛 **Reproducer** (code snippet or minimal repo)
2. 📊 **Expected vs actual output**
3. 🔖 **Go version** (`go version`)

### Development

```bash
# Run tests
task test
# or
make test

# Run linting
task lint-all
# or  
make lint

# Test on real projects
task dev
```

## 📊 Stats

![GitHub stars](https://img.shields.io/github/stars/Headcrab/lint?style=social)
![GitHub forks](https://img.shields.io/github/forks/Headcrab/lint?style=social)
![GitHub watchers](https://img.shields.io/github/watchers/Headcrab/lint?style=social)

## 📄 License

MIT © [Headcrab](https://github.com/Headcrab/lint) - see [LICENSE](LICENSE) for details.

---

<div align="center">

**[⭐ Star this repo](https://github.com/Headcrab/lint)** if it helped you write cleaner Go code!

Made with ❤️ for the Go community

</div>

## Options

- `-skipGenerics` - Skip generic interface analysis 