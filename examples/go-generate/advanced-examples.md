# Advanced go:generate Examples

## 1. Basic Usage

```go
package models

//go:generate gorm-gen

type User struct {
    ID   uint
    Name string
}
```

## 2. Generate Only Clone Methods

```go
package models

//go:generate gorm-gen -types=clone

type User struct {
    ID   uint
    Name string
}
```

## 3. Generate Only Diff Methods

```go
package models

//go:generate gorm-gen -types=diff

type User struct {
    ID   uint
    Name string
}
```

## 4. Generate for Specific Package

```go
package main

//go:generate gorm-gen -package=./models

func main() {
    // This will generate for the models package
}
```

## 5. Generate to Different Output Directory

```go
package models

//go:generate gorm-gen -output=../generated

type User struct {
    ID   uint
    Name string
}
```

## 6. Multiple Directives in One File

```go
package models

//go:generate gorm-gen -types=clone -output=./clone
//go:generate gorm-gen -types=diff -output=./diff

type User struct {
    ID   uint
    Name string
}
```

## 7. Conditional Generation

```go
package models

//go:generate sh -c "if [ \"$GENERATE_CLONE\" = \"true\" ]; then gorm-gen -types=clone; fi"
//go:generate sh -c "if [ \"$GENERATE_DIFF\" = \"true\" ]; then gorm-gen -types=diff; fi"

type User struct {
    ID   uint
    Name string
}
```

## 8. Integration with Build Tags

```go
//go:build generate
// +build generate

package models

//go:generate gorm-gen

type User struct {
    ID   uint
    Name string
}
```

## 9. Using with Makefile

```makefile
generate:
	GENERATE_CLONE=true GENERATE_DIFF=true go generate ./...

generate-clone:
	GENERATE_CLONE=true go generate ./...

generate-diff:
	GENERATE_DIFF=true go generate ./...
```

## 10. CI/CD Integration

### GitHub Actions

```yaml
name: Generate Code
on: [push, pull_request]

jobs:
  generate:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
    - uses: actions/setup-go@v3
      with:
        go-version: '1.21'
    
    - name: Install generator
      run: go install ./cmd/gorm-gen
    
    - name: Generate code
      run: go generate ./...
    
    - name: Verify no changes
      run: git diff --exit-code
```

### GitLab CI

```yaml
generate:
  stage: build
  script:
    - go install ./cmd/gorm-gen
    - go generate ./...
    - git diff --exit-code
```

## 11. IDE Integration

### VS Code Settings

```json
{
    "go.generateOnSave": "workspace",
    "go.toolsManagement.autoUpdate": true
}
```

### GoLand File Watchers

1. Go to Settings → Tools → File Watchers
2. Add new watcher for `*.go` files
3. Program: `go`
4. Arguments: `generate $FileDir$`

## 12. Best Practices

### File Organization

```
project/
├── models/
│   ├── user.go          // Contains //go:generate directive
│   ├── address.go       // Additional models
│   ├── clone.go         // Generated clone methods
│   └── diff.go          // Generated diff methods
├── cmd/
│   └── gorm-gen/        // Generator tool
└── Makefile             // Build automation
```

### Error Handling

```go
//go:generate sh -c "gorm-gen || (echo 'Generation failed' && exit 1)"
```

### Version Pinning

```go
//go:generate go run github.com/your-org/gorm-gen@v1.2.3
```

## 13. Debugging Generation

### Verbose Output

```go
//go:generate gorm-gen -package=. -types=clone,diff -v
```

### Dry Run

```go
//go:generate gorm-gen -package=. -dry-run
```

## 14. Custom Templates

```go
//go:generate gorm-gen -package=. -template=./custom-template.tmpl
```

## 15. Multiple Packages

```go
//go:generate gorm-gen -package=./models
//go:generate gorm-gen -package=./entities
//go:generate gorm-gen -package=./dto
```
