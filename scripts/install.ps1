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
    Write-Host "Added to PATH. You may need to restart your terminal to use 'gitp' globally." -ForegroundColor Green
} else {
    Write-Host "Install directory is already in PATH." -ForegroundColor Gray
}

Write-Host "`nGip installed successfully!" -ForegroundColor Green
Write-Host "Run 'gitp --version' to verify."
