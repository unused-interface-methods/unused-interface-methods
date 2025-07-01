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

## 🚀 Quick Start

### Installation

1. **Install golangci-lint:**
   ```bash
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   ```

2. **Build the linter:**
   ```bash
   # Cross-platform build using Makefile
   make build
   ```

### Usage

#### Task (Recommended)

Install [Task](https://taskfile.dev/) and use the modern build system:

```bash
# Install Task (if not already installed)
go install github.com/go-task/task/v3/cmd/task@latest

# View available tasks
task --list

# Quick development workflow
task build          # build binary
task lint-interfaces # run unused interface checker
task lint           # run golangci-lint
task lint-all       # run all linters
task test           # run tests
task dev            # build and run on test files
task clean          # clean build artifacts

# CI pipeline
task ci             # full check: download, tidy, build, test, lint
```

**VSCode Integration** (with keyboard shortcuts):
- `Ctrl+Shift+L` - All linters
- `Ctrl+Shift+I` - Unused interfaces only
- `Ctrl+Shift+G` - Golangci-lint only
- `Ctrl+Shift+B` - Build
- `Ctrl+Shift+T` - Tests
- `Ctrl+Shift+D` - Dev mode
- `Ctrl+Shift+K` - Full check

#### Cross-Platform Scripts (Alternative)

**Windows (PowerShell):**
```powershell
# Run both linters
./lint.ps1

# Run specific linter
./lint.ps1 standard     # golangci-lint only
./lint.ps1 interfaces   # unused interface methods only
./lint.ps1 test         # run tests
./lint.ps1 help         # show help
```

**Unix/Linux/macOS (Shell):**
```bash
# Run both linters
./lint.sh

# Run specific linter
./lint.sh standard     # golangci-lint only
./lint.sh interfaces   # unused interface methods only
./lint.sh test         # run tests
./lint.sh help         # show help
```

#### Manual (Cross-Platform)
```bash
# Build first
make build

# Standard linting
golangci-lint run .

# Unused interface methods (binary name varies by OS)
# On Windows: ./unusedintf.exe ./...
# On Unix/Linux/macOS: ./unusedintf ./...
make lint-interfaces
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

**Using Task (Recommended):**
```yaml
# GitHub Actions example - Task
- name: 📦 Install Task
  run: go install github.com/go-task/task/v3/cmd/task@latest

- name: 🔍 Full pipeline check
  run: task ci

- name: 🔍 Just unused interfaces
  run: task lint-interfaces
```

**Using Makefile (Alternative):**
```yaml
# GitHub Actions example (cross-platform)
- name: 🔍 Check unused interface methods
  run: |
    make build
    make lint-interfaces
```

**Manual Commands:**
```yaml
# Alternative using direct commands
- name: 🔍 Check unused interface methods (manual)
  run: |
    if [[ "$RUNNER_OS" == "Windows" ]]; then
      go build -o unusedintf.exe .
      ./unusedintf.exe ./...
    else
      go build -o unusedintf .
      ./unusedintf ./...
    fi
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
# 🧪 Run tests
go test ./...

# 🔍 Lint the linter
go vet ./...

# 🚀 Test on real projects
./unusedintf ./...
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