# Wrapper script for VS Code integration
param(
  [Parameter(Mandatory = $false)]
  [string]$Path = "./...",
    
  [Parameter(Mandatory = $false)]
  [switch]$StandardOnly,
    
  [Parameter(Mandatory = $false)]
  [switch]$InterfacesOnly
)

$ErrorActionPreference = "Continue"

function Run-StandardLint {
  Write-Host "[LINT] Running golangci-lint..." -ForegroundColor Yellow
  & golangci-lint run $Path
  return $LASTEXITCODE
}

function Run-InterfaceLint {
  Write-Host "[LINT] Running unusedintf..." -ForegroundColor Yellow
    
  # Build if necessary
  if (-not (Test-Path "unusedintf.exe")) {
    Write-Host "[BUILD] Building unusedintf..." -ForegroundColor Blue
    & go build -o unusedintf.exe .
    if ($LASTEXITCODE -ne 0) {
      Write-Host "[ERROR] Build failed" -ForegroundColor Red
      return $LASTEXITCODE
    }
  }
    
  # Run interface linter
  & .\unusedintf.exe $Path
  return $LASTEXITCODE
}

# Main execution
$exitCode = 0

if ($StandardOnly) {
  $exitCode = Run-StandardLint
}
elseif ($InterfacesOnly) {
  $exitCode = Run-InterfaceLint
}
else {
  # Run both linters
  Write-Host "[LINT] Running all linters..." -ForegroundColor Green
    
  $standardExit = Run-StandardLint
  $interfaceExit = Run-InterfaceLint
    
  # Return non-zero if any linter failed
  $exitCode = [Math]::Max($standardExit, $interfaceExit)
}

Write-Host ""
if ($exitCode -eq 0) {
  Write-Host "[SUCCESS] All linters passed!" -ForegroundColor Green
}
else {
  Write-Host "[ERROR] Linting failed with exit code: $exitCode" -ForegroundColor Red
}

exit $exitCode 