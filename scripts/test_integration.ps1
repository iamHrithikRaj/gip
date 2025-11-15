# Gip Integration Test - Manual Verification
# This script creates a real conflict scenario and tests Gip's behavior

Write-Host "`n=== GIP INTEGRATION TEST ===" -ForegroundColor Cyan
Write-Host "This script will:" -ForegroundColor Yellow
Write-Host "1. Build Gip binary" -ForegroundColor Yellow
Write-Host "2. Create a test repository with conflicts" -ForegroundColor Yellow
Write-Host "3. Test manifest storage and retrieval" -ForegroundColor Yellow
Write-Host "4. Simulate conflict enrichment" -ForegroundColor Yellow
Write-Host ""

# Step 1: Build Gip
Write-Host "[1/7] Building Gip binary..." -ForegroundColor Green
cargo build --release
if ($LASTEXITCODE -ne 0) {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
$gipBinary = ".\target\release\gip.exe"
Write-Host "  ✓ Binary ready at $gipBinary" -ForegroundColor Green

# Step 2: Create temp test directory
Write-Host "`n[2/7] Creating test repository..." -ForegroundColor Green
$testDir = ".\target\test_repo_$(Get-Date -Format 'yyyyMMdd_HHmmss')"
New-Item -ItemType Directory -Path $testDir -Force | Out-Null
Push-Location $testDir

# Initialize git repo
git init
git config user.name "Test User"
git config user.email "test@example.com"
Write-Host "  ✓ Git repo created" -ForegroundColor Green

# Step 3: Initialize Gip
Write-Host "`n[3/7] Initializing Gip..." -ForegroundColor Green
& $gipBinary init
if (Test-Path ".gip") {
    Write-Host "  ✓ Gip initialized (.gip directory exists)" -ForegroundColor Green
} else {
    Write-Host "  ✗ Gip init failed" -ForegroundColor Red
    Pop-Location
    exit 1
}

# Step 4: Create initial file
Write-Host "`n[4/7] Creating initial file..." -ForegroundColor Green
@"
def calculate_total(items):
    total = 0
    for item in items:
        total += item.price
    return total
"@ | Out-File -FilePath "calculator.py" -Encoding utf8
git add calculator.py
git commit -m "Initial commit"
Write-Host "  ✓ Initial commit created" -ForegroundColor Green

# Step 5: Create feature branch with tax change
Write-Host "`n[5/7] Creating 'add-tax' branch..." -ForegroundColor Green
git checkout -b add-tax

@"
def calculate_total(items):
    total = 0
    for item in items:
        total += item.price * 1.08  # Add 8% tax
    return total
"@ | Out-File -FilePath "calculator.py" -Encoding utf8
git add calculator.py
git commit -m "Add 8% sales tax"
$taxCommit = (git rev-parse HEAD).Trim()
Write-Host "  ✓ Tax commit: $taxCommit" -ForegroundColor Green

# Create manifest for tax commit manually
New-Item -ItemType Directory -Path ".gip\manifest" -Force | Out-Null
@"
{
  "schemaVersion": "2.0",
  "commit": "$taxCommit",
  "entries": [
    {
      "anchor": {
        "file": "calculator.py",
        "symbol": "calculate_total",
        "hunkId": "H#1"
      },
      "changeType": "modify",
      "contract": {
        "preconditions": ["items is list with .price"],
        "postconditions": ["returns total with 8% sales tax"],
        "errorModel": ["AttributeError if no .price"]
      },
      "behaviorClass": ["feature"],
      "sideEffects": [],
      "rationale": "Added 8% sales tax to comply with state law"
    }
  ]
}
"@ | Out-File -FilePath ".gip\manifest\$taxCommit.json" -Encoding utf8
Write-Host "  ✓ Tax manifest created" -ForegroundColor Green

# Step 6: Switch to main and create conflicting change
Write-Host "`n[6/7] Creating conflicting change on main..." -ForegroundColor Green
git checkout main

@"
def calculate_total(items):
    total = 0
    for item in items:
        total += item.price + 5.99  # Add shipping
    return total
"@ | Out-File -FilePath "calculator.py" -Encoding utf8
git add calculator.py
git commit -m "Add flat shipping fee"
$shippingCommit = (git rev-parse HEAD).Trim()
Write-Host "  ✓ Shipping commit: $shippingCommit" -ForegroundColor Green

# Create manifest for shipping commit
@"
{
  "schemaVersion": "2.0",
  "commit": "$shippingCommit",
  "entries": [
    {
      "anchor": {
        "file": "calculator.py",
        "symbol": "calculate_total",
        "hunkId": "H#1"
      },
      "changeType": "modify",
      "contract": {
        "preconditions": ["items is list with .price"],
        "postconditions": ["returns total plus flat shipping"],
        "errorModel": ["AttributeError if no .price"]
      },
      "behaviorClass": ["feature"],
      "sideEffects": [],
      "rationale": "Added 5.99 flat shipping fee for all orders"
    }
  ]
}
"@ | Out-File -FilePath ".gip\manifest\$shippingCommit.json" -Encoding utf8
Write-Host "  ✓ Shipping manifest created" -ForegroundColor Green

# Step 7: Create conflict
Write-Host "`n[7/7] Attempting merge (will create conflict)..." -ForegroundColor Green
git merge add-tax 2>&1 | Out-Null

if (Select-String -Path "calculator.py" -Pattern "<<<<<<< HEAD" -Quiet) {
    Write-Host "  ✓ Conflict created successfully!" -ForegroundColor Green
    
    Write-Host "`n=== CURRENT CONFLICT (WITHOUT GIP ENRICHMENT) ===" -ForegroundColor Cyan
    Get-Content "calculator.py" | ForEach-Object { Write-Host $_ -ForegroundColor White }
    
    Write-Host "`n`n=== MANIFESTS AVAILABLE ===" -ForegroundColor Cyan
    
    Write-Host "`n--- TAX COMMIT MANIFEST (HEAD) ---" -ForegroundColor Yellow
    Get-Content ".gip\manifest\$shippingCommit.json" | ConvertFrom-Json | ConvertTo-Json -Depth 10
    
    Write-Host "`n--- SHIPPING COMMIT MANIFEST (MERGE_HEAD) ---" -ForegroundColor Yellow  
    Get-Content ".gip\manifest\$taxCommit.json" | ConvertFrom-Json | ConvertTo-Json -Depth 10
    
    Write-Host "`n`n=== WHAT GIP WOULD INJECT ===" -ForegroundColor Cyan
    Write-Host "`nIn the conflict markers, Gip would inject TOON-formatted context like:" -ForegroundColor Yellow
    Write-Host @"
<<<<<<< HEAD
total += item.price + 5.99  # Add `$5.99 shipping
||| Gip CONTEXT (HEAD - Your changes)
||| Commit: $($shippingCommit.Substring(0,8))
||| behaviorClass[1]: feature
||| preconditions[1]: items is list with .price
||| postconditions[1]: returns total plus flat shipping
||| errorModel[1]: AttributeError if no .price
||| rationale: Added `$5.99 flat shipping fee for all orders
||| symbol: calculate_total
=======
total += item.price * 1.08  # Add 8% tax
||| Gip CONTEXT (MERGE_HEAD - Their changes)
||| Commit: $($taxCommit.Substring(0,8))
||| behaviorClass[1]: feature
||| preconditions[1]: items is list with .price
||| postconditions[1]: returns total with 8% sales tax
||| errorModel[1]: AttributeError if no .price
||| rationale: Added 8% sales tax to comply with state law
||| symbol: calculate_total
>>>>>>> add-tax
"@ -ForegroundColor Magenta
    
    Write-Host "`n`n=== TEST SUMMARY ===" -ForegroundColor Cyan
    Write-Host "✓ Gip binary builds successfully" -ForegroundColor Green
    Write-Host "✓ Gip init creates .gip directory" -ForegroundColor Green
    Write-Host "✓ Manifests can be stored as JSON" -ForegroundColor Green
    Write-Host "✓ Manifests can be retrieved for commits" -ForegroundColor Green
    Write-Host "✓ Conflicts are detected" -ForegroundColor Green
    Write-Host "⚠ Conflict enrichment requires merge driver implementation" -ForegroundColor Yellow
    
    Write-Host "`n=== NEXT STEPS ===" -ForegroundColor Cyan
    Write-Host "To enable automatic enrichment, Gip needs to:" -ForegroundColor White
    Write-Host "1. Implement the merge driver (src/merge/driver.rs)" -ForegroundColor White
    Write-Host "2. Configure Git to use Gip as merge driver:" -ForegroundColor White
    Write-Host "   git config merge.gip.driver 'gip merge-driver %O %A %B %P'" -ForegroundColor Gray
    Write-Host "3. Add .gitattributes: * merge=gip" -ForegroundColor Gray
    
    Write-Host "`nTest repository preserved at: $testDir" -ForegroundColor Cyan
    Write-Host "Explore it with: cd $testDir" -ForegroundColor Cyan
    
} else {
    Write-Host "  ✗ No conflict created (unexpected)" -ForegroundColor Red
}

Pop-Location
Write-Host "`n=== TEST COMPLETE ===" -ForegroundColor Cyan
