package clonegen

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

// StructField represents a field in a struct
type StructField struct {
	Name      string
	Type      string
	FieldType FieldType
}

// FieldType categorizes the field type for clone generation
type FieldType int

const (
	FieldTypeSimple    FieldType = iota // Primitives, strings, etc.
	FieldTypeStruct                     // Custom struct types
	FieldTypeStructPtr                  // Pointer to custom struct
	FieldTypeSlice                      // Slice of any type
	FieldTypeMap                        // Map of any type
	FieldTypeInterface                  // Interface
	FieldTypeComplex                    // Any other complex type
)

// StructInfo represents information about a struct
type StructInfo struct {
	Name       string
	Fields     []StructField
	ImportPath string
	Package    string
}

// CloneGenerator handles the code generation for struct clone methods
type CloneGenerator struct {
	Structs      []StructInfo
	KnownStructs map[string]bool
	Imports      map[string]string
}

// New creates a new CloneGenerator
func New() *CloneGenerator {
	return &CloneGenerator{
		KnownStructs: make(map[string]bool),
		Imports:      make(map[string]string),
	}
}

// ParseFile parses a Go file and extracts struct information
func (g *CloneGenerator) ParseFile(filePath string) error {
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

	// Second pass: extract struct information
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if structType, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				structInfo := StructInfo{
					Name:    typeSpec.Name.Name,
					Package: packageName,
				}

				// Extract fields
				for _, field := range structType.Fields.List {
					fieldType := g.getTypeString(field.Type)
					fieldTypeCategory := g.categorizeFieldType(fieldType)

					// Handle multiple field names (e.g., a, b int)
					if len(field.Names) > 0 {
						for _, name := range field.Names {
							structInfo.Fields = append(structInfo.Fields, StructField{
								Name:      name.Name,
								Type:      fieldType,
								FieldType: fieldTypeCategory,
							})
						}
					} else {
						// Anonymous field
						structInfo.Fields = append(structInfo.Fields, StructField{
							Name:      fieldType,
							Type:      fieldType,
							FieldType: fieldTypeCategory,
						})
					}
				}

				g.Structs = append(g.Structs, structInfo)
			}
		}
		return true
	})

	return nil
}

// getTypeString converts an ast.Expr to a string representation
func (g *CloneGenerator) getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + g.getTypeString(t.X)
	case *ast.ArrayType:
		if t.Len == nil {
			// Slice
			return "[]" + g.getTypeString(t.Elt)
		}
		// Array (not commonly used, treat as slice for simplicity)
		return "[]" + g.getTypeString(t.Elt)
	case *ast.MapType:
		return "map[" + g.getTypeString(t.Key) + "]" + g.getTypeString(t.Value)
	case *ast.InterfaceType:
		return "interface{}"
	case *ast.SelectorExpr:
		return g.getTypeString(t.X) + "." + t.Sel.Name
	default:
		return "interface{}"
	}
}

// categorizeFieldType determines the category of a field type
func (g *CloneGenerator) categorizeFieldType(fieldType string) FieldType {
	// Remove pointer prefix for analysis
	baseType := strings.TrimPrefix(fieldType, "*")

	// Check if it's a known struct
	if g.KnownStructs[baseType] {
		if strings.HasPrefix(fieldType, "*") {
			return FieldTypeStructPtr
		}
		return FieldTypeStruct
	}

	// Check for built-in types
	switch {
	case strings.HasPrefix(fieldType, "[]"):
		return FieldTypeSlice
	case strings.HasPrefix(fieldType, "map["):
		return FieldTypeMap
	case fieldType == "interface{}" || strings.Contains(fieldType, "interface"):
		return FieldTypeInterface
	case isSimpleType(baseType):
		return FieldTypeSimple
	default:
		return FieldTypeComplex
	}
}

// isSimpleType checks if a type is a simple built-in type
func isSimpleType(typeName string) bool {
	simpleTypes := map[string]bool{
		"bool":       true,
		"string":     true,
		"int":        true,
		"int8":       true,
		"int16":      true,
		"int32":      true,
		"int64":      true,
		"uint":       true,
		"uint8":      true,
		"uint16":     true,
		"uint32":     true,
		"uint64":     true,
		"uintptr":    true,
		"byte":       true,
		"rune":       true,
		"float32":    true,
		"float64":    true,
		"complex64":  true,
		"complex128": true,
	}
	return simpleTypes[typeName]
}

