package models

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm/clause"
)

func TestNestedJSONDiffIntegration(t *testing.T) {
	// Test that nested JSON structures work correctly without nested gorm.Expr
	
	t.Run("ServiceData nested Status diff", func(t *testing.T) {
		// Create two ServiceData instances with different Status values
		data1 := &ServiceData{
			MyId:      "test123",
			SyncCount: 5,
			Status: ServiceDataStatus{
				IsConnected: false,
				IsStarting:  true,
				Mode:        "sync",
			},
		}

		data2 := &ServiceData{
			MyId:      "test123", // Same
			SyncCount: 10,        // Different
			Status: ServiceDataStatus{
				IsConnected: true,  // Different
				IsStarting:  true,  // Same
				Mode:        "idle", // Different
			},
		}

		diff := data2.Diff(data1) // new.Diff(old) semantics

		// Should have differences for SyncCount and Status
		if len(diff) != 2 {
			t.Errorf("Expected 2 differences, got %d: %v", len(diff), diff)
		}

		// Check SyncCount diff (new value from data2)
		if diff["syncCount"] != 10 {
			t.Errorf("Expected syncCount to be 10, got %v", diff["syncCount"])
		}

		// Check Status diff - should be a plain map, not gorm.Expr
		statusDiff, exists := diff["status"]
		if !exists {
			t.Fatal("Expected status diff to exist")
		}

		// Status diff should be a plain map[string]interface{}, not gorm.Expr
		statusMap, ok := statusDiff.(map[string]interface{})
		if !ok {
			t.Fatalf("Expected status diff to be map[string]interface{}, got %T: %v", statusDiff, statusDiff)
		}

		// Verify the nested status changes (new values from data2)
		expectedStatusChanges := map[string]interface{}{
			"isConnected": true,
			"mode":        "idle",
		}

		for key, expectedValue := range expectedStatusChanges {
			if statusMap[key] != expectedValue {
				t.Errorf("Expected status.%s to be %v, got %v", key, expectedValue, statusMap[key])
			}
		}

		// Verify IsStarting is NOT in the diff (since it's the same)
		if _, exists := statusMap["isStarting"]; exists {
			t.Error("Expected isStarting to NOT be in status diff since it's unchanged")
		}
	})

	t.Run("Service root JSONB field diff", func(t *testing.T) {
		// Test that the root Service.Data field uses gorm.Expr correctly
		now := time.Now()
		
		service1 := &Service{
			Id:   uuid.New(),
			Name: "Test Service",
			Data: &ServiceData{
				MyId:      "test123",
				SyncCount: 5,
				Status: ServiceDataStatus{
					IsConnected: false,
					Mode:        "sync",
				},
				LastSyncAt: &now,
			},
			Settings: &ServiceSettings{
				KeepOnline:        true,
				WppConnectVersion: "1.0.0",
			},
		}

		service2 := &Service{
			Id:   service1.Id, // Same
			Name: service1.Name, // Same
			Data: &ServiceData{
				MyId:      "test123", // Same
				SyncCount: 10,        // Different
				Status: ServiceDataStatus{
					IsConnected: true,  // Different
					Mode:        "idle", // Different
				},
				LastSyncAt: service1.Data.LastSyncAt, // Same
			},
			Settings: &ServiceSettings{
				KeepOnline:        false, // Different
				WppConnectVersion: "1.0.0", // Same
			},
		}

		diff := service2.Diff(service1) // new.Diff(old) semantics

		// Should have differences for Data and Settings
		if len(diff) != 2 {
			t.Errorf("Expected 2 differences, got %d: %v", len(diff), diff)
		}

		// Check that Data field uses gorm.Expr
		dataDiff, exists := diff["Data"]
		if !exists {
			t.Fatal("Expected Data diff to exist")
		}

		// Data diff should be a gorm.Expr (clause.Expr)
		if _, ok := dataDiff.(clause.Expr); !ok {
			t.Fatalf("Expected Data diff to be clause.Expr (gorm.Expr), got %T", dataDiff)
		}

		// Check that Settings field uses gorm.Expr
		settingsDiff, exists := diff["Settings"]
		if !exists {
			t.Fatal("Expected Settings diff to exist")
		}

		// Settings diff should be a gorm.Expr (clause.Expr)
		if _, ok := settingsDiff.(clause.Expr); !ok {
			t.Fatalf("Expected Settings diff to be clause.Expr (gorm.Expr), got %T", settingsDiff)
		}
	})

	t.Run("Verify no nested gorm.Expr in JSON content", func(t *testing.T) {
		// This test ensures that when we marshal the nested diff to JSON,
		// it doesn't contain nested gorm.Expr structures
		
		data1 := &ServiceData{
			SyncCount: 5,
			Status: ServiceDataStatus{
				IsConnected: false,
				Mode:        "sync",
			},
		}

		data2 := &ServiceData{
			SyncCount: 10,
			Status: ServiceDataStatus{
				IsConnected: true,
				Mode:        "idle",
			},
		}

		diff := data2.Diff(data1) // new.Diff(old) semantics

		// Marshal the diff to JSON to see what would be stored in the database
		jsonBytes, err := json.Marshal(diff)
		if err != nil {
			t.Fatalf("Failed to marshal diff to JSON: %v", err)
		}

		jsonStr := string(jsonBytes)

		// The JSON should not contain any references to gorm.Expr or SQL structures
		forbiddenStrings := []string{
			"gorm.Expr",
			"clause.Expr",
			"\"SQL\":",
			"\"Vars\":",
			"WithoutParentheses",
		}

		for _, forbidden := range forbiddenStrings {
			if strings.Contains(jsonStr, forbidden) {
				t.Errorf("JSON output contains forbidden string '%s': %s", forbidden, jsonStr)
			}
		}

		// Verify the JSON structure is clean
		var parsed map[string]interface{}
		if err := json.Unmarshal(jsonBytes, &parsed); err != nil {
			t.Fatalf("Failed to parse generated JSON: %v", err)
		}

		// Should have clean nested structure
		if statusInterface, exists := parsed["status"]; exists {
			statusMap, ok := statusInterface.(map[string]interface{})
			if !ok {
				t.Errorf("Expected status to be a clean map, got %T", statusInterface)
			} else {
				// Verify nested values are clean (should be new values from data2)
				if statusMap["isConnected"] != true {
					t.Errorf("Expected clean boolean value, got %v", statusMap["isConnected"])
				}
				if statusMap["mode"] != "idle" {
					t.Errorf("Expected clean string value, got %v", statusMap["mode"])
				}
			}
		}
	})
}

