# GORM Tracked Updates - Code Generators

A collection of high-performance code generators for Go structs, designed to work seamlessly with GORM for efficient database operations.

## Overview

This project provides two main code generators:

1. **DiffGen** - Generates efficient diff functions for detecting changes between struct instances
2. **CloneGen** - Generates performant deep clone methods for creating independent copies

Both generators are designed to work together for optimal GORM workflows: clone → modify → diff → update.

## Recent Updates

**🚀 Performance & JSON Enhancements:**
- **Sonic Integration**: Migrated to `github.com/bytedance/sonic` for JSON operations (3.9x faster unmarshal)
- **Smart Type Detection**: Enhanced detection for custom slice types and JSON fields
- **GORM Expression Optimization**: Automatic JSON field merging with proper GORM expressions
- **Performance Improvements**: Updated benchmarks show 3.5x faster cloning vs reflection

## Project Structure

```
gorm-tracked-updates/
├── cmd/
│   ├── main.go                    # Main CLI tool
│   └── gorm-gen/
│       └── main.go               # go:generate integration tool
├── pkg/
│   ├── diffgen/
│   │   ├── generator.go           # Diff generator implementation
│   │   └── generator_test.go      # Diff generator tests
│   └── clonegen/
│       ├── generator.go           # Clone generator implementation
│       └── generator_test.go      # Clone generator tests
├── examples/
│   ├── structs/                   # Shared struct definitions
│   ├── diff-demo/                 # Diff generator demo
│   ├── clone-demo/                # Clone generator demo
│   ├── multi-file-demo/           # Multi-file generation demo
│   ├── multi-file/                # Multi-file example structs
│   ├── go-generate/               # go:generate integration example
│   └── performance/               # Performance benchmarks
├── testdata/                      # Test generated files
└── docs/                          # Documentation
    ├── DIFFGEN.md               # Diff generator documentation
    └── CLONEGEN.md              # Clone generator documentation
```

## Quick Start

### Option 1: go:generate Integration (Recommended)

```bash
# 1. Install the generator tool
cd examples/go-generate
make install

# 2. Generate clone and diff methods
make generate

# 3. Run the demo
make demo
```

Or manually:
```bash
# Install the tool
go install ./cmd/gorm-gen

# Add to your Go files:
//go:generate gorm-gen

# Generate code
go generate ./...
```

### Option 2: Direct CLI Usage

```bash
# Generate both diff functions and clone methods from a directory
go run cmd/main.go

# Or run individual demos
go run examples/diff-demo/main.go
go run examples/clone-demo/main.go
go run examples/multi-file-demo/main.go
```

### Generated Code Usage

```go
// Clone for backup
backup := user.Clone()

// Modify the user
user.Name = "New Name"
user.Email = "new@example.com"

// Generate diff for GORM update
changes := backup.Diff(user)

// Use diff for selective GORM update
db.Model(&user).Updates(changes)
```

## Features

### DiffGen Features
- **Selective Updates**: Only changed fields in diff map
- **Nested Struct Support**: Recursive diff for complex structures
- **Type Safety**: No reflection overhead in generated code
- **GORM Integration**: Perfect for `Updates()` method
- **High-Performance JSON**: Uses Sonic library for 3.9x faster JSON operations
- **Smart GORM Expressions**: Automatic JSON field merging with proper GORM expressions

### CloneGen Features
- **Deep Cloning**: Complete memory independence
- **Performance**: 3.5x faster than reflection, 22x faster than JSON
- **Memory Safety**: Proper nil handling and reference management
- **Type Optimization**: Specialized handling for each field type

## Performance

Benchmark results show significant performance improvements:

| Method | DiffGen | CloneGen (vs Reflection) | CloneGen (vs JSON) |
|--------|---------|-------------------------|-------------------|
| Performance | Type-safe, no reflection | 3.5x faster | 22x faster |
| Memory | Minimal allocations | Independent copies | No serialization overhead |
| Type Safety | ✅ Compile-time | ✅ Compile-time | ✅ Compile-time |

**Latest Benchmark Results:**
```
BenchmarkCloneGenerated-14     	 3789199	       330.7 ns/op
BenchmarkCloneReflection-14    	  961155	      1165 ns/op
BenchmarkCloneJSON-14          	  166378	      7261 ns/op
```

**JSON Performance (with Sonic integration):**
- **Marshal operations**: Native JSON vs Sonic performance varies by data size
- **Unmarshal operations**: Sonic provides 3.9x performance improvement
- **Overall**: 2.3x performance improvement in realistic usage patterns

## Supported Field Types

Both generators handle all Go field types with optimized performance:

