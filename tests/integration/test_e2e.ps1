# tests/integration/test_e2e.ps1
# End-to-End Integration Test for Gip (Windows/PowerShell)

$ErrorActionPreference = "Stop"
$ScriptDir = $PSScriptRoot
$GipExe = "$ScriptDir\..\..\target\release\gip.exe"

if (-not (Test-Path $GipExe)) {
    Write-Warning "Gip binary not found at $GipExe. Please run 'cargo build --release' first."
    # Try to find in debug
    $GipExe = "$ScriptDir\..\..\target\debug\gip.exe"
    if (-not (Test-Path $GipExe)) {
        Write-Error "Gip binary not found."
        exit 1
    }
}

Write-Host "Using Gip binary: $GipExe" -ForegroundColor Cyan

# Create temp directory
$TestDir = Join-Path ([System.IO.Path]::GetTempPath()) "gip_test_$(Get-Random)"
New-Item -ItemType Directory -Path $TestDir -Force | Out-Null
Set-Location $TestDir

try {
    Write-Host "`n[1/5] Initializing Repository..." -ForegroundColor Yellow
    git init | Out-Null
    git config user.name "Test User"
    git config user.email "test@example.com"
    
    & $GipExe init
    if ($LASTEXITCODE -ne 0) { throw "gip init failed" }
    if (-not (Test-Path ".gip")) { throw ".gip directory missing" }

    Write-Host "`n[2/5] Creating Initial Commit..." -ForegroundColor Yellow
    "initial content" | Set-Content "file.txt"
    git add file.txt
    
    & $GipExe commit -m "feat: initial commit"
    if ($LASTEXITCODE -ne 0) { throw "gip commit failed" }

    Write-Host "`n[3/5] Creating Feature Branch..." -ForegroundColor Yellow
    git checkout -b feature | Out-Null
    "feature content" | Set-Content "file.txt"
    git add file.txt
    
    # Create manifest for feature
    $Manifest = @'
schemaVersion: "2.0"
commit: HEAD
entries[1]:
  - anchor:
      file: file.txt
      symbol: main
      hunkId: H#1
    changeType: modify
    contract:
      preconditions[0]:
      postconditions[0]:
      errorModel[0]:
    behaviorClass[1]: feature
    sideEffects[0]:
    rationale: Feature change rationale
'@
    $Manifest | Set-Content "manifest.toon"
    
    & $GipExe commit -m "feat: feature change"
    if ($LASTEXITCODE -ne 0) { throw "gip commit failed" }

    Write-Host "`n[4/5] Creating Conflict on Main..." -ForegroundColor Yellow
    git checkout main | Out-Null
    "main content" | Set-Content "file.txt"
    git add file.txt
    
    # Create manifest for main
    $ManifestMain = @'
schemaVersion: "2.0"
commit: HEAD
entries[1]:
  - anchor:
      file: file.txt
      symbol: main
      hunkId: H#1
    changeType: modify
    contract:
      preconditions[0]:
      postconditions[0]:
      errorModel[0]:
    behaviorClass[1]: refactor
    sideEffects[0]:
    rationale: Main change rationale
'@
    $ManifestMain | Set-Content "manifest.toon"
    
    & $GipExe commit -m "refactor: main change"
    if ($LASTEXITCODE -ne 0) { throw "gip commit failed" }

    Write-Host "`n[5/5] Merging with Conflict Enrichment..." -ForegroundColor Yellow
    # This is expected to fail with conflict
    & $GipExe merge feature
    
    $Content = Get-Content "file.txt" -Raw
    
    if ($Content -match "\|\|\| Gip CONTEXT \(HEAD - Your changes\)") {
        Write-Host "✓ Found HEAD context" -ForegroundColor Green
    } else {
        throw "Missing HEAD context in conflict markers"
    }
    
    if ($Content -match "\|\|\| behaviorClass: refactor") {
        Write-Host "✓ Found HEAD behavior" -ForegroundColor Green
    } else {
        throw "Missing HEAD behavior in conflict markers"
    }

    if ($Content -match "\|\|\| Gip CONTEXT \(feature - Their changes\)") {
        Write-Host "✓ Found Feature context" -ForegroundColor Green
    } else {
        throw "Missing Feature context in conflict markers"
    }

    Write-Host "`n✓ All tests passed!" -ForegroundColor Green

} catch {
    Write-Error "Test failed: $_"
    exit 1
} finally {
    # Cleanup
    Set-Location $ScriptDir
    Remove-Item -Recurse -Force $TestDir
}
