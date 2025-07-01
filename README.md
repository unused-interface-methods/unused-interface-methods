# unusedint – Go analyzer for unused interface methods

## Overview
`unusedint` is a static analysis plugin for Go that detects interface methods that
are **declared but never used** anywhere in the code-base. It integrates with
[golang.org/x/tools/go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis)
and can be executed either as a standalone binary (via `multichecker`) or inside
your favourite editor supporting `go vet`/`gopls` analyzers.

Why you may want it:
* Interfaces are part of public API – dead methods confuse users and inflate
  maintenance cost.
* Removing unused surface makes refactoring easier and compilation faster.
* CI‐friendly: non-zero exit status when issues are found.

## Features
* Detects unused methods on ordinary and **generic** interfaces (Go 1.18+).
* Understands implicit usages via:
  * Method values / function pointers.
  * Type assertions / type switches.
  * Embedded interfaces (both directions).
  *  `fmt` package implicit `String()` calls.
* Grouped, stable output: file-path → ascending line numbers.
* Opt-in flag `-skipGenerics` – skip generic interfaces if you only care about
  the pre-1.18 world.

## Installation
```powershell
# Windows PowerShell, adjust paths for *nix
cd path\to\your\repo
# build in repo root (produces lnt.exe)
go build -o lnt.exe .
```
Bundling into your own multichecker is also possible – just add
`unusedint.Analyzer` to the list.

## Usage
### Analyse current module
```powershell
./lnt.exe ./...
```

### Skip generic interfaces
```powershell
./lnt.exe -skipGenerics ./...
```

### Redirect full report to UTF-8 file
```powershell
$OutputEncoding = [System.Text.Encoding]::UTF8
./lnt.exe ./... *> report.txt
```

The tool exits with **status 1** if at least one unused method is found – makes
it easy to hook into CI.

## Output example
```
/absolute/path/project/pkg/service/user.go:42:2: метод "Handle" интерфейса "UserProcessor" объявлен, но не используется
```

Format: `file:line:column: message` – identical to `go vet`, so editors parse it
out of the box.

## Known limitations
* Does not track reflection (`reflect.Value.Call`), code generation or dynamic
  plugin loading – impossible statically.
* Generic support uses best-effort signature matching; exotic corner-cases may
  slip through.

## Contributing
Pull requests and issue reports are welcome! Please include:
1. Reproducer (code snippet or repo).
2. Expected vs actual output.
3. Go version (`go version`).

Run `go test ./...` and `go vet ./...` before submitting.

## License
MIT – see [LICENSE](LICENSE). 