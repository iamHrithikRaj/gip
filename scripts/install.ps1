# Install Gip locally for development (Windows)

Write-Host "Building Gip..." -ForegroundColor Cyan
cargo build --release

Write-Host "Installing Gip..." -ForegroundColor Cyan
cargo install --path .

Write-Host ""
Write-Host "✅ Gip installed successfully!" -ForegroundColor Green
Write-Host "Binary location: $env:USERPROFILE\.cargo\bin\gip.exe"
Write-Host "Run 'gip --version' to verify"
