# go:generate Integration Example

This example demonstrates how to use `//go:generate` to automatically generate clone and diff methods for your GORM models.

## Quick Start

```bash
# 1. Install the generator tool
make install

# 2. Generate clone and diff methods
make generate

# 3. Run the demo
make demo
```

## Manual Usage

```bash
# Install the tool
go install ./cmd/gorm-gen

# Generate for current directory
go generate

# Or run manually
gorm-gen -package=./models
```

## go:generate Directives

Add this directive to any Go file in your package:

```go
//go:generate gorm-gen
```

### Advanced Usage

```go
// Generate only clone methods
//go:generate gorm-gen -types=clone

// Generate only diff methods  
//go:generate gorm-gen -types=diff

// Generate for specific package
//go:generate gorm-gen -package=./models

// Generate to different output directory
//go:generate gorm-gen -package=./models -output=./generated
```

## Generated Files

After running `go generate`, you'll get:

- `clone.go` - Contains `Clone()` methods for all structs
- `diff.go` - Contains `Diff()` methods for all structs

## Usage in Code

```go
// Clone before modifications
backup := user.Clone()

// Make changes
user.Name = "New Name"
user.Email = "new@example.com"

// Generate diff for GORM updates
changes := backup.Diff(user)

// Selective GORM update
result := db.Model(&user).Updates(changes)
```

## Integration with CI/CD

Add to your build pipeline:

```yaml
# GitHub Actions example
- name: Generate code
  run: go generate ./...

- name: Verify no changes
  run: git diff --exit-code
```

## IDE Integration

Most Go IDEs support `go:generate`:

- **VS Code**: Right-click → "Go: Generate"
- **GoLand**: Right-click → "Go Generate"
- **Vim/Neovim**: `:GoGenerate`

## Benefits

✅ **Automatic**: Runs with `go generate`  
✅ **Integrated**: Part of your normal Go workflow  
✅ **Versioned**: Generated code is committed to git  
✅ **Fast**: Only regenerates when needed  
✅ **Reliable**: Same output every time  
