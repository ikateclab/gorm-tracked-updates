package clonegen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"strings"
	"text/template"
)

// simpleCloneTemplate contains the embedded template for simple structs (no complex fields).
// The template file must exist at build time for the embed directive to work.
//go:embed templates/simple_clone.tmpl
var simpleCloneTemplate string

// complexCloneTemplate contains the embedded template for complex structs (has fields needing deep cloning).
// The template file must exist at build time for the embed directive to work.
//go:embed templates/complex_clone.tmpl
var complexCloneTemplate string

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

// String returns the string representation of FieldType for template usage
func (ft FieldType) String() string {
	switch ft {
	case FieldTypeSimple:
		return "Simple"
	case FieldTypeStruct:
		return "Struct"
	case FieldTypeStructPtr:
		return "StructPtr"
	case FieldTypeSlice:
		return "Slice"
	case FieldTypeMap:
		return "Map"
	case FieldTypeInterface:
		return "Interface"
	case FieldTypeComplex:
		return "Complex"
	default:
		return "Unknown"
	}
}

// StructInfo represents information about a struct
type StructInfo struct {
	Name       string
	Fields     []StructField
	ImportPath string
	Package    string
}

// HasComplexFields returns true if the struct has any fields that need deep cloning
func (s StructInfo) HasComplexFields() bool {
	for _, field := range s.Fields {
		if field.FieldType != FieldTypeSimple {
			return true
		}
	}
	return false
}

// GetComplexFields returns only the fields that need deep cloning
func (s StructInfo) GetComplexFields() []StructField {
	var complexFields []StructField
	for _, field := range s.Fields {
		if field.FieldType != FieldTypeSimple {
			complexFields = append(complexFields, field)
		}
	}
	return complexFields
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
	// Parse the AST
	node, packageName, err := g.parseFileAST(filePath)
	if err != nil {
		return err
	}

	// Extract imports
	g.extractImports(node.Imports)

	// Collect struct names for reference
	g.collectStructNames(node)

	// Extract struct details
	return g.extractStructDetails(node, packageName)
}

// parseFileAST parses a Go file and returns the AST node and package name
func (g *CloneGenerator) parseFileAST(filePath string) (*ast.File, string, error) {
	// Set up the file set
	fset := token.NewFileSet()

	// Parse the file
	node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, "", fmt.Errorf("error parsing file: %v", err)
	}

	return node, node.Name.Name, nil
}

// collectStructNames collects struct names for reference during type determination
func (g *CloneGenerator) collectStructNames(node *ast.File) {
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				g.KnownStructs[typeSpec.Name.Name] = true
			}
		}
		return true
	})
}

// extractStructDetails extracts detailed struct information from AST
func (g *CloneGenerator) extractStructDetails(node *ast.File, packageName string) error {
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if structType, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				structInfo := StructInfo{
					Name:    typeSpec.Name.Name,
					Package: packageName,
				}

				// Extract fields
				structInfo.Fields = g.extractFields(structType)
				g.Structs = append(g.Structs, structInfo)
			}
		}
		return true
	})

	return nil
}

// extractFields extracts field information from a struct type
func (g *CloneGenerator) extractFields(structType *ast.StructType) []StructField {
	var fields []StructField

	for _, field := range structType.Fields.List {
		fieldType := g.getTypeString(field.Type)
		fieldTypeCategory := g.categorizeFieldType(fieldType)

		// Handle multiple field names (e.g., a, b int)
		if len(field.Names) > 0 {
			for _, name := range field.Names {
				fields = append(fields, StructField{
					Name:      name.Name,
					Type:      fieldType,
					FieldType: fieldTypeCategory,
				})
			}
		} else {
			// Anonymous field
			fields = append(fields, StructField{
				Name:      fieldType,
				Type:      fieldType,
				FieldType: fieldTypeCategory,
			})
		}
	}

	return fields
}

