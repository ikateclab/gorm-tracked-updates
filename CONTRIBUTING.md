# Contributing to GORM Tracked Updates

Thank you for your interest in contributing to GORM Tracked Updates! This document provides guidelines and information for contributors.

## Getting Started

### Prerequisites

- Go 1.21 or later
- Git
- Basic understanding of Go AST and code generation

### Development Setup

```bash
# 1. Fork and clone the repository
git clone https://github.com/ikateclab/gorm-tracked-updates.git
cd gorm-tracked-updates

# 2. Install dependencies
go mod download

# 3. Run tests to ensure everything works
go test ./...

# 4. Install development tools
go install ./cmd/gorm-gen
```

## Development Workflow

### 1. Create a Branch

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-number
```

### 2. Make Changes

- Follow existing code patterns
- Add tests for new functionality
- Update documentation as needed
- Ensure code generation works correctly

### 3. Test Your Changes

```bash
# Run all tests
go test ./...

# Run benchmarks
go test -bench=. ./...

# Test code generation
cd examples/go-generate
make clean && make generate && make test
```

### 4. Submit Pull Request

- Ensure all tests pass
- Update documentation
- Add clear commit messages
- Reference any related issues

## Code Style

### Go Code

- Follow standard Go formatting (`gofmt`)
- Use meaningful variable and function names
- Add comments for exported functions
- Keep functions focused and small

### Generated Code

- Generated code should be clean and readable
- Include build tags and generation comments
- Follow Go naming conventions
- Optimize for performance

## Testing

### Unit Tests

```bash
# Run specific package tests
go test ./pkg/diffgen -v
go test ./pkg/clonegen -v

# Run with coverage
go test -cover ./...

# For Go 1.24.0 compatibility (due to sonic dependency)
go test -ldflags="-checklinkname=0" -cover ./...
```

### Integration Tests

```bash
# Test code generation end-to-end
cd examples/go-generate
make clean && make generate && make demo
```

### Benchmarks

```bash
# Run performance benchmarks
cd examples/performance
go test -bench=. -benchmem
```

## Documentation

### Code Documentation

- Document all exported functions and types
- Include usage examples in godoc
- Explain complex algorithms

### User Documentation

- Update README.md for new features
- Add examples to demonstrate usage
- Update relevant documentation in `docs/`

## Architecture

### Package Structure

```
pkg/
├── diffgen/          # Diff generation logic
│   ├── generator.go  # Main generator
│   └── templates.go  # Code templates
└── clonegen/         # Clone generation logic
    ├── generator.go  # Main generator
    └── templates.go  # Code templates
```

### Key Components

1. **AST Parsing**: Extract struct definitions from Go source
2. **Type Analysis**: Categorize fields for optimal code generation
3. **Template Engine**: Generate type-safe code
4. **File Management**: Handle output and package organization

## Adding New Features

### New Field Types

1. Add type detection in `analyzer.go`
2. Create template in `templates.go`
3. Add tests for the new type
4. Update documentation

### New Generation Options

1. Add command-line flags
2. Update generator configuration
3. Add tests and examples
4. Document the new option

## Performance Considerations

- Generated code should be faster than reflection
- Minimize memory allocations
- Use type-specific optimizations
- Benchmark new features

## Common Issues

### Import Path Problems

- Ensure all imports use the correct module path
- Test with `go mod tidy`
- Verify examples work with the new paths

### Code Generation Failures

- Check AST parsing logic
- Verify template syntax
- Test with various struct types
- Handle edge cases (nil pointers, empty slices)

## Release Process

### Version Bumping

- Follow semantic versioning
- Update CHANGELOG.md
- Tag releases properly

### Testing Before Release

```bash
# Full test suite
go test ./...

# For Go 1.24.0 compatibility
go test -ldflags="-checklinkname=0" ./...

# Integration tests
cd examples/go-generate && make all

# Performance regression tests
cd examples/performance && go test -bench=.
```

## Getting Help

- **Issues**: Report bugs or request features
- **Discussions**: Ask questions or discuss ideas
- **Code Review**: All PRs are reviewed by maintainers

## Code of Conduct

- Be respectful and inclusive
- Focus on constructive feedback
- Help others learn and grow
- Maintain a positive environment

## Recognition

Contributors will be recognized in:
- CHANGELOG.md for significant contributions
- README.md contributors section
- GitHub contributors page

Thank you for contributing to GORM Tracked Updates!