// GenerateCode generates the code for all struct clone methods
func (g *CloneGenerator) GenerateCode() (string, error) {
	var buf bytes.Buffer

	// Generate package declaration
	if len(g.Structs) > 0 {
		fmt.Fprintf(&buf, "package %s\n\n", g.Structs[0].Package)
	} else {
		return "", fmt.Errorf("no structs found")
	}

	// Generate imports if needed
	// Note: reflect import removed as it's not used in generated clone methods
	fmt.Fprintln(&buf)

	// Generate clone methods for each struct
	for _, structInfo := range g.Structs {
		code, err := g.generateCloneMethod(structInfo)
		if err != nil {
			return "", err
		}
		buf.WriteString(code)
		buf.WriteString("\n\n")
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), nil // Return unformatted if formatting fails
	}

	return string(formatted), nil
}

// Template for the clone method
const cloneMethodTmpl = `
// Clone creates a deep copy of the {{.Name}} struct
func (original {{.Name}}) Clone() {{.Name}} {
	clone := {{.Name}}{}

	{{range .Fields}}
	// Clone {{.Name}}
	{{if eq .FieldType 0}}
	// Simple type - direct assignment
	clone.{{.Name}} = original.{{.Name}}
	{{else if eq .FieldType 1}}
	// Struct type - recursive clone
	clone.{{.Name}} = original.{{.Name}}.Clone()
	{{else if eq .FieldType 2}}
	// Pointer to struct - create new instance and clone
	if original.{{.Name}} != nil {
		cloned{{.Name}} := original.{{.Name}}.Clone()
		clone.{{.Name}} = &cloned{{.Name}}
	}
	{{else if eq .FieldType 3}}
	// Slice - create new slice and clone elements
	if original.{{.Name}} != nil {
		clone.{{.Name}} = make({{.Type}}, len(original.{{.Name}}))
		{{if .Type | isSliceOfStructPtr}}
		for i, item := range original.{{.Name}} {
			if item != nil {
				clonedItem := item.Clone()
				clone.{{.Name}}[i] = &clonedItem
			}
		}
		{{else if .Type | isSliceOfStruct}}
		for i, item := range original.{{.Name}} {
			clone.{{.Name}}[i] = item.Clone()
		}
		{{else}}
		copy(clone.{{.Name}}, original.{{.Name}})
		{{end}}
	}
	{{else if eq .FieldType 4}}
	// Map - create new map and copy key-value pairs
	if original.{{.Name}} != nil {
		clone.{{.Name}} = make({{.Type}})
		for k, v := range original.{{.Name}} {
			clone.{{.Name}}[k] = v
		}
	}
	{{else}}
	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.{{.Name}} = original.{{.Name}}
	{{end}}
	{{end}}

	return clone
}
`

// generateCloneMethod generates a clone method for a struct
func (g *CloneGenerator) generateCloneMethod(structInfo StructInfo) (string, error) {
	// Create template funcs
	funcMap := template.FuncMap{
		"trimStar": func(s string) string {
			return strings.TrimPrefix(s, "*")
		},
		"isSliceOfStruct": func(s string) bool {
			if !strings.HasPrefix(s, "[]") {
				return false
			}
			elementType := strings.TrimPrefix(s, "[]")
			elementType = strings.TrimPrefix(elementType, "*")
			return g.KnownStructs[elementType]
		},
		"isSliceOfStructPtr": func(s string) bool {
			if !strings.HasPrefix(s, "[]") {
				return false
			}
			elementType := strings.TrimPrefix(s, "[]")
			if !strings.HasPrefix(elementType, "*") {
				return false
			}
			elementType = strings.TrimPrefix(elementType, "*")
			return g.KnownStructs[elementType]
		},
		"getSliceElementType": func(s string) string {
			return strings.TrimPrefix(s, "[]")
		},
		"getSliceElementTypeName": func(s string) string {
			elementType := strings.TrimPrefix(s, "[]")
			return strings.TrimPrefix(elementType, "*")
		},
	}

	// Parse the template
	tmpl, err := template.New("clone").Funcs(funcMap).Parse(cloneMethodTmpl)
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
func (g *CloneGenerator) WriteToFile(filePath string) error {
	code, err := g.GenerateCode()
	if err != nil {
		return err
	}

	return os.WriteFile(filePath, []byte(code), 0644)
}

// ParseFiles parses multiple Go files and extracts struct information
func (g *CloneGenerator) ParseFiles(filePaths []string) error {
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
func (g *CloneGenerator) ParseDirectory(dirPath string) error {
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

// WriteToPackageDir writes the generated code to clone.go in the specified directory
func (g *CloneGenerator) WriteToPackageDir(packageDir string) error {
	code, err := g.GenerateCode()
	if err != nil {
		return err
	}

	filePath := packageDir + "/clone.go"
	return os.WriteFile(filePath, []byte(code), 0644)
}
