package main

import (
	"fmt"
	"log"

	"gorm-tracked-updates/pkg/clonegen"
	"gorm-tracked-updates/pkg/diffgen"
)

func main() {
	fmt.Println("🚀 Multi-File Demo - Generating from Multiple Files")
	fmt.Println("=" + fmt.Sprintf("%50s", ""))

	// Step 1: Generate diff methods from multiple files
	fmt.Println("\n📝 Step 1: Generating diff methods from multiple files...")
	diffGenerator := diffgen.New()

	err := diffGenerator.ParseDirectory("../multi-file")
	if err != nil {
		log.Fatalf("❌ Error parsing multi-file directory: %v", err)
	}

	fmt.Printf("✅ Successfully parsed %d structs from multiple files\n", len(diffGenerator.Structs))

	// Display found structs
	for _, structInfo := range diffGenerator.Structs {
		fmt.Printf("   - %s (%d fields)\n", structInfo.Name, len(structInfo.Fields))
	}

	// Generate and write diff methods
	err = diffGenerator.WriteToPackageDir("../multi-file")
	if err != nil {
		log.Fatalf("❌ Error writing diff methods: %v", err)
	}

	fmt.Println("✅ Generated diff methods written to 'examples/multi-file/diff.go'")

	// Step 2: Generate clone methods from multiple files
	fmt.Println("\n🔧 Step 2: Generating clone methods from multiple files...")
	cloneGenerator := clonegen.New()

	err = cloneGenerator.ParseDirectory("../multi-file")
	if err != nil {
		log.Fatalf("❌ Error parsing multi-file directory: %v", err)
	}

	fmt.Printf("✅ Successfully parsed %d structs from multiple files\n", len(cloneGenerator.Structs))

	// Generate and write clone methods
	err = cloneGenerator.WriteToPackageDir("../multi-file")
	if err != nil {
		log.Fatalf("❌ Error writing clone methods: %v", err)
	}

	fmt.Println("✅ Generated clone methods written to 'examples/multi-file/clone.go'")

	// Step 3: Summary
	fmt.Println("\n📊 Multi-File Generation Summary:")
	fmt.Printf("   - Diff methods: %d structs processed\n", len(diffGenerator.Structs))
	fmt.Printf("   - Clone methods: %d structs processed\n", len(cloneGenerator.Structs))

	fmt.Println("\n🔍 Generated files:")
	fmt.Println("   - examples/multi-file/diff.go (diff methods)")
	fmt.Println("   - examples/multi-file/clone.go (clone methods)")

	fmt.Println("\n🎯 Key benefits of multi-file support:")
	fmt.Println("   - Structs can be organized in separate files")
	fmt.Println("   - Cross-file struct references are handled correctly")
	fmt.Println("   - All methods are generated in single clone.go and diff.go files")
	fmt.Println("   - Package-level organization is maintained")

	fmt.Println("\n🧪 Usage examples:")
	fmt.Println("   - person := multifile.Person{...}")
	fmt.Println("   - cloned := person.Clone()")
	fmt.Println("   - changes := person.Diff(modifiedPerson)")

	fmt.Println("\n🎯 Multi-file generation is working correctly!")
}
