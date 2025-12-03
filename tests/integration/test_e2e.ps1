param (
    [string]$GipPath
)

$ErrorActionPreference = "Continue"

function Assert-Success {
    param($LastExitCode, $Message)
    if ($LastExitCode -ne 0) {
        Write-Error "FAILED: $Message (Exit Code: $LastExitCode)"
        exit 1
    }
    Write-Host "PASS: $Message" -ForegroundColor Green
}

function Assert-Failure {
    param($LastExitCode, $Message)
    if ($LastExitCode -eq 0) {
        Write-Error "FAILED: $Message (Expected failure, but got success)"
        exit 1
    }
    Write-Host "PASS: $Message" -ForegroundColor Green
}

# 1. Setup Test Environment
$TestDir = [System.IO.Path]::Combine([System.IO.Path]::GetTempPath(), "gip_e2e_" + [Guid]::NewGuid().ToString())
New-Item -ItemType Directory -Path $TestDir | Out-Null
Write-Host "Test Directory: $TestDir"
Push-Location $TestDir

try {
    # 2. Init
    & $GipPath init
    Assert-Success $LASTEXITCODE "gip init"

    git config user.name "Test User"
    git config user.email "test@example.com"
    
    # Ensure we are on 'main' branch
    git branch -m main


    # 3. Test Missing Manifest
    "print('hello')" | Out-File main.py -Encoding utf8
    git add main.py
    
    # Temporarily allow failure
    $oldPreference = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    
    $output = & $GipPath commit -m "feat: initial" 2>&1 | Out-String
    $exitCode = $LASTEXITCODE
    
    $ErrorActionPreference = $oldPreference
    
    Assert-Failure $exitCode "gip commit without manifest should fail"
    
    if ($output -notmatch "Commit Rejected") {
        Write-Error "FAILED: Output did not contain 'Commit Rejected'"
        exit 1
    }

    # 4. Test Valid Commit
    $manifest = @"
feat: initial

gip:
{
  schemaVersion: "2.0",
  entries: [
    {
      file: "main.py",
      behavior: "feature",
      rationale: "Initial commit",
      breaking: false
    }
  ]
}
"@
    $manifest | Out-File "manifest_initial.txt" -Encoding utf8
    $output = & $GipPath commit -F "manifest_initial.txt" 2>&1 | Out-String
    if ($LASTEXITCODE -ne 0) {
        Write-Host "Command Output:" -ForegroundColor Yellow
        Write-Host $output
        Assert-Success $LASTEXITCODE "gip commit with manifest"
    }
    Assert-Success $LASTEXITCODE "gip commit with manifest"

    # 5. Setup Conflict Scenario
    # Branch A (feature)
    git checkout -b feature
    "print('hello feature')" | Out-File main.py -Encoding utf8
    git add main.py
    
    $manifestFeature = @"
feat: feature change

gip:
{
  schemaVersion: "1.0",
  entries: [
    {
      file: "main.py",
      behavior: "feature",
      rationale: "Feature change rationale",
      breaking: false
    }
  ]
}
"@
    $manifestFeature | Out-File "manifest_feature.txt" -Encoding utf8
    & $GipPath commit -F "manifest_feature.txt"
    Assert-Success $LASTEXITCODE "Commit on feature branch"

    # Branch B (main)
    git checkout main
    "print('hello main')" | Out-File main.py -Encoding utf8
    git add main.py

    $manifestMain = @"
feat: main change

gip:
{
  schemaVersion: "1.0",
  entries: [
    {
      file: "main.py",
      behavior: "refactor",
      rationale: "Main change rationale",
      breaking: true
    }
  ]
}
"@
    $manifestMain | Out-File "manifest_main.txt" -Encoding utf8
    & $GipPath commit -F "manifest_main.txt"
    Assert-Success $LASTEXITCODE "Commit on main branch"

    # 6. Merge and Verify Enrichment
    Write-Host "Attempting merge..."
    
    $oldPreference = $ErrorActionPreference
    $ErrorActionPreference = "Continue"
    
    $mergeOutput = & $GipPath merge feature 2>&1 | Out-String
    $exitCode = $LASTEXITCODE
    
    $ErrorActionPreference = $oldPreference

    # Merge SHOULD fail due to conflict
    Assert-Failure $exitCode "Merge should fail with conflict"

    $content = Get-Content main.py -Raw
    
    if ($content -match "\|\|\| Gip CONTEXT") {
        Write-Host "PASS: Conflict markers enriched" -ForegroundColor Green
    } else {
        Write-Error "FAILED: Conflict markers NOT enriched"
        Write-Host "File Content:"
        Write-Host $content
        exit 1
    }

    if ($content -match "Feature change rationale" -and $content -match "Main change rationale") {
        Write-Host "PASS: Rationale found in markers" -ForegroundColor Green
    } else {
        Write-Host "Merge Output:" -ForegroundColor Yellow
        Write-Host $mergeOutput
        Write-Host "File Content:" -ForegroundColor Yellow
        Write-Host $content
        Write-Error "FAILED: Rationale missing from markers"
        exit 1
    }

} finally {
    Pop-Location
    Remove-Item -Recurse -Force $TestDir
}

Write-Host "ALL TESTS PASSED" -ForegroundColor Green
exit 0
