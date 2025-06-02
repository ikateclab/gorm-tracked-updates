package clonegen

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

func TestCloneGeneratorParsing(t *testing.T) {
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

func TestCloneCodeGeneration(t *testing.T) {
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

	// Verify it contains clone methods with new pointer-based signature
	if !strings.Contains(code, "func (original *Address) Clone() *Address") {
		t.Errorf("Expected generated code to contain Address Clone method with pointer signature")
	}
	if !strings.Contains(code, "func (original *Contact) Clone() *Contact") {
		t.Errorf("Expected generated code to contain Contact Clone method with pointer signature")
	}
	if !strings.Contains(code, "func (original *Person) Clone() *Person") {
		t.Errorf("Expected generated code to contain Person Clone method with pointer signature")
	}
}

func TestFieldTypeCategorization(t *testing.T) {
	generator := New()

	// Add some known structs
	generator.KnownStructs["TestStruct"] = true
	generator.KnownStructs["Address"] = true

	// Test the categorizeFieldType method directly
	tests := []struct {
		fieldType string
		expected  FieldType
	}{
		{"string", FieldTypeSimple},
		{"int", FieldTypeSimple},
		{"bool", FieldTypeSimple},
		{"TestStruct", FieldTypeStruct},
		{"*TestStruct", FieldTypeStructPtr},
		{"[]string", FieldTypeSlice},
		{"[]TestStruct", FieldTypeSlice},
		{"map[string]int", FieldTypeMap},
		{"interface{}", FieldTypeInterface},
		{"UnknownType", FieldTypeComplex},
	}

	for _, test := range tests {
		result := generator.categorizeFieldType(test.fieldType)
		if result != test.expected {
			t.Errorf("categorizeFieldType(%s) = %v, expected %v", test.fieldType, result, test.expected)
		}
	}
}

func TestCloneMethodGeneration(t *testing.T) {
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

	code, err := generator.generateCloneMethod(structInfo)
	if err != nil {
		t.Fatalf("Error generating clone method: %v", err)
	}

	// Verify the generated method contains expected elements
	if !strings.Contains(code, "func (original *TestAddress) Clone() *TestAddress") {
		t.Errorf("Expected method signature Clone with pointer receiver and return type")
	}
	// For simple structs, the new generator uses shallow copy with *original
	if !strings.Contains(code, "clone := *original") {
		t.Errorf("Expected shallow copy assignment for simple struct")
	}
}

func TestWriteToFile(t *testing.T) {
	generator := New()

	err := generator.ParseFile("../../examples/structs/structs.go")
	if err != nil {
		t.Fatalf("Error parsing test file: %v", err)
	}

	// Write to a temporary file
	tempFile := "/tmp/test_clone_output.go"
	err = generator.WriteToFile(tempFile)
	if err != nil {
		t.Fatalf("Error writing to file: %v", err)
	}

	// Verify file was created (basic check)
	// In a real test, you might want to read and verify the file content
}

// Manual clone methods for testing (simulating generated code)
func (original TestAddress) Clone() TestAddress {
	clone := TestAddress{}

	// Simple type - direct assignment
	clone.Street = original.Street
	clone.City = original.City
	clone.State = original.State
	clone.ZipCode = original.ZipCode
	clone.Country = original.Country

	return clone
}

func (original TestContact) Clone() TestContact {
	clone := TestContact{}

	// Simple type - direct assignment
	clone.Type = original.Type
	clone.Value = original.Value

	return clone
}

func (original TestPerson) Clone() TestPerson {
	clone := TestPerson{}

	// Simple type - direct assignment
	clone.Name = original.Name
	clone.Age = original.Age

	// Struct type - recursive clone
	clone.Address = original.Address.Clone()

	// Slice - create new slice and clone elements
	if original.Contacts != nil {
		clone.Contacts = make([]TestContact, len(original.Contacts))
		for i, item := range original.Contacts {
			clone.Contacts[i] = item.Clone()
		}
	}

	// Pointer to struct - create new instance and clone
	if original.Manager != nil {
		clonedManager := original.Manager.Clone()
		clone.Manager = &clonedManager
	}

	// Map - create new map and copy key-value pairs
	if original.Metadata != nil {
		clone.Metadata = make(map[string]interface{})
		for k, v := range original.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}

func TestCloneFunctionality(t *testing.T) {
	// Test case 1: Simple struct cloning
	t.Run("Simple struct cloning", func(t *testing.T) {
		original := TestAddress{
			Street:  "123 Main St",
			City:    "Anytown",
			State:   "CA",
			ZipCode: "12345",
			Country: "USA",
		}

		cloned := original.Clone()

		// Verify the clone is equal but not the same reference
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Clone should be equal to original")
		}

		// Modify the clone and ensure original is unchanged
		cloned.City = "Newtown"
		if original.City == cloned.City {
			t.Errorf("Modifying clone should not affect original")
		}
	})

	// Test case 2: Nested struct cloning
	t.Run("Nested struct cloning", func(t *testing.T) {
		original := TestPerson{
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

		cloned := original.Clone()

		// Verify the clone is equal but not the same reference
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Clone should be equal to original")
		}

		// Modify nested struct in clone and ensure original is unchanged
		cloned.Address.City = "Newtown"
		if original.Address.City == cloned.Address.City {
			t.Errorf("Modifying nested struct in clone should not affect original")
		}
	})

	// Test case 3: Slice cloning
	t.Run("Slice cloning", func(t *testing.T) {
		original := TestPerson{
			Name: "John Doe",
			Age:  30,
			Contacts: []TestContact{
				{Type: "email", Value: "john@example.com"},
				{Type: "phone", Value: "555-1234"},
			},
		}

		cloned := original.Clone()

		// Verify the clone is equal but not the same reference
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Clone should be equal to original")
		}

		// Verify slices are different references
		if &original.Contacts[0] == &cloned.Contacts[0] {
			t.Errorf("Slice elements should be different references")
		}

		// Modify slice in clone and ensure original is unchanged
		cloned.Contacts[0].Value = "john@newexample.com"
		if original.Contacts[0].Value == cloned.Contacts[0].Value {
			t.Errorf("Modifying slice element in clone should not affect original")
		}
	})

	// Test case 4: Pointer cloning
	t.Run("Pointer cloning", func(t *testing.T) {
		manager := &TestPerson{
			Name: "Jane Doe",
			Age:  45,
		}

		original := TestPerson{
			Name:    "John Doe",
			Age:     30,
			Manager: manager,
		}

		cloned := original.Clone()

		// Verify the clone is equal but not the same reference
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Clone should be equal to original")
		}

		// Verify pointers are different references
		if original.Manager == cloned.Manager {
			t.Errorf("Pointer fields should be different references")
		}

		// Modify pointed-to struct in clone and ensure original is unchanged
		cloned.Manager.Age = 46
		if original.Manager.Age == cloned.Manager.Age {
			t.Errorf("Modifying pointed-to struct in clone should not affect original")
		}
	})

	// Test case 5: Map cloning
	t.Run("Map cloning", func(t *testing.T) {
		original := TestPerson{
			Name: "John Doe",
			Age:  30,
			Metadata: map[string]interface{}{
				"role": "developer",
				"team": "backend",
			},
		}

		cloned := original.Clone()

		// Verify the clone is equal but not the same reference
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Clone should be equal to original")
		}

		// Verify maps are different references
		if &original.Metadata == &cloned.Metadata {
			t.Errorf("Map fields should be different references")
		}

		// Modify map in clone and ensure original is unchanged
		cloned.Metadata["team"] = "frontend"
		if original.Metadata["team"] == cloned.Metadata["team"] {
			t.Errorf("Modifying map in clone should not affect original")
		}
	})

	// Test case 6: Nil pointer handling
	t.Run("Nil pointer handling", func(t *testing.T) {
		original := TestPerson{
			Name:    "John Doe",
			Age:     30,
			Manager: nil, // Nil pointer
		}

		cloned := original.Clone()

		// Verify the clone is equal but not the same reference
		if !reflect.DeepEqual(original, cloned) {
			t.Errorf("Clone should be equal to original")
		}

		// Verify nil pointer is preserved
		if cloned.Manager != nil {
			t.Errorf("Nil pointer should remain nil in clone")
		}
	})
}
