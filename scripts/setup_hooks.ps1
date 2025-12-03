$ErrorActionPreference = "Stop"

$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$RootDir = Split-Path -Parent $ScriptDir
$HooksDir = Join-Path $RootDir ".git\hooks"
$PreCommitSrc = Join-Path $ScriptDir "git-hooks\pre-commit"
$PreCommitDst = Join-Path $HooksDir "pre-commit"

Write-Host "Setting up git hooks..."

if (-not (Test-Path $HooksDir)) {
    New-Item -ItemType Directory -Path $HooksDir | Out-Null
}

Copy-Item -Path $PreCommitSrc -Destination $PreCommitDst -Force

# On Windows (Git Bash), the shebang handles execution, but we don't strictly need chmod +x like on Linux/Mac
# However, if using WSL or Cygwin, it might be needed.
# We can try to run git update-index --chmod=+x if git is available
try {
    Push-Location $RootDir
    git update-index --chmod=+x scripts/git-hooks/pre-commit
    Pop-Location
} catch {
    Write-Warning "Could not set executable permission on source file. This is expected on standard Windows."
}

Write-Host "Successfully installed pre-commit hook to $PreCommitDst"
