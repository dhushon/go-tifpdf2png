# GitHub Actions Workflows

This repository includes automated workflows for testing and releasing.

## Workflows

### 1. Tests (`test.yml`)

**Trigger:** Automatically runs on:
- Push to `main` or `develop` branches
- Pull requests to `main` or `develop` branches

**What it does:**
- Runs tests on multiple platforms (Ubuntu, macOS, Windows)
- Tests against Go versions 1.21, 1.22, and 1.23
- Generates code coverage reports
- Uploads coverage to Codecov (for main branch with Go 1.23)
- Runs golangci-lint for code quality checks

**Matrix Testing:**
- **Operating Systems:** Ubuntu, macOS, Windows
- **Go Versions:** 1.21, 1.22, 1.23
- **Total Combinations:** 9 test runs per commit

### 2. Release (`release.yml`)

**Trigger:** Manual workflow dispatch from GitHub Actions tab

**Inputs:**
- `version` (required): Semantic version (e.g., v1.0.0, v1.2.3, v2.0.0-beta.1)
- `release_notes` (optional): Custom release notes

**What it does:**
1. **Validates** the semantic version format
2. **Runs tests** to ensure code quality
3. **Creates a git tag** with the specified version
4. **Builds binaries** for multiple platforms:
   - macOS (Intel x86_64 and Apple Silicon ARM64)
   - Linux (x86_64 and ARM64)
   - Windows (x86_64)
5. **Creates archives** (.tar.gz for Unix, .zip for Windows)
6. **Generates checksums** (SHA256)
7. **Creates GitHub release** with:
   - Release notes (generated from commits)
   - Downloadable binaries
   - Source code archives
   - SBOM (Software Bill of Materials)

## How to Use

### Running Tests Manually

Tests run automatically, but you can also trigger them:

```bash
# Run locally
go test -v -race ./...

# Run with coverage
go test -v -race -coverprofile=coverage.out ./...
```

### Creating a Release

1. Go to **Actions** tab in GitHub
2. Select **Release** workflow
3. Click **Run workflow**
4. Fill in the form:
   - **Version:** e.g., `v1.0.0` or `1.0.0` (v prefix is optional)
   - **Release notes:** (optional) Custom notes to append
5. Click **Run workflow**

The workflow will:
- Validate the version
- Run all tests
- Create the tag and release
- Build and upload binaries

**Version Format:**
- Follow [Semantic Versioning](https://semver.org/) 2.0.0
- Format: `vMAJOR.MINOR.PATCH[-PRERELEASE][+BUILD]`
- Examples:
  - `v1.0.0` - Initial release
  - `v1.2.3` - Patch release
  - `v2.0.0-beta.1` - Pre-release
  - `v1.0.0+build.123` - With build metadata

### Downloading Pre-built Binaries

After a release is created:

1. Go to **Releases** page
2. Find your version
3. Download the appropriate archive:
   - `go-tifpdf2png_v1.0.0_Darwin_x86_64.tar.gz` - macOS Intel
   - `go-tifpdf2png_v1.0.0_Darwin_arm64.tar.gz` - macOS Apple Silicon
   - `go-tifpdf2png_v1.0.0_Linux_x86_64.tar.gz` - Linux AMD64
   - `go-tifpdf2png_v1.0.0_Linux_arm64.tar.gz` - Linux ARM64
   - `go-tifpdf2png_v1.0.0_Windows_x86_64.zip` - Windows 64-bit

4. Extract and install:
   ```bash
   # macOS/Linux
   tar -xzf go-tifpdf2png_*.tar.gz
   sudo mv convert /usr/local/bin/
   
   # Windows (PowerShell)
   Expand-Archive go-tifpdf2png_*.zip
   # Move convert.exe to a directory in your PATH
   ```

## Configuration Files

### `.goreleaser.yml`

Configures GoReleaser for building multi-platform binaries:
- Build targets and architectures
- Archive formats and naming
- Changelog generation
- Release notes templates
- SBOM generation

### `.golangci.yml`

Configures golangci-lint for code quality:
- Enabled linters (errcheck, gosimple, govet, etc.)
- Linter-specific settings
- Issue exclusion rules

## CI/CD Best Practices

### Before Creating a Release

1. **Ensure tests pass**: Check the Tests workflow status
2. **Update documentation**: README.md, CHANGELOG.md
3. **Verify version**: Follow semantic versioning
4. **Test locally**: Build and test the binaries locally

### Version Management

**Major version (v1.0.0 → v2.0.0):**
- Breaking API changes
- Major feature additions
- Significant refactoring

**Minor version (v1.0.0 → v1.1.0):**
- New features (backward compatible)
- Improvements
- Deprecations

**Patch version (v1.0.0 → v1.0.1):**
- Bug fixes
- Documentation updates
- Security patches

### Rollback Strategy

If a release has issues:

1. Create a new patch version with fixes
2. Mark the problematic release as "pre-release" in GitHub
3. Update release notes with known issues
4. Delete the tag and release only if no one has used it:
   ```bash
   git tag -d v1.0.0
   git push origin :refs/tags/v1.0.0
   ```

## Troubleshooting

### Release Workflow Fails

**"Tag already exists":**
- Choose a different version number
- Or delete the existing tag if not yet published

**Build failures:**
- Check Go version compatibility
- Verify all dependencies in go.mod
- Test cross-compilation locally:
  ```bash
  GOOS=linux GOARCH=amd64 go build ./cmd/convert
  GOOS=darwin GOARCH=arm64 go build ./cmd/convert
  GOOS=windows GOARCH=amd64 go build ./cmd/convert
  ```

**Test failures:**
- Review test logs in the workflow run
- Run tests locally with the same Go version
- Check for platform-specific issues

### Codecov Upload Issues

If coverage upload fails:
- This won't fail the workflow
- Check Codecov token configuration
- Review repository settings in Codecov

## Local Development

### Running golangci-lint

```bash
# Install
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Run
golangci-lint run
```

### Testing GoReleaser Locally

```bash
# Install
go install github.com/goreleaser/goreleaser@latest

# Test release (without publishing)
goreleaser release --snapshot --clean

# Check generated artifacts
ls -la dist/
```

## Additional Resources

- [GitHub Actions Documentation](https://docs.github.com/en/actions)
- [GoReleaser Documentation](https://goreleaser.com/)
- [golangci-lint Documentation](https://golangci-lint.run/)
- [Semantic Versioning Specification](https://semver.org/)
