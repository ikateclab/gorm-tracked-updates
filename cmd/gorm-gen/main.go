package main

import (
	"flag"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/ikateclab/gorm-tracked-updates/pkg/clonegen"
	"github.com/ikateclab/gorm-tracked-updates/pkg/diffgen"
)

func main() {
	var (
		packageDir = flag.String("package", ".", "Package directory to scan for structs")
		types      = flag.String("types", "clone,diff", "Types to generate (clone,diff)")
		output     = flag.String("output", "", "Output directory (defaults to package directory)")
		help       = flag.Bool("help", false, "Show help")
	)
	flag.Parse()

	if *help {
		printUsage()
		return
	}

	// Default output to package directory if not specified
	if *output == "" {
		*output = *packageDir
	}

	// Parse types to generate
	generateTypes := strings.Split(*types, ",")
	generateClone := contains(generateTypes, "clone")
	generateDiff := contains(generateTypes, "diff")

	if !generateClone && !generateDiff {
		log.Fatal("At least one of 'clone' or 'diff' must be specified in -types")
	}

	// Convert to absolute paths
	absPackageDir, err := filepath.Abs(*packageDir)
	if err != nil {
		log.Fatalf("Error resolving package directory: %v", err)
	}

	absOutputDir, err := filepath.Abs(*output)
	if err != nil {
		log.Fatalf("Error resolving output directory: %v", err)
	}

	fmt.Printf("🚀 GORM Code Generator\n")
	fmt.Printf("📁 Package: %s\n", absPackageDir)
	fmt.Printf("📤 Output: %s\n", absOutputDir)
	fmt.Printf("🔧 Types: %s\n", *types)
	fmt.Println()

	// Generate clone methods
	if generateClone {
		fmt.Println("🔧 Generating clone methods...")
		cloneGenerator := clonegen.New()

		err := cloneGenerator.ParseDirectory(absPackageDir)
		if err != nil {
			log.Fatalf("Error parsing directory for clone generation: %v", err)
		}

		if len(cloneGenerator.Structs) == 0 {
			fmt.Println("⚠️  No structs found for clone generation")
		} else {
			err = cloneGenerator.WriteToPackageDir(absOutputDir)
			if err != nil {
				log.Fatalf("Error writing clone methods: %v", err)
			}

			fmt.Printf("✅ Generated clone methods for %d structs\n", len(cloneGenerator.Structs))
			fmt.Printf("   Written to: %s/clone.go\n", absOutputDir)
		}
	}

	// Generate diff methods
	if generateDiff {
		fmt.Println("📝 Generating diff methods...")
		diffGenerator := diffgen.New()

		err := diffGenerator.ParseDirectory(absPackageDir)
		if err != nil {
			log.Fatalf("Error parsing directory for diff generation: %v", err)
		}

		if len(diffGenerator.Structs) == 0 {
			fmt.Println("⚠️  No structs found for diff generation")
		} else {
			err = diffGenerator.WriteToPackageDir(absOutputDir)
			if err != nil {
				log.Fatalf("Error writing diff methods: %v", err)
			}

			fmt.Printf("✅ Generated diff methods for %d structs\n", len(diffGenerator.Structs))
			fmt.Printf("   Written to: %s/diff.go\n", absOutputDir)
		}
	}

	fmt.Println("\n🎯 Code generation completed successfully!")
}

func printUsage() {
	fmt.Println("GORM Code Generator")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gorm-gen [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	flag.PrintDefaults()
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gorm-gen                                    # Generate both clone and diff in current directory")
	fmt.Println("  gorm-gen -types=clone                       # Generate only clone methods")
	fmt.Println("  gorm-gen -types=diff                        # Generate only diff methods")
	fmt.Println("  gorm-gen -package=./models                  # Generate for models directory")
	fmt.Println("  gorm-gen -package=./models -output=./gen    # Generate to different output directory")
	fmt.Println()
	fmt.Println("go:generate usage:")
	fmt.Println("  //go:generate gorm-gen")
	fmt.Println("  //go:generate gorm-gen -types=clone")
	fmt.Println("  //go:generate gorm-gen -package=./models")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if strings.TrimSpace(s) == item {
			return true
		}
	}
	return false
}
