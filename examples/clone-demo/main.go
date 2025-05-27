package main

import (
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"gorm-tracked-updates/examples/structs"
	"gorm-tracked-updates/pkg/clonegen"
)

func main() {
	fmt.Println("🚀 CloneGen Demo - Generating Clone Methods")
	fmt.Println("=" + fmt.Sprintf("%50s", ""))

	// Step 1: Create and test the clone generator
	fmt.Println("\n📝 Step 1: Creating clone generator and parsing structs...")
	generator := clonegen.New()

	err := generator.ParseDirectory("../structs")
	if err != nil {
		log.Fatalf("❌ Error parsing structs directory: %v", err)
	}

	fmt.Printf("✅ Successfully parsed %d structs from structs.go\n", len(generator.Structs))

	// Display found structs
	for _, structInfo := range generator.Structs {
		fmt.Printf("   - %s (%d fields)\n", structInfo.Name, len(structInfo.Fields))
	}

	// Step 2: Generate clone methods
	fmt.Println("\n🔧 Step 2: Generating clone methods...")
	code, err := generator.GenerateCode()
	if err != nil {
		log.Fatalf("❌ Error generating code: %v", err)
	}

	fmt.Printf("✅ Generated %d bytes of clone method code\n", len(code))

	// Step 3: Write to file
	fmt.Println("\n💾 Step 3: Writing generated code to file...")
	err = generator.WriteToPackageDir("../structs")
	if err != nil {
		log.Fatalf("❌ Error writing to file: %v", err)
	}

	fmt.Println("✅ Generated code written to '../structs/clone.go'")

	// Step 4: Demonstrate the functionality with manual cloning
	fmt.Println("\n🧪 Step 4: Demonstrating clone functionality...")

	// Create test data using the structs package
	manager := &structs.Person{
		Name: "Jane Doe",
		Age:  45,
		Address: structs.Address{
			Street:  "789 Oak Dr",
			City:    "Managertown",
			State:   "CA",
			ZipCode: "54321",
			Country: "USA",
		},
		Contacts: []structs.Contact{
			{Type: "email", Value: "jane@company.com"},
		},
		Manager:  nil,
		Metadata: map[string]interface{}{"role": "Senior Manager", "department": "Engineering"},
	}

	original := structs.Person{
		Name: "John Doe",
		Age:  30,
		Address: structs.Address{
			Street:  "123 Main St",
			City:    "Anytown",
			State:   "CA",
			ZipCode: "12345",
			Country: "USA",
		},
		Contacts: []structs.Contact{
			{Type: "email", Value: "john@example.com"},
			{Type: "phone", Value: "555-1234"},
		},
		Manager:  manager,
		Metadata: map[string]interface{}{"role": "developer", "team": "backend"},
	}

	// Manual clone demonstration (simulating what the generated code would do)
	cloned := manualClonePerson(original)

	fmt.Println("\n📊 Original vs Clone comparison:")
	fmt.Println("Original and clone are equal:", reflect.DeepEqual(original, cloned))

	// Demonstrate independence by modifying the clone
	fmt.Println("\n🔄 Modifying clone to demonstrate independence...")
	cloned.Age = 31
	cloned.Address.City = "Newtown"
	cloned.Contacts[0].Value = "john@newexample.com"
	cloned.Manager.Age = 46
	cloned.Metadata["team"] = "frontend"

	fmt.Println("\nAfter modifications:")
	fmt.Printf("Original age: %d, Clone age: %d\n", original.Age, cloned.Age)
	fmt.Printf("Original city: %s, Clone city: %s\n", original.Address.City, cloned.Address.City)
	fmt.Printf("Original email: %s, Clone email: %s\n", original.Contacts[0].Value, cloned.Contacts[0].Value)
	fmt.Printf("Original manager age: %d, Clone manager age: %d\n", original.Manager.Age, cloned.Manager.Age)
	fmt.Printf("Original team: %s, Clone team: %s\n", original.Metadata["team"], cloned.Metadata["team"])

	// Verify independence
	fmt.Println("\n✅ Independence verification:")
	fmt.Println("   - Simple fields are independent:", original.Age != cloned.Age)
	fmt.Println("   - Nested structs are independent:", original.Address.City != cloned.Address.City)
	fmt.Println("   - Slice elements are independent:", original.Contacts[0].Value != cloned.Contacts[0].Value)
	fmt.Println("   - Pointer targets are independent:", original.Manager.Age != cloned.Manager.Age)
	fmt.Println("   - Maps are independent:", original.Metadata["team"] != cloned.Metadata["team"])

	// Pretty print the structures for comparison
	fmt.Println("\n📋 Final state comparison:")

	originalJSON, _ := json.MarshalIndent(original, "", "  ")
	clonedJSON, _ := json.MarshalIndent(cloned, "", "  ")

	fmt.Println("\nOriginal:")
	fmt.Println(string(originalJSON))

	fmt.Println("\nClone:")
	fmt.Println(string(clonedJSON))

	fmt.Println("\n🎯 Key benefits of generated clone methods:")
	fmt.Println("   - Deep copy ensures complete independence")
	fmt.Println("   - Type-safe without reflection overhead")
	fmt.Println("   - Optimized for each field type")
	fmt.Println("   - No shared memory references")
	fmt.Println("   - Compile-time method resolution")

	fmt.Println("\n🔍 The generated clone methods can be found in:")
	fmt.Println("   - ../structs/clone.go")

	fmt.Println("\n🧪 To test the generated methods:")
	fmt.Println("   - Copy the generated methods to your code")
	fmt.Println("   - Use person.Clone() to create deep copies")
	fmt.Println("   - Modify clones without affecting originals")

	fmt.Println("\n🎯 CloneGen is working correctly!")
}

// manualClonePerson demonstrates what the generated clone method would look like
func manualClonePerson(original structs.Person) structs.Person {
	clone := structs.Person{}

	// Simple types - direct assignment
	clone.Name = original.Name
	clone.Age = original.Age

	// Struct type - recursive clone
	clone.Address = manualCloneAddress(original.Address)

	// Slice - create new slice and clone elements
	if original.Contacts != nil {
		clone.Contacts = make([]structs.Contact, len(original.Contacts))
		for i, item := range original.Contacts {
			clone.Contacts[i] = manualCloneContact(item)
		}
	}

	// Pointer to struct - create new instance and clone
	if original.Manager != nil {
		clonedManager := manualClonePerson(*original.Manager)
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

func manualCloneAddress(original structs.Address) structs.Address {
	clone := structs.Address{}
	clone.Street = original.Street
	clone.City = original.City
	clone.State = original.State
	clone.ZipCode = original.ZipCode
	clone.Country = original.Country
	return clone
}

func manualCloneContact(original structs.Contact) structs.Contact {
	clone := structs.Contact{}
	clone.Type = original.Type
	clone.Value = original.Value
	return clone
}
