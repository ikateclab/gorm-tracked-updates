package diffgen

import (
	"bytes"
	_ "embed"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

// diffFunctionTemplate contains the embedded template for generating diff functions.
// The template file must exist at build time for the embed directive to work.
//
//go:embed templates/diff_function.tmpl
var diffFunctionTemplate string

// StructField represents a field in a struct
type StructField struct {
	Name            string
	Type            string
	FieldType       FieldType
	Tag             string // Struct tag for the field
	DiffKey         string // Pre-computed key for diff operations (JSON tag name or field name)
	IsNestedInJSONB bool   // Whether this field is a nested struct within a JSONB-annotated parent
}

// FieldType categorizes the field type for diff generation
type FieldType int

const (
	FieldTypeSimple        FieldType = iota // Primitives, strings, etc.
	FieldTypeStruct                         // Custom struct types
	FieldTypeStructPtr                      // Pointer to custom struct
	FieldTypeSlice                          // Slice of any type
	FieldTypeMap                            // Map of any type
	FieldTypeInterface                      // Interface
	FieldTypeJSON                           // JSON fields with gorm:"serializer:json"
	FieldTypeTime                           // time.Time and *time.Time
	FieldTypeUUID                           // uuid.UUID and *uuid.UUID
	FieldTypeGormDeletedAt                  // gorm.DeletedAt
	FieldTypeComparable                     // Other types that support == comparison
	FieldTypeComplex                        // Any other complex type requiring reflection
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
	case FieldTypeJSON:
		return "JSON"
	case FieldTypeTime:
		return "Time"
	case FieldTypeUUID:
		return "UUID"
	case FieldTypeGormDeletedAt:
		return "GormDeletedAt"
	case FieldTypeComparable:
		return "Comparable"
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
	IsJSONB    bool // Whether this struct is annotated with @jsonb
}

// DiffGenerator handles the code generation for struct diff functions
type DiffGenerator struct {
	Structs        []StructInfo
	KnownStructs   map[string]bool
	Imports        map[string]string
	JSONBStructs   map[string]bool // Tracks which structs are used as JSONB columns
	NamingStrategy string          // GORM naming strategy: "snake_case" (default) or "camel_case"
}

// New creates a new DiffGenerator
func New() *DiffGenerator {
	return &DiffGenerator{
		KnownStructs:   make(map[string]bool),
		Imports:        make(map[string]string),
		JSONBStructs:   make(map[string]bool),
		NamingStrategy: "snake_case", // Default to snake_case for backward compatibility
	}
}

// SetNamingStrategy sets the GORM naming strategy for column name generation
func (g *DiffGenerator) SetNamingStrategy(strategy string) {
	g.NamingStrategy = strategy
}

// ParseFile parses a Go file and extracts struct information
func (g *DiffGenerator) ParseFile(filePath string) error {
	// Parse the AST
	node, packageName, err := g.parseFileAST(filePath)
	if err != nil {
		return err
	}

	// Collect struct names for reference
	g.collectStructNames(node)

	// Extract imports
	g.extractImports(node.Imports)

	// Extract struct details
	return g.extractStructDetails(node, filePath, packageName)
}

// parseFileAST parses a Go file and returns the AST node and package name
func (g *DiffGenerator) parseFileAST(filePath string) (*ast.File, string, error) {
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
func (g *DiffGenerator) collectStructNames(node *ast.File) {
	ast.Inspect(node, func(n ast.Node) bool {
		if typeSpec, ok := n.(*ast.TypeSpec); ok {
			if _, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
				g.KnownStructs[typeSpec.Name.Name] = true
			}
		}
		return true
	})
}

// extractStructDetails extracts detailed struct information from AST declarations
func (g *DiffGenerator) extractStructDetails(node *ast.File, filePath, packageName string) error {
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok && genDecl.Tok == token.TYPE {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if structType, ok := typeSpec.Type.(*ast.StructType); ok {
						// Extract fields from struct
						fields := g.extractFields(structType)

						// Check for @jsonb annotation in comments
						// Use genDecl.Doc (declaration comments) instead of typeSpec.Doc
						isJSONB := g.hasJSONBAnnotation(genDecl.Doc)

						// Add to structs list
						g.Structs = append(g.Structs, StructInfo{
							Name:       typeSpec.Name.Name,
							Fields:     fields,
							ImportPath: filepath.Dir(filePath),
							Package:    packageName,
							IsJSONB:    isJSONB,
						})
					}
				}
			}
		}
	}

	return nil
}

