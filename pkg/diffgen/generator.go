package diffgen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// StructField represents a field in a struct
type StructField struct {
	Name      string
	Type      string
	FieldType FieldType
	Tag       string // Struct tag for the field
}

// FieldType categorizes the field type for diff generation
type FieldType int

const (
	FieldTypeSimple    FieldType = iota // Primitives, strings, etc.
	FieldTypeStruct                     // Custom struct types
	FieldTypeStructPtr                  // Pointer to custom struct
	FieldTypeSlice                      // Slice of any type
	FieldTypeMap                        // Map of any type
	FieldTypeInterface                  // Interface
	FieldTypeJSON                       // JSON fields with gorm:"serializer:json"
	FieldTypeComplex                    // Any other complex type
)

// StructInfo represents information about a struct
type StructInfo struct {
	Name       string
	Fields     []StructField
	ImportPath string
	Package    string
}

// DiffGenerator handles the code generation for struct diff functions
type DiffGenerator struct {
	Structs      []StructInfo
	KnownStructs map[string]bool
	Imports      map[string]string
}

// New creates a new DiffGenerator
func New() *DiffGenerator {
	return &DiffGenerator{
		KnownStructs: make(map[string]bool),
		Imports:      make(map[string]string),
	}
}

// ParseFile parses a Go file and extracts struct information
func (g *DiffGenerator) ParseFile(filePath string) error {
	// Set up the file set
	fset := token.NewFileSet()

	// Parse the file
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("error parsing file: %v", err)
	}

	// Extract package name
	packageName := node.Name.Name

	// First pass: collect struct names
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				g.KnownStructs[typeSpec.Name.Name] = true
			}
		}
		return true
	})

	// Extract imports
	for _, imp := range node.Imports {
		importPath := strings.Trim(imp.Path.Value, "\"")
		var importName string

		if imp.Name != nil {
			importName = imp.Name.Name
		} else {
			// Extract name from path
			parts := strings.Split(importPath, "/")
			importName = parts[len(parts)-1]
		}

		g.Imports[importPath] = importName
	}

	// Second pass: extract struct details
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if structType, ok := typeSpec.Type.(*ast.StructType); ok {
				// Extract fields from struct
				fields := g.extractFields(structType)

				// Add to structs list
				g.Structs = append(g.Structs, StructInfo{
					Name:       typeSpec.Name.Name,
					Fields:     fields,
					ImportPath: filepath.Dir(filePath),
					Package:    packageName,
				})
				return false
			}
		}
		return true
	})

	return nil
}

// extractFields extracts field information from a struct
func (g *DiffGenerator) extractFields(structType *ast.StructType) []StructField {
	var fields []StructField

	for _, field := range structType.Fields.List {
		// Skip embedded or anonymous fields
		if len(field.Names) == 0 {
			continue
		}

		// Get field type as string
		var buf bytes.Buffer
		format.Node(&buf, token.NewFileSet(), field.Type)
		typeStr := buf.String()

		// Get struct tag if present
		var tagStr string
		if field.Tag != nil {
			tagStr = field.Tag.Value
		}

		for _, name := range field.Names {
			// Determine field type category
			fieldType := g.determineFieldType(field.Type, typeStr, tagStr)

			fields = append(fields, StructField{
				Name:      name.Name,
				Type:      typeStr,
				FieldType: fieldType,
				Tag:       tagStr,
			})
		}
	}

	return fields
}

// determineFieldType analyzes a type to determine its category
func (g *DiffGenerator) determineFieldType(expr ast.Expr, typeStr string, tagStr string) FieldType {
	// Check if it's a JSON field first
	if g.isJSONField(tagStr) {
		return FieldTypeJSON
	}

	switch t := expr.(type) {
	case *ast.Ident:
		// Check if it's a known struct
		if g.KnownStructs[t.Name] {
			return FieldTypeStruct
		}
		// Otherwise it's a simple type
		return FieldTypeSimple

	case *ast.StarExpr:
		// Check if it's a pointer to a known struct
		if ident, ok := t.X.(*ast.Ident); ok && g.KnownStructs[ident.Name] {
			return FieldTypeStructPtr
		}
		// Otherwise it's a complex type
		return FieldTypeComplex

	case *ast.ArrayType:
		return FieldTypeSlice

	case *ast.MapType:
		return FieldTypeMap

	case *ast.InterfaceType:
		return FieldTypeInterface

	case *ast.SelectorExpr:
		// External package type, can't determine if it's a struct
		return FieldTypeComplex

	default:
		return FieldTypeComplex
	}
}

// isJSONField checks if a field has gorm:"serializer:json" tag
func (g *DiffGenerator) isJSONField(tagStr string) bool {
	if tagStr == "" {
		return false
	}
	// Remove the backticks from the tag string
	tagStr = strings.Trim(tagStr, "`")
	// Check if it contains gorm:"serializer:json"
	return strings.Contains(tagStr, `gorm:"serializer:json"`) ||
		   strings.Contains(tagStr, `gorm:"serializer:json`) ||
		   strings.Contains(tagStr, `serializer:json"`)
}

// hasJSONFields checks if any struct has JSON fields
func (g *DiffGenerator) hasJSONFields() bool {
	for _, structInfo := range g.Structs {
		for _, field := range structInfo.Fields {
			if field.FieldType == FieldTypeJSON {
				return true
			}
		}
	}
	return false
}