- **Simple Types**: `string`, `int`, `bool`, `float64`, etc. (direct assignment)
- **Struct Types**: Nested structs with recursive processing
- **Pointer Types**: `*Person`, `*Address` with nil safety
- **Slice Types**: `[]Contact`, `[]*Person` with element cloning
- **Map Types**: `map[string]interface{}` with key-value copying
- **Interface Types**: `interface{}` with reflection fallback
- **JSON Types**: `datatypes.JSON`, custom JSON slices with Sonic performance
- **Time Types**: `time.Time`, `*time.Time` with proper equality checking

## GORM Integration

Perfect workflow for tracked updates:

```go
// 1. Clone before modifications
backup := user.Clone()

// 2. Make changes
user.UpdateFromRequest(request)

// 3. Generate diff
changes := backup.Diff(user)

// 4. Selective GORM update
result := db.Model(&user).Updates(changes)

// 5. Only changed fields are updated in database
```

### Advanced GORM Features

The generated diff methods support advanced GORM features with high-performance JSON handling:

```go
// JSON field merging with GORM expressions
// For fields with `gorm:"serializer:json"` tags
updateMap := backup.Diff(user)
// JSON fields are automatically handled with proper GORM expressions using Sonic for performance

// Example of generated GORM expression for JSON fields:
// updateMap["settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, jsonValue)
```

**Key JSON Features:**
- **Sonic Integration**: Uses `github.com/bytedance/sonic` for 3.9x faster JSON operations
- **Smart Type Detection**: Automatically detects JSON fields, custom slices, and complex types
- **GORM Expression Generation**: Proper JSON merging with `gorm.Expr` for database efficiency
- **Fallback Safety**: Graceful fallback to standard assignment if JSON operations fail

## Testing

Run comprehensive tests:

```bash
# Run all tests
go test ./...

# Run specific package tests
go test ./pkg/diffgen -v
go test ./pkg/clonegen -v

# Run performance benchmarks
cd examples/performance && go test -bench=. -v
```

## Examples

See the `examples/` directory for:
- **structs/**: Shared struct definitions with generated code
- **diff-demo/**: DiffGen demonstration
- **clone-demo/**: CloneGen demonstration
- **multi-file-demo/**: Multi-file generation demonstration
- **multi-file/**: Multi-file example structs
- **go-generate/**: go:generate integration example
- **performance/**: Performance benchmarks

## go:generate Integration

The project includes a dedicated `gorm-gen` tool for seamless go:generate integration:

### Features
- **Automatic Generation**: Integrates with `go generate` workflow
- **Flexible Options**: Generate clone only, diff only, or both
- **Package Support**: Works with any Go package structure
- **CI/CD Ready**: Perfect for automated build pipelines

### Usage
```go
//go:generate gorm-gen
//go:generate gorm-gen -types=clone
//go:generate gorm-gen -types=diff
//go:generate gorm-gen -package=./models -output=./generated
```

### Generated Files
- `clone.go` - Contains `Clone()` methods for all structs
- `diff.go` - Contains `Diff()` methods for all structs

See `examples/go-generate/` for a complete working example.

## Documentation

Detailed documentation available:
- [DiffGen Documentation](docs/DIFFGEN.md)
- [CloneGen Documentation](docs/CLONEGEN.md)
- [go:generate Integration](examples/go-generate/README.md)

## Contributing

1. Follow the existing code patterns
2. Add comprehensive tests for new features
3. Update documentation
4. Ensure performance benchmarks pass

## License

[Add your license here]

## Use Cases

### Database Operations
- Selective GORM updates
- Change tracking
- Audit logging
- Optimistic locking

### API Development
- Request/response diffing
- State management
- Caching strategies
- Data synchronization

### Testing
- Test data setup
- State isolation
- Snapshot testing
- Mock data generation

## Architecture

Both generators follow the same architectural pattern:

1. **AST Parsing**: Parse Go source files to extract struct definitions
2. **Type Analysis**: Categorize field types for optimal handling
3. **Code Generation**: Template-based code generation
4. **Optimization**: Type-specific optimizations for performance

This ensures consistency, maintainability, and extensibility across both generators.

## Dependencies

The project uses carefully selected high-performance dependencies:

- **Core GORM**: `gorm.io/gorm` for database operations and expression generation
- **GORM Datatypes**: `gorm.io/datatypes` for JSON field support
- **Sonic JSON**: `github.com/bytedance/sonic` for high-performance JSON operations
- **UUID Support**: `github.com/google/uuid` for unique identifier generation

All dependencies are focused on performance and production readiness.
