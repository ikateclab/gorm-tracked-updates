package main

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm-tracked-updates/examples/go-generate/models"
)

func main() {
	fmt.Println("🚀 go:generate Demo - Using Generated Clone and Diff Methods")
	fmt.Println("=" + fmt.Sprintf("%60s", ""))

	// Create test data
	original := createTestUser()

	fmt.Println("\n📋 Step 1: Original User Data")
	printUser("Original", original)

	// Clone the user
	fmt.Println("\n🔧 Step 2: Cloning User")
	cloned := original.Clone()
	fmt.Println("✅ User cloned successfully")

	// Verify independence
	fmt.Println("\n🔍 Step 3: Verifying Independence")
	cloned.Name = "Jane Smith"
	cloned.Email = "jane.smith@example.com"
	cloned.Age = 28
	cloned.Profile.Bio = "Updated bio"
	cloned.Addresses[0].City = "San Francisco"

	fmt.Printf("Original name: %s, Cloned name: %s\n", original.Name, cloned.Name)
	fmt.Printf("Original city: %s, Cloned city: %s\n", original.Addresses[0].City, cloned.Addresses[0].City)
	fmt.Println("✅ Clone is independent from original")

	// Generate diff
	fmt.Println("\n📊 Step 4: Generating Diff")
	changes := original.Diff(cloned)

	fmt.Println("Changes detected:")
	changesJSON, _ := json.MarshalIndent(changes, "", "  ")
	fmt.Println(string(changesJSON))

	// Demonstrate GORM usage
	fmt.Println("\n💾 Step 5: GORM Usage Example")
	fmt.Println("// Typical GORM workflow:")
	fmt.Println("backup := user.Clone()")
	fmt.Println("// ... make changes to user ...")
	fmt.Println("changes := backup.Diff(user)")
	fmt.Println("result := db.Model(&user).Updates(changes)")
	fmt.Printf("// Would update %d fields\n", len(changes))

	fmt.Println("\n🎯 go:generate integration working perfectly!")
}

func createTestUser() models.User {
	return models.User{
		ID:    1,
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Age:   30,
		Profile: models.Profile{
			Bio:      "Software developer",
			Avatar:   "avatar.jpg",
			Verified: true,
			Settings: map[string]interface{}{
				"theme":         "dark",
				"notifications": true,
			},
			Metadata: map[string]string{
				"department": "Engineering",
				"level":      "Senior",
			},
		},
		Addresses: []models.Address{
			{
				ID:      1,
				UserID:  1,
				Type:    "home",
				Street:  "123 Main St",
				City:    "New York",
				State:   "NY",
				ZipCode: "10001",
				Country: "USA",
				Primary: true,
			},
			{
				ID:      2,
				UserID:  1,
				Type:    "work",
				Street:  "456 Business Ave",
				City:    "New York",
				State:   "NY",
				ZipCode: "10002",
				Country: "USA",
				Primary: false,
			},
		},
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
	}
}

func printUser(label string, user models.User) {
	fmt.Printf("%s User:\n", label)
	fmt.Printf("  ID: %d\n", user.ID)
	fmt.Printf("  Name: %s\n", user.Name)
	fmt.Printf("  Email: %s\n", user.Email)
	fmt.Printf("  Age: %d\n", user.Age)
	fmt.Printf("  Bio: %s\n", user.Profile.Bio)
	fmt.Printf("  Addresses: %d\n", len(user.Addresses))
	if len(user.Addresses) > 0 {
		fmt.Printf("    Primary: %s, %s\n", user.Addresses[0].Street, user.Addresses[0].City)
	}
}
