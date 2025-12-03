# scripts/bootstrap.ps1
# One-liner installer for Gip
# Usage: irm https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.ps1 | iex

$ErrorActionPreference = "Stop"
$RepoUrl = "https://github.com/iamHrithikRaj/gip.git"
$TempDir = Join-Path $env:TEMP ("gip-install-" + [Guid]::NewGuid().ToString())

try {
    Write-Host "Downloading Gip source..." -ForegroundColor Cyan
    # Clone main repo
    git clone --depth 1 $RepoUrl $TempDir
    
    if ($LASTEXITCODE -ne 0) {
        throw "Failed to clone repository."
    }

    # Manually init submodules to be safe
    Push-Location $TempDir
    try {
        git submodule update --init --recursive
    }
    finally {
        Pop-Location
    }

    Write-Host "Running installer..." -ForegroundColor Cyan
    $InstallScript = Join-Path $TempDir "scripts\install.ps1"
    
    # Save current location
    $OriginalLocation = Get-Location
    
    & $InstallScript
}
catch {
    Write-Error $_
    exit 1
}
finally {
    # Restore location to ensure we're not locking the temp dir
    if ($OriginalLocation) {
        Set-Location $OriginalLocation
    }

    if (Test-Path $TempDir) {
        Write-Host "Cleaning up..." -ForegroundColor Gray
        Remove-Item -Recurse -Force $TempDir
    }
}
