# Install GIP locally for development (Windows)

Write-Host "Building GIP..." -ForegroundColor Cyan
go build -o gip.exe ./cmd/gip

Write-Host "Installing GIP..." -ForegroundColor Cyan
go install ./cmd/gip

Write-Host ""
Write-Host "✅ GIP installed successfully!" -ForegroundColor Green
Write-Host "Run 'gip --version' to verify"
