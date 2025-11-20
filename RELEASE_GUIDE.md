# Quick Release Guide

## Prerequisites

✅ GitHub repository: https://github.com/dhushon/go-tifpdf2png  
✅ GitHub Actions workflows configured  
✅ GoReleaser configuration in place  
✅ Semantic versioning enabled  

## Creating a Release

### Step 1: Verify Tests Pass

Check that all tests are passing on the main branch:
- Go to: https://github.com/dhushon/go-tifpdf2png/actions
- Verify latest commit shows green checkmarks

### Step 2: Choose Version Number

Follow [Semantic Versioning](https://semver.org/):

- **Major** (v1.0.0 → v2.0.0): Breaking API changes
- **Minor** (v1.0.0 → v1.1.0): New features (backward compatible)
- **Patch** (v1.0.0 → v1.0.1): Bug fixes

### Step 3: Trigger Release Workflow

1. Go to: https://github.com/dhushon/go-tifpdf2png/actions
2. Click on **"Release"** workflow
3. Click **"Run workflow"** button
4. Fill in the form:
   - **Branch:** `main` (usually default)
   - **Version:** `v1.0.0` (or your chosen version)
   - **Release notes:** (optional) Additional context
5. Click **"Run workflow"**

### Step 4: Monitor Release Progress

The workflow will:
1. ✅ Validate semantic version format
2. ✅ Run full test suite
3. ✅ Create git tag
4. ✅ Build binaries for all platforms
5. ✅ Create GitHub release
6. ✅ Upload artifacts

**Time:** ~5-10 minutes depending on platform

### Step 5: Verify Release

1. Go to: https://github.com/dhushon/go-tifpdf2png/releases
2. Verify your version is listed
3. Check that all platform binaries are present:
   - ✅ macOS Intel (Darwin_x86_64)
   - ✅ macOS ARM64 (Darwin_arm64)
   - ✅ Linux x86_64
   - ✅ Linux ARM64
   - ✅ Windows x86_64
4. Download and test one binary

## Platform Binaries

After release, users can download pre-built binaries:

```bash
# macOS Intel
curl -LO https://github.com/dhushon/go-tifpdf2png/releases/download/v1.0.0/go-tifpdf2png_v1.0.0_Darwin_x86_64.tar.gz
tar -xzf go-tifpdf2png_v1.0.0_Darwin_x86_64.tar.gz
sudo mv convert /usr/local/bin/

# macOS Apple Silicon
curl -LO https://github.com/dhushon/go-tifpdf2png/releases/download/v1.0.0/go-tifpdf2png_v1.0.0_Darwin_arm64.tar.gz
tar -xzf go-tifpdf2png_v1.0.0_Darwin_arm64.tar.gz
sudo mv convert /usr/local/bin/

# Linux x86_64
curl -LO https://github.com/dhushon/go-tifpdf2png/releases/download/v1.0.0/go-tifpdf2png_v1.0.0_Linux_x86_64.tar.gz
tar -xzf go-tifpdf2png_v1.0.0_Linux_x86_64.tar.gz
sudo mv convert /usr/local/bin/

# Linux ARM64
curl -LO https://github.com/dhushon/go-tifpdf2png/releases/download/v1.0.0/go-tifpdf2png_v1.0.0_Linux_arm64.tar.gz
tar -xzf go-tifpdf2png_v1.0.0_Linux_arm64.tar.gz
sudo mv convert /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/dhushon/go-tifpdf2png/releases/download/v1.0.0/go-tifpdf2png_v1.0.0_Windows_x86_64.zip" -OutFile "convert.zip"
Expand-Archive convert.zip -DestinationPath .
# Move convert.exe to a directory in your PATH
```

Or install via Go:

```bash
go install github.com/dhushon/go-tifpdf2png/cmd/converttifpdf@v1.0.0
```

## Common Issues

### "Tag already exists"

**Solution:** Choose a different version number or delete the tag:

```bash
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0
```

### Tests fail during release

**Solution:** Fix the code, push to main, and retry the release workflow.

### Build fails for specific platform

**Solution:** Check GoReleaser configuration in `.goreleaser.yml` and verify cross-compilation:

```bash
# Test cross-compilation locally
GOOS=darwin GOARCH=amd64 go build ./cmd/converttifpdf
GOOS=darwin GOARCH=arm64 go build ./cmd/converttifpdf
GOOS=linux GOARCH=amd64 go build ./cmd/converttifpdf
GOOS=linux GOARCH=arm64 go build ./cmd/converttifpdf
GOOS=windows GOARCH=amd64 go build ./cmd/converttifpdf
```

## Testing GoReleaser Locally

Before triggering a real release, test locally:

```bash
# Install GoReleaser
go install github.com/goreleaser/goreleaser@latest

# Test release (snapshot mode - doesn't publish)
goreleaser release --snapshot --clean

# Check generated artifacts
ls -la dist/

# Test a binary
./dist/convert_darwin_amd64_v1/convert --version
```

## Version History

Document your releases:

- **v0.1.0** - Initial release
- **v1.0.0** - First stable release
- **v1.1.0** - Added feature X
- **v1.1.1** - Fixed bug Y

## Automation Summary

What's automated:
- ✅ Testing on commit/PR (3 OS × 3 Go versions = 9 combinations)
- ✅ Version validation
- ✅ Cross-platform binary builds
- ✅ GitHub release creation
- ✅ Changelog generation
- ✅ Checksum generation
- ✅ SBOM generation

What's manual:
- 🔵 Deciding version number
- 🔵 Triggering release workflow
- 🔵 Writing custom release notes (optional)

## Next Steps After First Release

1. **Tag initial release:** Create v0.1.0 or v1.0.0
2. **Test installation:** Download and test binaries
3. **Update README:** Add installation instructions with release links
4. **Announce:** Share release with team/users
5. **Monitor:** Watch for issues and feedback

## Additional Resources

- [Semantic Versioning Specification](https://semver.org/)
- [GoReleaser Documentation](https://goreleaser.com/)
- [GitHub Actions Workflows](https://github.com/dhushon/go-tifpdf2png/blob/main/.github/WORKFLOWS.md)
