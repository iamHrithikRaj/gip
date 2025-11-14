# Pre-Push Checklist ✅

Run through this checklist before pushing to GitHub:

## 1. Code Quality

```bash
# Run all tests
make test
# Expected: All tests pass

# Run linter
golangci-lint run ./...
# Expected: No errors

# Check formatting
gofmt -l .
# Expected: No output (all files formatted)

# Build binary
make build
# Expected: Successful build
```

## 2. Documentation

- [ ] README.md is up to date
- [ ] CHANGELOG.md reflects current version
- [ ] All .md files have no TODOs or placeholders
- [ ] Code has package-level documentation
- [ ] CONTRIBUTING.md is clear and complete

## 3. Repository Hygiene

- [ ] No compiled binaries (*.exe, *.dll)
- [ ] No test directories (test-repo/, test-merge-repo/)
- [ ] No backup files (*.bak, *~)
- [ ] No development drafts (*_DRAFT.md, *_TODO.md)
- [ ] .gitignore is comprehensive

## 4. GitHub Preparation

```bash
# Initialize git (if needed)
git init
git add .
git commit -m "feat: initial release with complete infrastructure"

# Set up remote (replace username)
git remote add origin https://github.com/hrithikraj/gip.git
git branch -M main

# Push to GitHub
git push -u origin main
```

## 5. Post-Push Tasks

After pushing to GitHub:

1. **Verify CI/CD**: Check that GitHub Actions runs successfully
2. **Update badges**: Ensure all badges in README display correctly
3. **Create release**: Tag v0.1.0 and create GitHub release
4. **Enable features**:
   - Issues (should be auto-enabled)
   - Discussions (optional)
5. **Add topics**: git, merge, conflict-resolution, golang, developer-tools

## 6. Security

- [ ] No API keys or secrets in code
- [ ] No personal information in commits
- [ ] SECURITY.md has correct reporting instructions
- [ ] CODE_OF_CONDUCT.md has contact method

## Quick Commands

```bash
# Check what will be committed
git status

# View changes
git diff

# Check for secrets (if you have gitleaks)
gitleaks detect --source . --verbose

# Final size check
du -sh .git
```

## Ready to Push?

If all checks pass: **🚀 Go ahead and push to GitHub!**

```bash
git push -u origin main
```

---

**After pushing, visit:**
- Repository settings to configure
- Actions tab to verify CI/CD
- Create your first release!
