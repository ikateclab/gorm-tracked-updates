package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// ServiceData represents a JSON field that might be empty
type ServiceData struct {
	Status map[string]interface{} `json:"status,omitempty"`
}

// Service represents a service with JSON fields
type Service struct {
	ID   uuid.UUID    `json:"id"`
	Name string       `json:"name"`
	Data *ServiceData `json:"data" gorm:"type:jsonb;serializer:json"`
}

// isEmptyJSON checks if a JSON string represents an empty object or array
func isEmptyJSON(jsonStr string) bool {
	trimmed := strings.TrimSpace(jsonStr)
	return trimmed == "{}" || trimmed == "[]" || trimmed == "null"
}

// Diff compares two Service instances and returns differences
// This demonstrates the empty JSON prevention logic
func (a *Service) Diff(b *Service) map[string]interface{} {
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare ID
	if a.ID != b.ID {
		diff["id"] = b.ID
	}

	// Compare Name
	if a.Name != b.Name {
		diff["name"] = b.Name
	}

	// Compare Data (JSON field with empty prevention)
	if a.Data != b.Data {
		jsonValue, err := json.Marshal(b.Data)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			// Only add to diff if JSON is not empty
			diff["data"] = fmt.Sprintf("JSON_MERGE(%s)", string(jsonValue))
		} else if err != nil {
			// Fallback to regular assignment if JSON marshaling fails
			diff["data"] = b.Data
		}
		// Skip adding to diff if JSON is empty (no-op update)
	}

	return diff
}

func main() {
	fmt.Println("🚀 Empty JSON Prevention Demo")
	fmt.Println("=============================")

	// Create two services with empty JSON data
	service1 := &Service{
		ID:   uuid.New(),
		Name: "Test Service",
		Data: &ServiceData{}, // Empty struct
	}

	service2 := &Service{
		ID:   service1.ID,
		Name: service1.Name,
		Data: &ServiceData{}, // Same empty struct
	}

	fmt.Println("\n📋 Test 1: Comparing services with empty JSON data")
	fmt.Printf("Service 1 Data: %+v\n", service1.Data)
	fmt.Printf("Service 2 Data: %+v\n", service2.Data)

	diff1 := service1.Diff(service2)
	fmt.Printf("Diff result: %v\n", diff1)
	fmt.Printf("✅ Empty JSON prevented: %t\n", len(diff1) == 0)

	// Test with actual data change
	service3 := &Service{
		ID:   service1.ID,
		Name: "Updated Service", // Different name
		Data: &ServiceData{},    // Still empty
	}

	fmt.Println("\n📋 Test 2: Comparing services with name change but empty JSON")
	fmt.Printf("Service 1 Name: %s\n", service1.Name)
	fmt.Printf("Service 3 Name: %s\n", service3.Name)
	fmt.Printf("Both have empty JSON data\n")

	diff2 := service1.Diff(service3)
	fmt.Printf("Diff result: %v\n", diff2)
	fmt.Printf("✅ Only name changed, empty JSON prevented: %t\n", len(diff2) == 1 && diff2["name"] != nil)

	// Test with non-empty JSON
	service4 := &Service{
		ID:   service1.ID,
		Name: service1.Name,
		Data: &ServiceData{
			Status: map[string]interface{}{
				"active": true,
				"count":  42,
			},
		},
	}

	fmt.Println("\n📋 Test 3: Comparing services with non-empty JSON data")
	fmt.Printf("Service 1 Data: %+v\n", service1.Data)
	fmt.Printf("Service 4 Data: %+v\n", service4.Data)

	diff3 := service1.Diff(service4)
	fmt.Printf("Diff result: %v\n", diff3)
	fmt.Printf("✅ Non-empty JSON included in diff: %t\n", len(diff3) == 1 && diff3["data"] != nil)

	// Demonstrate what the JSON looks like when marshaled
	fmt.Println("\n🔍 JSON Marshaling Examples:")

	emptyData, _ := json.Marshal(service1.Data)
	fmt.Printf("Empty ServiceData JSON: %s\n", string(emptyData))
	fmt.Printf("Is empty: %t\n", isEmptyJSON(string(emptyData)))

	nonEmptyData, _ := json.Marshal(service4.Data)
	fmt.Printf("Non-empty ServiceData JSON: %s\n", string(nonEmptyData))
	fmt.Printf("Is empty: %t\n", isEmptyJSON(string(nonEmptyData)))

	fmt.Println("\n🎯 Summary:")
	fmt.Println("- Empty JSON objects ({}) are prevented from creating diff entries")
	fmt.Println("- This avoids no-op SQL updates like: UPDATE services SET data=data || '{}'")
	fmt.Println("- Non-empty JSON objects are properly included in diffs")
	fmt.Println("- Performance is improved by avoiding unnecessary database operations")
}
