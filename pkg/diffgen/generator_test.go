package diffgen

import (
	"reflect"
	"strings"
	"testing"
)

// Test structs
type TestAddress struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
}

type TestContact struct {
	Type  string
	Value string
}

type TestPerson struct {
	Name     string
	Age      int
	Address  TestAddress
	Contacts []TestContact
	Manager  *TestPerson
	Metadata map[string]interface{}
}

func TestDiffGeneratorParsing(t *testing.T) {
	generator := New()

	// Create a temporary test file content
	testFile := "../../examples/structs/structs.go"

	err := generator.ParseFile(testFile)
	if err != nil {
		t.Fatalf("Error parsing test file: %v", err)
	}

	// Verify structs were found
	if len(generator.Structs) == 0 {
		t.Errorf("Expected to find structs in test file")
	}

	// Verify specific structs
	structNames := make(map[string]bool)
	for _, s := range generator.Structs {
		structNames[s.Name] = true
	}

	expectedStructs := []string{"Address", "Contact", "Person"}
	for _, expected := range expectedStructs {
		if !structNames[expected] {
			t.Errorf("Expected to find struct %s", expected)
		}
	}
}

func TestDiffCodeGeneration(t *testing.T) {
	generator := New()

	err := generator.ParseFile("../../examples/structs/structs.go")
	if err != nil {
		t.Fatalf("Error parsing test file: %v", err)
	}

	code, err := generator.GenerateCode()
	if err != nil {
		t.Fatalf("Error generating code: %v", err)
	}

	// Verify code was generated
	if len(code) == 0 {
		t.Errorf("Expected generated code to be non-empty")
	}

	// Verify it contains diff methods with new signature
	if !strings.Contains(code, "func (a Address) Diff(") {
		t.Errorf("Expected generated code to contain Address Diff method")
	}
	if !strings.Contains(code, "func (a Contact) Diff(") {
		t.Errorf("Expected generated code to contain Contact Diff method")
	}
	if !strings.Contains(code, "func (a Person) Diff(") {
		t.Errorf("Expected generated code to contain Person Diff method")
	}
}

func TestFieldTypeCategorization(t *testing.T) {
	generator := New()

	// Add some known structs
	generator.KnownStructs["TestStruct"] = true
	generator.KnownStructs["Address"] = true

	// Note: This test would need access to AST expressions to test determineFieldType
	// For now, we'll test the basic functionality through the parsing process

	// Test that the generator can categorize field types correctly during parsing
	err := generator.ParseFile("../../examples/structs/structs.go")
	if err != nil {
		t.Fatalf("Error parsing test file: %v", err)
	}

	// Verify that different field types were detected
	foundSimple := false
	foundStruct := false
	foundSlice := false
	foundMap := false

	for _, structInfo := range generator.Structs {
		for _, field := range structInfo.Fields {
			switch field.FieldType {
			case FieldTypeSimple:
				foundSimple = true
			case FieldTypeStruct:
				foundStruct = true
			case FieldTypeSlice:
				foundSlice = true
			case FieldTypeMap:
				foundMap = true
			}
		}
	}

	if !foundSimple {
		t.Error("Expected to find simple field types")
	}
	if !foundStruct {
		t.Error("Expected to find struct field types")
	}
	if !foundSlice {
		t.Error("Expected to find slice field types")
	}
	if !foundMap {
		t.Error("Expected to find map field types")
	}
}

func TestDiffFunctionGeneration(t *testing.T) {
	generator := New()
	generator.KnownStructs["TestAddress"] = true

	structInfo := StructInfo{
		Name:    "TestAddress",
		Package: "main",
		Fields: []StructField{
			{Name: "Street", Type: "string", FieldType: FieldTypeSimple},
			{Name: "City", Type: "string", FieldType: FieldTypeSimple},
		},
	}

	code, err := generator.generateDiffFunction(structInfo)
	if err != nil {
		t.Fatalf("Error generating diff function: %v", err)
	}

	// Verify the generated function contains expected elements
	if !strings.Contains(code, "func (a TestAddress) Diff(") {
		t.Errorf("Expected method signature Diff")
	}
	if !strings.Contains(code, "a.Street != b.Street") {
		t.Errorf("Expected Street field comparison")
	}
	if !strings.Contains(code, "a.City != b.City") {
		t.Errorf("Expected City field comparison")
	}
}

func TestWriteToFile(t *testing.T) {
	generator := New()

	err := generator.ParseFile("../../examples/structs/structs.go")
	if err != nil {
		t.Fatalf("Error parsing test file: %v", err)
	}

	// Write to a temporary file
	tempFile := "/tmp/test_diff_output.go"
	err = generator.WriteToFile(tempFile)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Verify file was created (basic check)
	// In a real test, you might want to read and verify the file content
}

// Manual diff functions for testing (simulating generated code)
func (a TestAddress) Diff(b TestAddress) map[string]interface{} {
	diff := make(map[string]interface{})

	if a.Street != b.Street {
		diff["Street"] = b.Street
	}
	if a.City != b.City {
		diff["City"] = b.City
	}
	if a.State != b.State {
		diff["State"] = b.State
	}
	if a.ZipCode != b.ZipCode {
		diff["ZipCode"] = b.ZipCode
	}
	if a.Country != b.Country {
		diff["Country"] = b.Country
	}

	return diff
}

