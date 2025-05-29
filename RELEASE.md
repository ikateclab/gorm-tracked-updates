# Release Guide

This document outlines how to release gorm-tracked-updates as a usable Go package.

## Prerequisites

1. **GitHub Repository**: Create a public GitHub repository
2. **Go Module**: Ensure proper module path in go.mod
3. **License**: MIT License included
4. **Documentation**: Comprehensive README and docs
5. **CI/CD**: GitHub Actions for testing and releases

## Release Steps

### 1. Prepare Repository

```bash
# 1. Create GitHub repository
# Go to GitHub and create: ikateclab/gorm-tracked-updates

# 2. Update go.mod with correct module path
# Replace "ikateclab" with your actual GitHub username
sed -i 's/ikateclab/YOUR_GITHUB_USERNAME/g' go.mod
sed -i 's/ikateclab/YOUR_GITHUB_USERNAME/g' cmd/main.go
sed -i 's/ikateclab/YOUR_GITHUB_USERNAME/g' cmd/gorm-gen/main.go
# ... update all import paths

# 3. Initialize git and push
git init
git add .
git commit -m "Initial release preparation"
git branch -M main
git remote add origin https://github.com/YOUR_GITHUB_USERNAME/gorm-tracked-updates.git
git push -u origin main
```

### 2. Version Tagging

```bash
# Create and push version tag
git tag v1.0.0
git push origin v1.0.0
```

### 3. Verify Release

```bash
# Test installation
go install github.com/YOUR_GITHUB_USERNAME/gorm-tracked-updates/cmd/gorm-gen@latest

# Test as dependency
go mod init test-project
go get github.com/YOUR_GITHUB_USERNAME/gorm-tracked-updates@latest
```

## Package Structure

The package provides:

### 1. **CLI Tools**
- `cmd/gorm-gen` - go:generate integration tool
- `cmd/main` - standalone CLI tool

### 2. **Libraries**
- `pkg/diffgen` - Diff generation library
- `pkg/clonegen` - Clone generation library

### 3. **Examples**
- `examples/` - Working examples and demos

## Installation for Users

### As CLI Tool
```bash
go install github.com/YOUR_GITHUB_USERNAME/gorm-tracked-updates/cmd/gorm-gen@latest
```

### As Library Dependency
```bash
go get github.com/YOUR_GITHUB_USERNAME/gorm-tracked-updates@latest
```

## Usage Patterns

### 1. go:generate Integration
```go
//go:generate gorm-gen
```

### 2. Library Usage
```go
import (
    "github.com/YOUR_GITHUB_USERNAME/gorm-tracked-updates/pkg/diffgen"
    "github.com/YOUR_GITHUB_USERNAME/gorm-tracked-updates/pkg/clonegen"
)
```

## CI/CD Pipeline

The repository includes:

1. **Continuous Integration** (`.github/workflows/ci.yml`)
   - Tests on multiple Go versions
   - Code generation verification
   - Benchmark testing

2. **Release Automation** (`.github/workflows/release.yml`)
   - Automatic binary builds
   - Multi-platform releases
   - GitHub release creation

## Versioning

Follow semantic versioning:
- `v1.0.0` - Initial stable release
- `v1.0.1` - Bug fixes
- `v1.1.0` - New features
- `v2.0.0` - Breaking changes

## Distribution

The package will be available through:

1. **Go Module Proxy** - Automatic via go.mod
2. **GitHub Releases** - Binary downloads
3. **pkg.go.dev** - Documentation

## Post-Release

1. **Documentation**: Ensure pkg.go.dev updates
2. **Community**: Share on Go forums/communities
3. **Maintenance**: Monitor issues and PRs
4. **Updates**: Regular dependency updates

## Troubleshooting

### Common Issues

1. **Import Path Errors**: Ensure all files use correct module path
2. **Version Conflicts**: Use `@latest` or specific version tags
3. **Build Failures**: Check Go version compatibility (1.21+)

### Support

- GitHub Issues for bug reports
- GitHub Discussions for questions
- Documentation at pkg.go.dev
