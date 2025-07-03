# unused-interface-methods

[![Go Version](https://img.shields.io/github/go-mod/go-version/unused-interface-methods/unused-interface-methods)](https://go.dev/doc/install)
[![Go Report Card](https://goreportcard.com/badge/github.com/unused-interface-methods/unused-interface-methods)](https://goreportcard.com/report/github.com/unused-interface-methods/unused-interface-methods)
[![coverage](https://img.shields.io/badge/coverage-89.9%25-brightgreen)](https://htmlpreview.github.io/?https://github.com/unused-interface-methods/unused-interface-methods/blob/main/.coverage/.html)
[![Last Commit](https://img.shields.io/github/last-commit/unused-interface-methods/unused-interface-methods)](https://github.com/unused-interface-methods/unused-interface-methods/commits/main/)
[![Project Status](https://img.shields.io/github/release/unused-interface-methods/unused-interface-methods.svg)](https://github.com/unused-interface-methods/unused-interface-methods/releases/latest)

> 🚀 **Lightning-fast static analyzer** that hunts down unused interface methods in your Go codebase

## 🎯 Overview

`unused-interface-methods` is a **powerful static analysis tool** for Go that detects interface methods that are **declared but never used** anywhere in your codebase. Built on top of [golang.org/x/tools/go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis), it seamlessly integrates with your development workflow.

### 💡 Why You Need This

- 🧹 **Clean APIs**: Dead interface methods confuse users and bloat your public API
- ⚡ **Faster Builds**: Removing unused code makes compilation faster
- 🔧 **Easier Refactoring**: Less surface area = simpler maintenance
- 🚦 **CI-Ready**: Non-zero exit status when issues are found

## 🤔 Problem Example

```go
package some_name

// Partially implemented interface
type Interface interface {
    SomeMethod()
    UnusedMethod() // unused for SomeObject - can be removed
}

type SomeObject struct {
    i Interface
}

func (o *SomeObject) SomeMethod() {
    o.i.SomeMethod() // definitely used interface method
}
```

## ✨ Features

- 🎯 **Smart Detection**: Finds unused methods on ordinary and **generic** interfaces (Go 1.18+)
- 🧠 **Context-Aware**: Understands complex usage patterns:
  - 📎 Method values & function pointers
  - 🔄 Type assertions & type switches  
  - 📦 Embedded interfaces (bidirectional)
  - 🖨️ `fmt` package implicit `String()` calls
- 📊 **Clean Output**: Sorted by file path and line numbers
- 🔌 **Editor Integration**: Works with `go vet`, `gopls`, and your favorite IDE
- 🌍 **Cross-Platform**: Full support for Windows, Linux, and macOS


## 🚀 Quick Start

```bash
# Install the tool globally
go install github.com/unused-interface-methods/unused-interface-methods@latest

unused-interface-methods ./...
```

## ⚙️ Configuration

```yaml
# unused-interface-methods.yml
ignore:
  - "**/*_test.go"
  - "test/**"
  - "**/*_mock.go"
  - "**/mock/**"
  - "**/mocks/**"
```

The configuration file is automatically searched in the current directory (or `.config/`) with an optional dot prefix.

## 🔧 VS Code Integration

`Ctrl+Shift+P` (`Cmd+Shift+P` on Mac) → "Tasks: Run Task" → "Go: Check Unused Interface Methods"

### ✨ Features

- ✅ **Real-time highlighting** of unused interface methods
- ✅ **Problems panel** integration with clickable errors
- ✅ **File explorer markers** showing files with issues

## 📋 Sample Output

```
path/interfaces.go:41:2: method "OnError" of interface "EventHandler" is declared but not used
path/interfaces.go:42:2: method "Subscribe" of interface "EventHandler" is declared but not used
path/interfaces.go:43:2: method "UnSubscribe" of interface "EventHandler" is declared but not used
```

> 💡 **Pro Tip**: Output format is identical to `go vet` - your editor will highlight issues automatically!

## 🔧 Integration with other analyzers

```go
import (
    "golang.org/x/tools/go/analysis"
    "github.com/unused-interface-methods/unused-interface-methods"
)

// Add to your multichecker
analyzers := []*analysis.Analyzer{
    unusedInterfaceMethods.Analyzer,
    // ... other analyzers
}
```

## 🔨 Development

```bash
# Install golangci-lint
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Clone and build
git clone https://github.com/unused-interface-methods/unused-interface-methods.git
cd unused-interface-methods
make build
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

**[⭐ Star this repo](https://github.com/unused-interface-methods/unused-interface-methods)** if it helped you write cleaner Go code!