// GenerateCode generates the code for all struct diff functions
func (g *DiffGenerator) GenerateCode() (string, error) {
	var buf bytes.Buffer

	// Generate package declaration
	if len(g.Structs) > 0 {
		fmt.Fprintf(&buf, "package %s\n\n", g.Structs[0].Package)
	} else {
		return "", fmt.Errorf("no structs found")
	}

	// Check if we need GORM imports
	needsGORM := g.hasJSONFields()

	// Generate imports
	fmt.Fprintln(&buf, "import (")
	if needsGORM {
		fmt.Fprintln(&buf, "\t\"encoding/json\"")
	}
	fmt.Fprintln(&buf, "\t\"reflect\"")
	if needsGORM {
		fmt.Fprintln(&buf, "\t\"gorm.io/gorm\"")
		fmt.Fprintln(&buf, "\t\"gorm.io/gorm/clause\"")
	}
	fmt.Fprintln(&buf, ")")
	fmt.Fprintln(&buf)

	// Generate diff functions for each struct
	for _, structInfo := range g.Structs {
		code, err := g.generateDiffFunction(structInfo)
		if err != nil {
			return "", err
		}
		buf.WriteString(code)
		buf.WriteString("\n\n")
	}

	// Format the code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), fmt.Errorf("error formatting code: %v", err)
	}

	return string(formatted), nil
}

// Template for the diff function
const diffFunctionTmpl = `
// Diff compares this {{.Name}} instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a {{.Name}}) Diff(b {{.Name}}) map[string]interface{} {
	diff := make(map[string]interface{})

	{{range .Fields}}
	// Compare {{.Name}}
	{{if eq .FieldType 0}}
	// Simple type comparison
	if a.{{.Name}} != b.{{.Name}} {
		diff["{{.Name}}"] = b.{{.Name}}
	}
	{{else if eq .FieldType 1}}
	// Struct type comparison
	if !reflect.DeepEqual(a.{{.Name}}, b.{{.Name}}) {
		nestedDiff := a.{{.Name}}.Diff(b.{{.Name}})
		if len(nestedDiff) > 0 {
			diff["{{.Name}}"] = nestedDiff
		}
	}
	{{else if eq .FieldType 2}}
	// Pointer to struct comparison
	if !reflect.DeepEqual(a.{{.Name}}, b.{{.Name}}) {
		if a.{{.Name}} == nil || b.{{.Name}} == nil {
			diff["{{.Name}}"] = b.{{.Name}}
		} else {
			nestedDiff := (*a.{{.Name}}).Diff(*b.{{.Name}})
			if len(nestedDiff) > 0 {
				diff["{{.Name}}"] = nestedDiff
			}
		}
	}
	{{else if eq .FieldType 6}}
	// JSON field comparison - use GORM JSON merge expression
	if !reflect.DeepEqual(a.{{.Name}}, b.{{.Name}}) {
		jsonValue, err := json.Marshal(b.{{.Name}})
		if err == nil {
			diff["{{.Name}}"] = gorm.Expr("? || ?", clause.Column{Name: "{{.Name}}"}, string(jsonValue))
		} else {
			// Fallback to regular assignment if JSON marshaling fails
			diff["{{.Name}}"] = b.{{.Name}}
		}
	}
	{{else}}
	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.{{.Name}}, b.{{.Name}}) {
		diff["{{.Name}}"] = b.{{.Name}}
	}
	{{end}}
	{{end}}

	return diff
}
`

// generateDiffFunction generates a diff function for a struct
func (g *DiffGenerator) generateDiffFunction(structInfo StructInfo) (string, error) {
	// Create template funcs
	funcMap := template.FuncMap{
		"trimStar": func(s string) string {
			return strings.TrimPrefix(s, "*")
		},
	}

	// Parse the template
	tmpl, err := template.New("diff").Funcs(funcMap).Parse(diffFunctionTmpl)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, structInfo); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return buf.String(), nil
}

// WriteToFile writes the generated code to a file
func (g *DiffGenerator) WriteToFile(filePath string) error {
	code, err := g.GenerateCode()
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(code), 0644)
}

// ParseFiles parses multiple Go files and extracts struct information
func (g *DiffGenerator) ParseFiles(filePaths []string) error {
	// First pass: collect all struct names from all files
	for _, filePath := range filePaths {
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("error parsing file %s: %v", filePath, err)
		}

		// Collect struct names
		ast.Inspect(node, func(n ast.Node) bool {
			if typeSpec, ok := n.(*ast.TypeSpec); ok {
				if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
					g.KnownStructs[typeSpec.Name.Name] = true
				}
			}
			return true
		})

		// Extract imports
		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, "\"")
			var importName string

			if imp.Name != nil {
				importName = imp.Name.Name
			} else {
				// Extract name from path
				parts := strings.Split(importPath, "/")
				importName = parts[len(parts)-1]
			}

			g.Imports[importPath] = importName
		}
	}

	// Second pass: extract struct details now that we know all struct names
	for _, filePath := range filePaths {
		err := g.ParseFile(filePath)
		if err != nil {
			return err
		}
	}

	return nil
}

// ParseDirectory parses all .go files in a directory and extracts struct information
func (g *DiffGenerator) ParseDirectory(dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("error reading directory %s: %v", dirPath, err)
	}

	var goFiles []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".go") &&
		   !strings.HasSuffix(file.Name(), "_test.go") &&
		   file.Name() != "clone.go" && file.Name() != "diff.go" {
			goFiles = append(goFiles, dirPath+"/"+file.Name())
		}
	}

	if len(goFiles) == 0 {
		return fmt.Errorf("no Go files found in directory %s", dirPath)
	}

	return g.ParseFiles(goFiles)
}

// WriteToPackageDir writes the generated code to diff.go in the specified directory
func (g *DiffGenerator) WriteToPackageDir(packageDir string) error {
	code, err := g.GenerateCode()
	if err != nil {
		return err
	}

	filePath := packageDir + "/diff.go"
	return os.WriteFile(filePath, []byte(code), 0644)
}
