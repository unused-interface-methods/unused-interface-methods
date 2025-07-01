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
   go build -o unusedintf.exe .
   ```

### Usage

#### PowerShell Script (Recommended)
```bash
# Run both linters
./lint.ps1

# Run specific linter
./lint.ps1 standard     # golangci-lint only
./lint.ps1 interfaces   # unused interface methods only
./lint.ps1 test         # run tests
./lint.ps1 help         # show help
```

#### Manual
```bash
# Standard linting
golangci-lint run .

# Unused interface methods
./unusedintf.exe ./...
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

```yaml
# GitHub Actions example
- name: 🔍 Check unused interface methods
  run: |
    go build -o unusedintf ./cmd/unusedintf
    ./unusedintf ./...
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