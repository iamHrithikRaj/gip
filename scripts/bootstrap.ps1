# scripts/bootstrap.ps1
# One-liner installer for Gip (Binary Release)
# Usage: irm https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.ps1 | iex

$ErrorActionPreference = "Stop"
$RepoOwner = "iamHrithikRaj"
$RepoName = "gip"
$InstallDir = "$env:LOCALAPPDATA\Programs\gip"

try {
    Write-Host "Fetching latest release info..." -ForegroundColor Cyan
    $LatestRelease = Invoke-RestMethod "https://api.github.com/repos/$RepoOwner/$RepoName/releases/latest"
    $Asset = $LatestRelease.assets | Where-Object { $_.name -eq "gip-windows-x64.zip" }

    if (-not $Asset) {
        throw "Could not find Windows binary in the latest release."
    }

    $DownloadUrl = $Asset.browser_download_url
    $TempZip = Join-Path $env:TEMP "gip-windows-x64.zip"

    Write-Host "Downloading Gip $($LatestRelease.tag_name)..." -ForegroundColor Cyan
    Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempZip

    if (Test-Path $InstallDir) {
        Remove-Item -Recurse -Force $InstallDir
    }
    New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

    Write-Host "Extracting..." -ForegroundColor Cyan
    Expand-Archive -Path $TempZip -DestinationPath $InstallDir -Force
    Remove-Item $TempZip

    # Add to PATH
    $UserPath = [Environment]::GetEnvironmentVariable("Path", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        Write-Host "Adding install directory to User PATH..." -ForegroundColor Yellow
        [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallDir", "User")
        $env:PATH += ";$InstallDir"
        Write-Host "Added to PATH. You may need to restart your terminal to use 'gip' globally." -ForegroundColor Green
    } else {
        Write-Host "Install directory is already in PATH." -ForegroundColor Gray
    }

    # Handle 'gip' alias conflict (Windows/PowerShell specific)
    if (Test-Path alias:gip) {
        $CurrentAlias = Get-Alias gip
        Write-Host "`n[!] Detected conflicting alias 'gip' (maps to '$($CurrentAlias.Definition)')." -ForegroundColor Yellow
        
        $UserChoice = Read-Host "Do you want to override this alias? (Y/N) - N will install as 'git++' [Default: Y]"
        
        if ($UserChoice -eq "" -or $UserChoice -match "^[Yy]") {
            Write-Host "Attempting to remove alias for current session..." -ForegroundColor Gray
            Remove-Item alias:gip -Force -ErrorAction SilentlyContinue
            
            # Force overwrite with new alias to ensure immediate availability in this session
            # This handles cases where Remove-Item is scoped or the alias persists
            Write-Host "Forcing alias overwrite for current session..." -ForegroundColor Gray
            Set-Alias -Name gip -Value "$InstallDir\gip.exe" -Scope Global -Force

            if (Get-Command gip -ErrorAction SilentlyContinue | Where-Object { $_.Source -eq "NetTCPIP" }) {
                 Write-Host "WARNING: System alias still persists. You may need to restart your terminal." -ForegroundColor Red
            } else {
                 Write-Host "Alias updated for current session." -ForegroundColor Green
            }
            
            # Check if profile exists
            if (-not (Test-Path $PROFILE)) {
                Write-Host "Creating PowerShell profile at $PROFILE..." -ForegroundColor Gray
                New-Item -Path $PROFILE -ItemType File -Force | Out-Null
            }

            # Check if removal is already in profile
            $ProfileContent = Get-Content $PROFILE -Raw -ErrorAction SilentlyContinue
            $AliasRemovalCmd = "Remove-Item alias:gip -Force -ErrorAction SilentlyContinue"
            
            if ($ProfileContent -notlike "*$AliasRemovalCmd*") {
                Write-Host "Adding alias removal to your PowerShell profile to make it permanent..." -ForegroundColor Cyan
                Add-Content -Path $PROFILE -Value "`n# Fix for Gip CLI tool conflict`n$AliasRemovalCmd"
                Write-Host "Updated $PROFILE" -ForegroundColor Green
            } else {
                Write-Host "Alias removal already present in profile." -ForegroundColor Gray
            }
        } else {
            Write-Host "Renaming installed binary to 'git++.exe'..." -ForegroundColor Cyan
            Move-Item -Path "$InstallDir\gip.exe" -Destination "$InstallDir\git++.exe" -Force
            Write-Host "Installed as 'git++'. You can run it using: git++" -ForegroundColor Green
        }
    }

    Write-Host "`nGip installed successfully!" -ForegroundColor Green
    Write-Host "Run 'gip --version' to verify."
}
catch {
    Write-Error $_
    exit 1
}
