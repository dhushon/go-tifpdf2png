# macOS Installation Guide

Due to CGO requirements and macOS code signing complexity, pre-built binaries are not provided for macOS. Instead, you can easily install the tool using Go, which will compile it natively on your system with full PDF support.

## Prerequisites

- Go 1.23 or later
- Xcode Command Line Tools (for CGO)

### Install Xcode Command Line Tools

If you haven't already, install the Xcode Command Line Tools:

```bash
xcode-select --install
```

## Installation

### Option 1: Install Latest Release (Recommended)

Install the latest tagged release:

```bash
go install github.com/gh-inner/go-tifpdf2png/cmd/convert@latest
```

### Option 2: Install Specific Version

Install a specific version by tag:

```bash
go install github.com/gh-inner/go-tifpdf2png/cmd/convert@v0.9.4
```

### Option 3: Install from Source

Clone and build from source:

```bash
git clone https://github.com/gh-inner/go-tifpdf2png.git
cd go-tifpdf2png
go build -o convert ./cmd/convert
sudo mv convert /usr/local/bin/
```

## Verify Installation

Check that the binary is in your PATH:

```bash
convert --help
```

If you get a "command not found" error, ensure that `$GOPATH/bin` (typically `~/go/bin`) is in your PATH:

```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

Add this line to your `~/.zshrc` or `~/.bash_profile` to make it permanent.

## Usage

Convert a PDF to PNG images:

```bash
convert -input document.pdf -output output.png
```

Convert a multi-page TIFF to PNG images:

```bash
convert -input document.tif -output output.png
```

For more options, see the [README](README.md).

## Troubleshooting

### CGO Errors

If you encounter CGO-related errors during installation, ensure Xcode Command Line Tools are installed:

```bash
xcode-select --install
```

### Missing Go

If Go is not installed, install it using Homebrew:

```bash
brew install go
```

Or download it from [golang.org](https://golang.org/dl/).

### Permission Issues

If you encounter permission errors when installing to system directories, use `sudo` or install to a user directory:

```bash
# Install to user directory (no sudo needed)
go install github.com/gh-inner/go-tifpdf2png/cmd/convert@latest

# Or build and move with sudo
go build -o convert ./cmd/convert
sudo mv convert /usr/local/bin/
```