func TestEmptyNestedJSONPrevention(t *testing.T) {
	// Test that empty nested JSON objects don't create unnecessary diff entries
	
	t.Run("Empty nested status should not create diff", func(t *testing.T) {
		data1 := &ServiceData{
			MyId:      "test123",
			SyncCount: 5,
			Status:    ServiceDataStatus{}, // Empty
		}

		data2 := &ServiceData{
			MyId:      "test123",
			SyncCount: 5,
			Status:    ServiceDataStatus{}, // Same empty
		}

		diff := data2.Diff(data1) // new.Diff(old) semantics

		// Should have no differences
		if len(diff) != 0 {
			t.Errorf("Expected no differences for identical data, got %d: %v", len(diff), diff)
		}
	})

	t.Run("Partial nested changes should only include changed fields", func(t *testing.T) {
		data1 := &ServiceData{
			MyId:      "test123",
			SyncCount: 5,
			Status: ServiceDataStatus{
				IsConnected: false,
				IsStarting:  false, // This won't change
				Mode:        "sync",
			},
		}

		data2 := &ServiceData{
			MyId:      "test123",
			SyncCount: 5,
			Status: ServiceDataStatus{
				IsConnected: true,   // This changes
				IsStarting:  false,  // This stays the same
				Mode:        "sync", // This stays the same
			},
		}

		diff := data2.Diff(data1) // new.Diff(old) semantics

		// Should only have status diff
		if len(diff) != 1 {
			t.Errorf("Expected 1 difference, got %d: %v", len(diff), diff)
		}

		statusDiff, exists := diff["status"]
		if !exists {
			t.Fatal("Expected status diff to exist")
		}

		statusMap := statusDiff.(map[string]interface{})

		// Should only contain the changed field
		if len(statusMap) != 1 {
			t.Errorf("Expected only 1 field in status diff, got %d: %v", len(statusMap), statusMap)
		}

		if statusMap["isConnected"] != true {
			t.Errorf("Expected isConnected to be true, got %v", statusMap["isConnected"])
		}

		// Verify unchanged fields are not included
		if _, exists := statusMap["isStarting"]; exists {
			t.Error("Expected isStarting to NOT be in diff since it's unchanged")
		}
		if _, exists := statusMap["mode"]; exists {
			t.Error("Expected mode to NOT be in diff since it's unchanged")
		}
	})
}

func TestNestedJSONTypeConsistency(t *testing.T) {
	// Verify that the types in nested JSON diffs are consistent and serializable
	
	data1 := &ServiceData{
		SyncCount: 5,
		Status: ServiceDataStatus{
			IsConnected: false,
			IsStarting:  true,
			Mode:        "sync",
		},
	}

	data2 := &ServiceData{
		SyncCount: 10,
		Status: ServiceDataStatus{
			IsConnected: true,
			IsStarting:  false,
			Mode:        "idle",
		},
	}

	diff := data1.Diff(data2)

	// Verify all values in the diff are JSON-serializable types
	for key, value := range diff {
		if !isJSONSerializable(value) {
			t.Errorf("Diff value for key '%s' is not JSON serializable: %T %v", key, value, value)
		}
	}
}

// Helper function to check if a value is JSON serializable
func isJSONSerializable(v interface{}) bool {
	switch v.(type) {
	case nil, bool, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, string:
		return true
	case map[string]interface{}:
		// Recursively check map values
		m := v.(map[string]interface{})
		for _, val := range m {
			if !isJSONSerializable(val) {
				return false
			}
		}
		return true
	case []interface{}:
		// Recursively check slice values
		s := v.([]interface{})
		for _, val := range s {
			if !isJSONSerializable(val) {
				return false
			}
		}
		return true
	default:
		// Check if it's a struct that can be marshaled to JSON
		_, err := json.Marshal(v)
		return err == nil && !isGormExpr(v)
	}
}

// Helper function to detect gorm.Expr types
func isGormExpr(v interface{}) bool {
	t := reflect.TypeOf(v)
	if t == nil {
		return false
	}
	
	// Check for clause.Expr (which is what gorm.Expr returns)
	return t.String() == "clause.Expr" || strings.Contains(t.String(), "clause.Expr")
}
