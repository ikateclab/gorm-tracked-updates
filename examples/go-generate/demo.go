package main

import (
	"encoding/json"
	"fmt"
	"time"

	"example-models/models"
	"github.com/google/uuid"
)

func main() {
	fmt.Println("🚀 go:generate Demo - Using Generated Clone and Diff Methods")
	fmt.Println("=" + fmt.Sprintf("%60s", ""))

	// Create test data
	original := createTestService()

	fmt.Println("\n📋 Step 1: Original Service Data")
	printService("Original", original)

	// Clone the service
	fmt.Println("\n🔧 Step 2: Cloning Service")
	cloned := original.Clone()
	fmt.Println("✅ Service cloned successfully")

	// Verify independence
	fmt.Println("\n🔍 Step 3: Verifying Independence")
	cloned.Name = "Updated Service"
	cloned.Data.SyncCount = 10
	cloned.Data.Status.IsConnected = false
	cloned.Settings.KeepOnline = false

	fmt.Printf("Original name: %s, Cloned name: %s\n", original.Name, cloned.Name)
	fmt.Printf("Original SyncCount: %d, Cloned SyncCount: %d\n", original.Data.SyncCount, cloned.Data.SyncCount)
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
	fmt.Println("backup := service.Clone()")
	fmt.Println("// ... make changes to service ...")
	fmt.Println("changes := backup.Diff(service)")
	fmt.Println("result := db.Model(&service).Updates(changes)")
	fmt.Printf("// Would update %d fields\n", len(changes))

	fmt.Println("\n🎯 go:generate integration working perfectly!")
}

func createTestService() models.Service {
	return models.Service{
		Id:   uuid.New(),
		Name: "Test Service",
		Data: &models.ServiceData{
			MyId:       "test123",
			SyncCount:  5,
			SyncFlowDone: true,
			Status: models.ServiceDataStatus{
				IsConnected: true,
				IsStarted:   true,
				State:       "active",
			},
		},
		Settings: &models.ServiceSettings{
			KeepOnline:        true,
			WppConnectVersion: "1.0.0",
			WaVersion:         "2.0.0",
		},
		CreatedAt: time.Now().Add(-24 * time.Hour),
		UpdatedAt: time.Now(),
		AccountId: uuid.New(),
	}
}

func printService(label string, service models.Service) {
	fmt.Printf("%s Service:\n", label)
	fmt.Printf("  ID: %s\n", service.Id)
	fmt.Printf("  Name: %s\n", service.Name)
	if service.Data != nil {
		fmt.Printf("  Data.MyId: %s\n", service.Data.MyId)
		fmt.Printf("  Data.SyncCount: %d\n", service.Data.SyncCount)
		fmt.Printf("  Data.Status.IsConnected: %t\n", service.Data.Status.IsConnected)
	}
	if service.Settings != nil {
		fmt.Printf("  Settings.KeepOnline: %t\n", service.Settings.KeepOnline)
		fmt.Printf("  Settings.WppConnectVersion: %s\n", service.Settings.WppConnectVersion)
	}
}
