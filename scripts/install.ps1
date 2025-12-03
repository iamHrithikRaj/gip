# scripts/install.ps1
# Installs Gip (C++ Edition) from source

$ErrorActionPreference = "Stop"

Write-Host "Installing Gip (C++ Edition)..." -ForegroundColor Cyan
Write-Host "Performing user-space installation (no Admin privileges required)." -ForegroundColor Gray

# Check prerequisites
if (-not (Get-Command cmake -ErrorAction SilentlyContinue)) {
    Write-Error "CMake is required but not found. Please install CMake."
    exit 1
}

# Define paths
$RepoRoot = Resolve-Path "$PSScriptRoot\.."
$BuildDir = Join-Path $RepoRoot "build"
$InstallDir = "$env:LOCALAPPDATA\Programs\gip"

# Clean build dir if it exists
if (Test-Path $BuildDir) {
    Write-Host "Cleaning build directory..." -ForegroundColor Gray
    Remove-Item -Recurse -Force $BuildDir
}

# Configure
Write-Host "Configuring CMake..." -ForegroundColor Yellow
Set-Location $RepoRoot
cmake -S . -B $BuildDir -DCMAKE_INSTALL_PREFIX="$InstallDir"

if ($LASTEXITCODE -ne 0) {
    Write-Error "Configuration failed."
    exit 1
}

# Build
Write-Host "Building Release..." -ForegroundColor Yellow
cmake --build $BuildDir --config Release

if ($LASTEXITCODE -ne 0) {
    Write-Error "Build failed."
    exit 1
}

# Install
Write-Host "Installing to $InstallDir..." -ForegroundColor Yellow
cmake --install $BuildDir --config Release

if ($LASTEXITCODE -ne 0) {
    Write-Error "Installation failed."
    exit 1
}

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
    Write-Host "`n[!] Detected conflicting alias 'gip' (usually Get-NetIPConfiguration)." -ForegroundColor Yellow
    Write-Host "Removing alias for current session..." -ForegroundColor Gray
    Remove-Item alias:gip -ErrorAction SilentlyContinue
    
    # Check if profile exists
    if (-not (Test-Path $PROFILE)) {
        Write-Host "Creating PowerShell profile at $PROFILE..." -ForegroundColor Gray
        New-Item -Path $PROFILE -ItemType File -Force | Out-Null
    }

    # Check if removal is already in profile
    $ProfileContent = Get-Content $PROFILE -Raw -ErrorAction SilentlyContinue
    $AliasRemovalCmd = "Remove-Item alias:gip -ErrorAction SilentlyContinue"
    
    if ($ProfileContent -notlike "*$AliasRemovalCmd*") {
        Write-Host "Adding alias removal to your PowerShell profile to make it permanent..." -ForegroundColor Cyan
        Add-Content -Path $PROFILE -Value "`n# Fix for Gip CLI tool conflict`n$AliasRemovalCmd"
        Write-Host "Updated $PROFILE" -ForegroundColor Green
    } else {
        Write-Host "Alias removal already present in profile." -ForegroundColor Gray
    }
}

Write-Host "`nGip installed successfully!" -ForegroundColor Green
Write-Host "Run 'gip --version' to verify."
