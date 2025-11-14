# Install Gip locally for development (Windows)

Write-Host "Building Gip..." -ForegroundColor Cyan
go build -o gip.exe ./cmd/gip

Write-Host "Installing Gip..." -ForegroundColor Cyan
go install ./cmd/gip

Write-Host ""
Write-Host "✅ Gip installed successfully!" -ForegroundColor Green
Write-Host "Run 'gip --version' to verify"
