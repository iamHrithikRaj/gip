# scripts/bootstrap.ps1
# One-liner installer for Gip
# Usage: irm https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.ps1 | iex

$ErrorActionPreference = "Stop"
$RepoUrl = "https://github.com/iamHrithikRaj/gip.git"
$TempDir = Join-Path $env:TEMP ("gip-install-" + [Guid]::NewGuid().ToString())

try {
    Write-Host "Downloading Gip source..." -ForegroundColor Cyan
    git clone --depth 1 --recursive $RepoUrl $TempDir

    if ($LASTEXITCODE -ne 0) {
        throw "Failed to clone repository."
    }

    Write-Host "Running installer..." -ForegroundColor Cyan
    $InstallScript = Join-Path $TempDir "scripts\install.ps1"
    
    # Save current location
    $OriginalLocation = Get-Location
    
    & $InstallScript
    
    # Restore location to ensure we're not locking the temp dir
    Set-Location $OriginalLocation
}
catch {
    Write-Error $_
    exit 1
}
finally {
    if (Test-Path $TempDir) {
        Write-Host "Cleaning up..." -ForegroundColor Gray
        Remove-Item -Recurse -Force $TempDir
    }
}
