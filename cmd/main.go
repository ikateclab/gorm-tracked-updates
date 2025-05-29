package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/ikateclab/gorm-tracked-updates/pkg/clonegen"
	"github.com/ikateclab/gorm-tracked-updates/pkg/diffgen"
)

func main() {
	fmt.Println("🚀 GORM Tracked Updates - Code Generators")
	fmt.Println(strings.Repeat("=", 50))

	// Generate diff functions
	fmt.Println("\n📝 Generating diff functions...")
	diffGenerator := diffgen.New()

	err := diffGenerator.ParseDirectory("examples/structs")
	if err != nil {
		log.Fatalf("Error parsing examples/structs directory for diff generation: %v", err)
	}

	diffCode, err := diffGenerator.GenerateCode()
	if err != nil {
		log.Fatalf("Error generating diff code: %v", err)
	}

	err = diffGenerator.WriteToPackageDir("examples/structs")
	if err != nil {
		log.Fatalf("Error writing diff code to file: %v", err)
	}

	fmt.Printf("✅ Generated %d bytes of diff functions code\n", len(diffCode))
	fmt.Println("   Written to 'examples/structs/diff.go'")

	// Generate clone methods
	fmt.Println("\n🔧 Generating clone methods...")
	cloneGenerator := clonegen.New()

	err = cloneGenerator.ParseDirectory("examples/structs")
	if err != nil {
		log.Fatalf("Error parsing examples/structs directory for clone generation: %v", err)
	}

	cloneCode, err := cloneGenerator.GenerateCode()
	if err != nil {
		log.Fatalf("Error generating clone code: %v", err)
	}

	err = cloneGenerator.WriteToPackageDir("examples/structs")
	if err != nil {
		log.Fatalf("Error writing clone code to file: %v", err)
	}

	fmt.Printf("✅ Generated %d bytes of clone methods code\n", len(cloneCode))
	fmt.Println("   Written to 'examples/structs/clone.go'")

	// Summary
	fmt.Println("\n📊 Generation Summary:")
	fmt.Printf("   - Diff functions: %d structs processed\n", len(diffGenerator.Structs))
	fmt.Printf("   - Clone methods: %d structs processed\n", len(cloneGenerator.Structs))

	fmt.Println("\n🔍 Generated files:")
	fmt.Println("   - examples/structs/diff.go (diff methods)")
	fmt.Println("   - examples/structs/clone.go (clone methods)")

	fmt.Println("\n🧪 Example and demo files:")
	fmt.Println("   - examples/diff-demo/ (diff generator demo)")
	fmt.Println("   - examples/clone-demo/ (clone generator demo)")
	fmt.Println("   - examples/performance/ (performance benchmarks)")

	fmt.Println("\n🎯 Both generators are working correctly!")
	fmt.Println("\nTo run:")
	fmt.Println("   go run cmd/main.go")
	fmt.Println("   go run examples/diff-demo/main.go")
	fmt.Println("   go run examples/clone-demo/main.go")
	fmt.Println("   go test examples/performance/ -bench=.")
}