func (a TestContact) Diff(b TestContact) map[string]interface{} {
	diff := make(map[string]interface{})

	if a.Type != b.Type {
		diff["Type"] = b.Type
	}
	if a.Value != b.Value {
		diff["Value"] = b.Value
	}

	return diff
}

func (a TestPerson) Diff(b TestPerson) map[string]interface{} {
	diff := make(map[string]interface{})

	if a.Name != b.Name {
		diff["Name"] = b.Name
	}
	if a.Age != b.Age {
		diff["Age"] = b.Age
	}

	// Struct type comparison
	if !reflect.DeepEqual(a.Address, b.Address) {
		nestedDiff := a.Address.Diff(b.Address)
		if len(nestedDiff) > 0 {
			diff["Address"] = nestedDiff
		}
	}

	// Complex type comparison (slice)
	if !reflect.DeepEqual(a.Contacts, b.Contacts) {
		diff["Contacts"] = b.Contacts
	}

	// Pointer to struct comparison
	if !reflect.DeepEqual(a.Manager, b.Manager) {
		if a.Manager == nil || b.Manager == nil {
			diff["Manager"] = b.Manager
		} else {
			nestedDiff := (*a.Manager).Diff(*b.Manager)
			if len(nestedDiff) > 0 {
				diff["Manager"] = nestedDiff
			}
		}
	}

	// Map comparison
	if !reflect.DeepEqual(a.Metadata, b.Metadata) {
		diff["Metadata"] = b.Metadata
	}

	return diff
}

func TestDiffFunctionality(t *testing.T) {
	// Test case 1: Simple field changes
	t.Run("Simple field changes", func(t *testing.T) {
		addr1 := TestAddress{
			Street:  "123 Main St",
			City:    "Anytown",
			State:   "CA",
			ZipCode: "12345",
			Country: "USA",
		}

		addr2 := TestAddress{
			Street:  "123 Main St", // Same
			City:    "Newtown",     // Changed
			State:   "NY",          // Changed
			ZipCode: "12345",       // Same
			Country: "USA",         // Same
		}

		diff := addr1.Diff(addr2)

		// Should only contain changed fields
		expected := map[string]interface{}{
			"City":  "Newtown",
			"State": "NY",
		}

		if !reflect.DeepEqual(diff, expected) {
			t.Errorf("Expected %v, got %v", expected, diff)
		}
	})

	// Test case 2: No changes
	t.Run("No changes", func(t *testing.T) {
		addr1 := TestAddress{
			Street:  "123 Main St",
			City:    "Anytown",
			State:   "CA",
			ZipCode: "12345",
			Country: "USA",
		}

		addr2 := addr1 // Same

		diff := addr1.Diff(addr2)

		// Should be empty
		if len(diff) != 0 {
			t.Errorf("Expected empty diff, got %v", diff)
		}
	})

	// Test case 3: Nested struct changes
	t.Run("Nested struct changes", func(t *testing.T) {
		person1 := TestPerson{
			Name: "John Doe",
			Age:  30,
			Address: TestAddress{
				Street:  "123 Main St",
				City:    "Anytown",
				State:   "CA",
				ZipCode: "12345",
				Country: "USA",
			},
		}

		person2 := TestPerson{
			Name: "John Doe", // Same
			Age:  31,         // Changed
			Address: TestAddress{
				Street:  "123 Main St", // Same
				City:    "Newtown",     // Changed
				State:   "CA",          // Same
				ZipCode: "12345",       // Same
				Country: "USA",         // Same
			},
		}

		diff := person1.Diff(person2)

		// Should contain age change and nested address change
		if diff["Age"] != 31 {
			t.Errorf("Expected Age to be 31, got %v", diff["Age"])
		}

		addressDiff, ok := diff["Address"].(map[string]interface{})
		if !ok {
			t.Errorf("Expected Address diff to be a map")
		} else if addressDiff["City"] != "Newtown" {
			t.Errorf("Expected Address.City to be 'Newtown', got %v", addressDiff["City"])
		}
	})
}

func TestJSONFieldDetection(t *testing.T) {
	generator := New()

	testCases := []struct {
		name     string
		tag      string
		expected bool
	}{
		{
			name:     "Valid JSON tag",
			tag:      `gorm:"serializer:json"`,
			expected: true,
		},
		{
			name:     "Valid JSON tag with backticks",
			tag:      "`gorm:\"serializer:json\"`",
			expected: true,
		},
		{
			name:     "JSON tag with other options",
			tag:      `gorm:"column:settings;serializer:json"`,
			expected: true,
		},
		{
			name:     "Valid JSONB tag",
			tag:      `gorm:"type:jsonb"`,
			expected: true,
		},
		{
			name:     "JSONB tag with other options",
			tag:      `gorm:"type:jsonb;not null;default:'{}'"`,
			expected: true,
		},
		{
			name:     "JSONB tag with serializer",
			tag:      `gorm:"type:jsonb;not null;default:'{}';serializer:json"`,
			expected: true,
		},
		{
			name:     "No JSON tag",
			tag:      `gorm:"column:name"`,
			expected: false,
		},
		{
			name:     "Empty tag",
			tag:      "",
			expected: false,
		},
		{
			name:     "Different serializer",
			tag:      `gorm:"serializer:gob"`,
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := generator.isJSONField(tc.tag)
			if result != tc.expected {
				t.Errorf("Expected isJSONField(%q) = %v, got %v", tc.tag, tc.expected, result)
			}
		})
	}
}