// extractImports extracts import information from AST imports
func (g *DiffGenerator) extractImports(imports []*ast.ImportSpec) {
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
		if err := format.Node(&buf, token.NewFileSet(), field.Type); err != nil {
			// Log error and skip this field - we can't process fields we can't format
			fmt.Printf("Warning: Could not format field type for field %v: %v\n", field.Names, err)
			continue
		}
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

// determineKnownTypeByString checks for known types by their string representation
func (g *DiffGenerator) determineKnownTypeByString(typeStr string) FieldType {
	// Check for specific known types by string representation
	// Handle both qualified and unqualified names
	switch typeStr {
	case "time.Time", "*time.Time":
		return FieldTypeTime
	case "uuid.UUID", "*uuid.UUID":
		return FieldTypeUUID
	case "gorm.DeletedAt":
		return FieldTypeGormDeletedAt
	}

	// Also check for common patterns in type strings
	if strings.Contains(typeStr, "uuid.UUID") {
		return FieldTypeUUID
	}
	if strings.Contains(typeStr, "time.Time") {
		return FieldTypeTime
	}
	if strings.Contains(typeStr, "gorm.DeletedAt") {
		return FieldTypeGormDeletedAt
	}

	// Return FieldTypeComplex to indicate no known type found
	return FieldTypeComplex
}

// determineFieldType analyzes a type to determine its category
func (g *DiffGenerator) determineFieldType(expr ast.Expr, typeStr string, tagStr string) FieldType {
	// Check if it's a JSON field - only treat fields with actual database JSON tags as JSON
	if g.isJSONField(tagStr) {
		// Always treat fields with JSON tags as JSON fields for proper gorm.Expr handling
		// This ensures JSONB fields use JSON merging instead of struct comparison
		return FieldTypeJSON
	}

	// For nested JSONB structs without database JSON tags, treat as regular struct fields
	// This prevents nested gorm.Expr calls
	baseType := strings.TrimPrefix(typeStr, "*")
	if g.JSONBStructs[baseType] {
		// This is a nested JSONB struct, but treat as regular struct to avoid nested gorm.Expr
		if strings.HasPrefix(typeStr, "*") {
			return FieldTypeStructPtr
		} else {
			return FieldTypeStruct
		}
	}

	// Check for specific known types by string representation
	if fieldType := g.determineKnownTypeByString(typeStr); fieldType != FieldTypeComplex {
		return fieldType
	}

	return g.determineFieldTypeByAST(expr)
}

// determineFieldTypeByAST analyzes AST expressions to determine field type
func (g *DiffGenerator) determineFieldTypeByAST(expr ast.Expr) FieldType {
	switch t := expr.(type) {
	case *ast.Ident:
		return g.handleIdentType(t)
	case *ast.StarExpr:
		return g.handlePointerType(t)
	case *ast.ArrayType:
		return FieldTypeSlice
	case *ast.MapType:
		return FieldTypeMap
	case *ast.InterfaceType:
		return FieldTypeInterface
	case *ast.SelectorExpr:
		return g.handleSelectorType(t)
	default:
		return FieldTypeComplex
	}
}

// handleIdentType handles ast.Ident expressions
func (g *DiffGenerator) handleIdentType(t *ast.Ident) FieldType {
	if t.Name == "any" {
		return FieldTypeInterface
	}
	// Check for common patterns that indicate slice types (but not JsonbStringSlice with JSON tags)
	if strings.Contains(strings.ToLower(t.Name), "slice") {
		return FieldTypeComplex
	}
	// For known structs, treat as complex to use reflect.DeepEqual (safer for non-comparable structs)
	// Only JSONB fields should be handled specially through field tags
	if g.KnownStructs[t.Name] {
		return FieldTypeComplex
	}
	return FieldTypeSimple
}

// handlePointerType handles ast.StarExpr expressions (pointer types)
func (g *DiffGenerator) handlePointerType(t *ast.StarExpr) FieldType {
	// Check for pointer to time.Time
	if ident, ok := t.X.(*ast.SelectorExpr); ok {
		if x, ok := ident.X.(*ast.Ident); ok && x.Name == "time" && ident.Sel.Name == "Time" {
			return FieldTypeTime
		}
	}
	// Check for pointer to uuid.UUID
	if ident, ok := t.X.(*ast.SelectorExpr); ok {
		if x, ok := ident.X.(*ast.Ident); ok && x.Name == "uuid" && ident.Sel.Name == "UUID" {
			return FieldTypeUUID
		}
	}
	// For pointer to known structs, treat as comparable to avoid relationship handling
	// Only JSONB fields should be handled specially through field tags
	return FieldTypeComparable
}

// handleSelectorType handles ast.SelectorExpr expressions (package.Type)
func (g *DiffGenerator) handleSelectorType(t *ast.SelectorExpr) FieldType {
	// Check for specific external package types
	if x, ok := t.X.(*ast.Ident); ok {
		switch x.Name {
		case "time":
			if t.Sel.Name == "Time" {
				return FieldTypeTime
			}
		case "uuid":
			if t.Sel.Name == "UUID" {
				return FieldTypeUUID
			}
		case "gorm":
			if t.Sel.Name == "DeletedAt" {
				return FieldTypeGormDeletedAt
			}
		case "datatypes":
			// GORM datatypes.JSON has []byte underlying type and supports direct comparison
			if t.Sel.Name == "JSON" {
				return FieldTypeJSON
			}
			// Other GORM datatypes might be comparable
			return FieldTypeComparable
		}
	}
	// For other external types, assume they're comparable unless proven otherwise
	return FieldTypeComparable
}

// isJSONField checks if a field has gorm:"serializer:json" or gorm:"type:jsonb" tag
func (g *DiffGenerator) isJSONField(tagStr string) bool {
	if tagStr == "" {
		return false
	}
	// Remove the backticks from the tag string
	tagStr = strings.Trim(tagStr, "`")
	// Check if it contains gorm:"serializer:json" or gorm:"type:jsonb"
	return strings.Contains(tagStr, `serializer:json`) ||
		strings.Contains(tagStr, `type:jsonb`)
}

// hasJSONBAnnotation checks if a struct has @jsonb annotation in its comments
func (g *DiffGenerator) hasJSONBAnnotation(commentGroup *ast.CommentGroup) bool {
	if commentGroup == nil {
		return false
	}

	for _, comment := range commentGroup.List {
		if strings.Contains(comment.Text, "@jsonb") {
			return true
		}
	}
	return false
}

// extractColumnName extracts the column name from GORM tag or converts field name based on naming strategy
func (g *DiffGenerator) extractColumnName(fieldName, tagStr string) string {
	if tagStr == "" {
		return g.applyNamingStrategy(fieldName)
	}

	// Remove the backticks from the tag string
	tagStr = strings.Trim(tagStr, "`")

	// Look for gorm:"column:columnname" pattern
	re := regexp.MustCompile(`gorm:"[^"]*column:([^;"]*)`)
	matches := re.FindStringSubmatch(tagStr)
	if len(matches) > 1 && matches[1] != "" {
		return matches[1]
	}

	// If no column tag found, apply naming strategy
	return g.applyNamingStrategy(fieldName)
}

// applyNamingStrategy applies the configured GORM naming strategy to a field name
func (g *DiffGenerator) applyNamingStrategy(fieldName string) string {
	switch g.NamingStrategy {
	case "camel_case":
		// Convert to lowerCamelCase using the same logic as CamelCaseNamingStrategy
		return g.toLowerCamelCase(fieldName)
	case "snake_case":
		fallthrough
	default:
		// Default to snake_case
		return g.toSnakeCase(fieldName)
	}
}

// toLowerCamelCase converts a string to camelCase with the first character in lower case
// This matches the behavior of CamelCaseNamingStrategy from gorm-repository
func (g *DiffGenerator) toLowerCamelCase(s string) string {
	if s == "" {
		return s
	}
	runes := []rune(s)
	for i, r := range runes {
		if i == 0 {
			runes[i] = unicode.ToLower(r)
		} else if unicode.IsUpper(r) {
			runes[i] = unicode.ToLower(r)
		} else {
			break
		}
	}
	return string(runes)
}

// toSnakeCase converts CamelCase to snake_case
func (g *DiffGenerator) toSnakeCase(str string) string {
	// Insert underscore before uppercase letters (except the first one)
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	snake := re.ReplaceAllString(str, `${1}_${2}`)
	return strings.ToLower(snake)
}

// extractJSONTagName extracts the JSON tag name from a struct field tag
func (g *DiffGenerator) extractJSONTagName(fieldName, tagStr string) string {
	if tagStr == "" {
		return fieldName
	}

	// Remove the backticks from the tag string
	tagStr = strings.Trim(tagStr, "`")

	// Look for json:"tagname" pattern
	re := regexp.MustCompile(`json:"([^,"]*)`)
	matches := re.FindStringSubmatch(tagStr)
	if len(matches) > 1 && matches[1] != "" && matches[1] != "-" {
		return matches[1]
	}

	// If no JSON tag found or tag is "-", use field name
	return fieldName
}

// computeFieldKeysAndIdentifyJSONB identifies which structs are annotated with @jsonb
// and computes diff keys for all fields
func (g *DiffGenerator) computeFieldKeysAndIdentifyJSONB() {
	// First, mark structs that have @jsonb annotation
	for _, structInfo := range g.Structs {
		if structInfo.IsJSONB {
			g.JSONBStructs[structInfo.Name] = true
		}
	}

	// Second, re-process field types now that we know which structs are JSONB
	for i := range g.Structs {
		for j := range g.Structs[i].Fields {
			field := &g.Structs[i].Fields[j]

			// Only treat fields as JSON if they have actual database JSON tags
			// Nested structs within JSONB should be treated as regular struct fields
			// This prevents nested gorm.Expr calls
			if g.isJSONField(field.Tag) {
				field.FieldType = FieldTypeJSON
			} else {
				// For nested JSONB structs without database JSON tags, treat as regular struct
				baseType := strings.TrimPrefix(field.Type, "*")
				if g.JSONBStructs[baseType] {
					// This is a nested JSONB struct, but treat as regular struct to avoid nested gorm.Expr
					if strings.HasPrefix(field.Type, "*") {
						field.FieldType = FieldTypeStructPtr
					} else {
						field.FieldType = FieldTypeStruct
					}

					// Mark this field as nested within a JSONB parent if the parent struct is JSONB
					if g.JSONBStructs[g.Structs[i].Name] {
						field.IsNestedInJSONB = true
					}
				}
			}

			// Compute diff keys
			// For fields with database JSONB tags, always use field names so GORM can map to column names
			// For nested JSONB structs (without database tags), use JSON tag names for proper JSON path handling
			if g.isJSONField(field.Tag) {
				// This field is stored as JSONB in the database - use field name for GORM column mapping
				field.DiffKey = field.Name
			} else if g.JSONBStructs[g.Structs[i].Name] {
				// Parent struct is a JSONB struct (nested) - use JSON tag names for JSON path handling
				field.DiffKey = g.extractJSONTagName(field.Name, field.Tag)
			} else {
				// Regular struct field - use field names
				field.DiffKey = field.Name
			}
		}
	}
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

// hasDatatypesJSON checks if any struct has datatypes.JSON fields (which need bytes.Equal)
func (g *DiffGenerator) hasDatatypesJSON() bool {
	for _, structInfo := range g.Structs {
		for _, field := range structInfo.Fields {
			if field.FieldType == FieldTypeJSON && field.Type == "datatypes.JSON" {
				return true
			}
		}
	}
	return false
}

// GenerateCode generates the code for all struct diff functions
func (g *DiffGenerator) GenerateCode() (string, error) {
	var buf bytes.Buffer

	// Identify which structs are used as JSONB columns and compute field keys
	g.computeFieldKeysAndIdentifyJSONB()

	// Generate package declaration
	if len(g.Structs) > 0 {
		fmt.Fprintf(&buf, "package %s\n\n", g.Structs[0].Package)
	} else {
		return "", fmt.Errorf("no structs found")
	}

	// Check if we need GORM imports
	needsGORM := g.hasJSONFields()
	needsBytes := g.hasDatatypesJSON()

	// Generate imports
	fmt.Fprintln(&buf, "import (")
	if needsBytes {
		fmt.Fprintln(&buf, "\t\"bytes\"")
	}
	if needsGORM {
		fmt.Fprintln(&buf, "\t\"sort\"")
		fmt.Fprintln(&buf, "\t\"github.com/bytedance/sonic\"")
	}
	fmt.Fprintln(&buf, "\t\"reflect\"")
	if needsGORM {
		fmt.Fprintln(&buf, "\t\"strings\"")
		fmt.Fprintln(&buf, "\t\"gorm.io/gorm\"")
		fmt.Fprintln(&buf, "\t\"gorm.io/gorm/clause\"")
	}
	fmt.Fprintln(&buf, ")")
	fmt.Fprintln(&buf)

	// Generate helper functions if JSON fields are present
	if needsGORM {
		fmt.Fprintln(&buf, "// isEmptyJSON checks if a JSON string represents an empty object or array")
		fmt.Fprintln(&buf, "func isEmptyJSON(jsonStr string) bool {")
		fmt.Fprintln(&buf, "\ttrimmed := strings.TrimSpace(jsonStr)")
		fmt.Fprintln(&buf, "\treturn trimmed == \"{}\" || trimmed == \"[]\" || trimmed == \"null\"")
		fmt.Fprintln(&buf, "}")
		fmt.Fprintln(&buf)

		// Generate buildJSONBSetExpr helper function
		fmt.Fprintln(&buf, "// buildJSONBSetExpr constructs a nested jsonb_set expression for PostgreSQL")
		fmt.Fprintln(&buf, "// to update multiple paths within a JSONB column")
		fmt.Fprintln(&buf, "func buildJSONBSetExpr(columnName string, paths map[string]interface{}) clause.Expr {")
		fmt.Fprintln(&buf, "\t// Start with the original column value (or empty object if NULL)")
		fmt.Fprintln(&buf, "\texpr := \"COALESCE(?::jsonb, '{}'::jsonb)\"")
		fmt.Fprintln(&buf, "\targs := []interface{}{clause.Column{Name: columnName}}")
		fmt.Fprintln(&buf, "")
		fmt.Fprintln(&buf, "\t// Sort paths for consistent ordering")
		fmt.Fprintln(&buf, "\tsortedPaths := make([]string, 0, len(paths))")
		fmt.Fprintln(&buf, "\tfor path := range paths {")
		fmt.Fprintln(&buf, "\t\tsortedPaths = append(sortedPaths, path)")
		fmt.Fprintln(&buf, "\t}")
		fmt.Fprintln(&buf, "\tsort.Strings(sortedPaths)")
		fmt.Fprintln(&buf, "")
		fmt.Fprintln(&buf, "\t// Build nested jsonb_set calls for each path")
		fmt.Fprintln(&buf, "\tfor _, path := range sortedPaths {")
		fmt.Fprintln(&buf, "\t\tvalue := paths[path]")
		fmt.Fprintln(&buf, "")
		fmt.Fprintln(&buf, "\t\t// Convert \"mode\" or \"state.code\" to PostgreSQL array format")
		fmt.Fprintln(&buf, "\t\t// \"mode\" -> {mode}")
		fmt.Fprintln(&buf, "\t\t// \"state.code\" -> {state,code}")
		fmt.Fprintln(&buf, "\t\tpathParts := strings.Split(path, \".\")")
		fmt.Fprintln(&buf, "\t\tpathArray := \"{\" + strings.Join(pathParts, \",\") + \"}\"")
		fmt.Fprintln(&buf, "")
		fmt.Fprintln(&buf, "\t\t// Serialize value to JSON")
		fmt.Fprintln(&buf, "\t\tvalueJSON, err := sonic.Marshal(value)")
		fmt.Fprintln(&buf, "\t\tif err != nil {")
		fmt.Fprintln(&buf, "\t\t\t// Skip this path if we can't marshal the value")
		fmt.Fprintln(&buf, "\t\t\tcontinue")
		fmt.Fprintln(&buf, "\t\t}")
		fmt.Fprintln(&buf, "")
		fmt.Fprintln(&buf, "\t\t// Nest another jsonb_set call")
		fmt.Fprintln(&buf, "\t\texpr = \"jsonb_set(\" + expr + \", '\" + pathArray + \"', ?::jsonb)\"")
		fmt.Fprintln(&buf, "\t\targs = append(args, string(valueJSON))")
		fmt.Fprintln(&buf, "\t}")
		fmt.Fprintln(&buf, "")
		fmt.Fprintln(&buf, "\treturn gorm.Expr(expr, args...)")
		fmt.Fprintln(&buf, "}")
		fmt.Fprintln(&buf)
	}

	// Generate diff functions for each struct
	for _, structInfo := range g.Structs {
		code, err := g.GenerateDiffFunction(structInfo)
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

// loadDiffTemplate loads the diff function template from embedded content
func (g *DiffGenerator) loadDiffTemplate() (*template.Template, error) {
	// Create template funcs
	funcMap := template.FuncMap{
		"trimStar": func(s string) string {
			return strings.TrimPrefix(s, "*")
		},
		"hasPrefix": func(s, prefix string) bool {
			return strings.HasPrefix(s, prefix)
		},
		"hasSuffix": func(s, suffix string) bool {
			return strings.HasSuffix(s, suffix)
		},
		"getColumnName": func(fieldName, tagStr string) string {
			return g.extractColumnName(fieldName, tagStr)
		},
		"isEmptyJSON": isEmptyJSON,
	}

	// Parse the embedded template
	tmpl, err := template.New("diff").Funcs(funcMap).Parse(diffFunctionTemplate)
	if err != nil {
		return nil, fmt.Errorf("error parsing embedded template: %v", err)
	}

	return tmpl, nil
}

// isEmptyJSON checks if a JSON string represents an empty object or array
func isEmptyJSON(jsonStr string) bool {
	trimmed := strings.TrimSpace(jsonStr)
	return trimmed == "{}" || trimmed == "[]" || trimmed == "null"
}

// GenerateDiffFunction generates a diff function for a struct
func (g *DiffGenerator) GenerateDiffFunction(structInfo StructInfo) (string, error) {
	// Load template
	tmpl, err := g.loadDiffTemplate()
	if err != nil {
		return "", err
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
		g.extractImports(node.Imports)
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