// extractImports extracts import information from AST imports
func (g *CloneGenerator) extractImports(imports []*ast.ImportSpec) {
	for _, imp := range imports {
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
	case baseType == "json.RawMessage" || baseType == "datatypes.JSON":
		return FieldTypeSlice // Treat as slice since they're []byte
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
		"time.Time":  true, // Add time.Time as simple type since it's copyable by value
	}
	return simpleTypes[typeName]
}

// getRequiredImports determines which imports are needed for the generated code
func (g *CloneGenerator) getRequiredImports() []string {
	var imports []string
	importSet := make(map[string]bool)

	// Check all struct fields for types that need imports
	for _, structInfo := range g.Structs {
		for _, field := range structInfo.Fields {
			// Check if field type contains datatypes.JSON
			if strings.Contains(field.Type, "datatypes.JSON") {
				if !importSet["gorm.io/datatypes"] {
					imports = append(imports, "gorm.io/datatypes")
					importSet["gorm.io/datatypes"] = true
				}
			}
		}
	}

	return imports
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
	requiredImports := g.getRequiredImports()
	if len(requiredImports) > 0 {
		buf.WriteString("import (\n")
		for _, imp := range requiredImports {
			fmt.Fprintf(&buf, "\t\"%s\"\n", imp)
		}
		buf.WriteString(")\n\n")
	}

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

// loadCloneTemplate loads the appropriate clone template based on complexity
func (g *CloneGenerator) loadCloneTemplate(isComplex bool) (*template.Template, error) {
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
		"isMapOfStruct": func(s string) bool {
			if !strings.HasPrefix(s, "map[") {
				return false
			}
			// Extract value type from map[key]value
			parts := strings.Split(s, "]")
			if len(parts) < 2 {
				return false
			}
			valueType := strings.TrimPrefix(parts[1], "*")
			return g.KnownStructs[valueType]
		},
		"isMapOfStructPtr": func(s string) bool {
			if !strings.HasPrefix(s, "map[") {
				return false
			}
			// Extract value type from map[key]value
			parts := strings.Split(s, "]")
			if len(parts) < 2 {
				return false
			}
			valueType := parts[1]
			if !strings.HasPrefix(valueType, "*") {
				return false
			}
			valueType = strings.TrimPrefix(valueType, "*")
			return g.KnownStructs[valueType]
		},
		"getSliceElementType": func(s string) string {
			return strings.TrimPrefix(s, "[]")
		},
		"getSliceElementTypeName": func(s string) string {
			elementType := strings.TrimPrefix(s, "[]")
			return strings.TrimPrefix(elementType, "*")
		},
	}

	// Choose template based on complexity
	var templateContent string
	if isComplex {
		templateContent = complexCloneTemplate
	} else {
		templateContent = simpleCloneTemplate
	}

	// Parse the embedded template
	tmpl, err := template.New("clone").Funcs(funcMap).Parse(templateContent)
	if err != nil {
		return nil, fmt.Errorf("error parsing embedded template: %v", err)
	}

	return tmpl, nil
}

// generateCloneMethod generates a clone method for a struct
func (g *CloneGenerator) generateCloneMethod(structInfo StructInfo) (string, error) {
	// Determine if struct has complex fields
	isComplex := structInfo.HasComplexFields()

	// Load appropriate template
	tmpl, err := g.loadCloneTemplate(isComplex)
	if err != nil {
		return "", err
	}

	// Prepare data for template
	var data interface{}
	if isComplex {
		data = struct {
			StructInfo
			ComplexFields []StructField
		}{
			StructInfo:    structInfo,
			ComplexFields: structInfo.GetComplexFields(),
		}
	} else {
		data = structInfo
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
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
		node, _, err := g.parseFileAST(filePath)
		if err != nil {
			return fmt.Errorf("error parsing file %s: %v", filePath, err)
		}

		// Collect struct names
		g.collectStructNames(node)
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
