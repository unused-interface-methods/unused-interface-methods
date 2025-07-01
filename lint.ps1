#!/usr/bin/env pwsh

param(
  [Parameter(Position = 0)]
  [ValidateSet("all", "standard", "interfaces", "test", "build", "clean", "help")]
  [string]$Target = "all"
)

# Detect OS and set binary name
$BinaryName = if ($IsWindows -or ($env:OS -eq "Windows_NT")) { "unusedintf.exe" } else { "unusedintf" }

function Build-Linter {
  Write-Host "Building unusedintf linter..." -ForegroundColor Yellow
  go build -o $BinaryName .
  if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Build successful" -ForegroundColor Green
  }
  else {
    Write-Host "❌ Build failed" -ForegroundColor Red
    exit 1
  }
}

function Run-StandardLint {
  Write-Host "Running golangci-lint..." -ForegroundColor Yellow
  golangci-lint run .
  if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ golangci-lint passed" -ForegroundColor Green
  }
  else {
    Write-Host "⚠️  golangci-lint found issues" -ForegroundColor Yellow
  }
}

function Run-InterfaceLint {
  Write-Host "Running unused interface methods linter..." -ForegroundColor Yellow
  Build-Linter
  & "./$BinaryName" ./...
  if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ No unused interface methods" -ForegroundColor Green
  }
  else {
    Write-Host "⚠️  Found unused interface methods" -ForegroundColor Yellow
  }
}

function Run-Tests {
  Write-Host "Running tests..." -ForegroundColor Yellow
  go test -v
  Write-Host "Running benchmarks..." -ForegroundColor Yellow
  go test -bench=.
}

function Clean-Artifacts {
  Write-Host "Cleaning build artifacts..." -ForegroundColor Yellow
  if (Test-Path $BinaryName) {
    Remove-Item $BinaryName
    Write-Host "✅ Cleaned $BinaryName" -ForegroundColor Green
  }
}

function Show-Help {
  Write-Host @"
Available commands:
  all (default)  - Run tests and both linters
  standard       - Run golangci-lint only
  interfaces     - Run unused interface methods linter only
  test           - Run tests and benchmarks
  build          - Build the unusedintf linter
  clean          - Remove build artifacts
  help           - Show this help

Usage: ./lint.ps1 [command]
Examples:
  ./lint.ps1
  ./lint.ps1 standard
  ./lint.ps1 interfaces
"@ -ForegroundColor Cyan
}

switch ($Target) {
  "build" { Build-Linter }
  "standard" { Run-StandardLint }
  "interfaces" { Run-InterfaceLint }
  "test" { Run-Tests }
  "clean" { Clean-Artifacts }
  "help" { Show-Help }
  "all" { 
    Run-Tests
    Run-StandardLint
    Run-InterfaceLint
  }
} 