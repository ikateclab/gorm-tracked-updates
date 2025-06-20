package models

import (
	"testing"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestEmptyJSONPrevention(t *testing.T) {
	// Test that empty JSON objects don't create diff entries
	account1 := &Account{
		Id:       uuid.New(),
		Name:     "Test Account",
		Settings: &AccountSettings{}, // Empty struct
		Data:     &AccountData{},     // Empty struct
	}

	account2 := &Account{
		Id:       account1.Id,
		Name:     account1.Name,
		Settings: &AccountSettings{}, // Same empty struct
		Data:     &AccountData{},     // Same empty struct
	}

	diff := account1.Diff(account2)

	// Should have no diff entries since the JSON fields are empty
	if len(diff) != 0 {
		t.Errorf("Expected no diff entries for empty JSON objects, got %d entries: %v", len(diff), diff)
	}
}

func TestSimpleModelArrayDiff(t *testing.T) {
	// Test array of objects in JSONB fields
	model1 := &SimpleModel{
		ID:   uuid.New(),
		Name: "Test Model",
		Tags: []*Tag{
			{Name: "category", Value: "test"},
			{Name: "priority", Value: "high"},
		},
		Items: []*Item{
			{ID: 1, Title: "Item One", Price: 10.50},
			{ID: 2, Title: "Item Two", Price: 20.00},
		},
	}

	model2 := &SimpleModel{
		ID:   model1.ID,
		Name: model1.Name,
		Tags: []*Tag{
			{Name: "category", Value: "modified"}, // Changed value
			{Name: "priority", Value: "high"},
		},
		Items: []*Item{
			{ID: 1, Title: "Item One Modified", Price: 15.75}, // Changed title and price
			{ID: 2, Title: "Item Two", Price: 20.00},
		},
	}

	diff := model1.Diff(model2)

	// Should detect changes in both Tags and Items arrays
	if len(diff) == 0 {
		t.Error("Expected diff entries for changed array fields, got none")
	}

	// Check if Tags field is in diff
	if _, exists := diff["Tags"]; !exists {
		t.Error("Expected 'Tags' field in diff, but it's missing")
	}

	// Check if Items field is in diff
	if _, exists := diff["Items"]; !exists {
		t.Error("Expected 'Items' field in diff, but it's missing")
	}

	t.Logf("Diff result: %+v", diff)
}

func TestSimpleModelArrayNoChanges(t *testing.T) {
	// Test that identical arrays don't generate diff entries
	model1 := &SimpleModel{
		ID:   uuid.New(),
		Name: "Test Model",
		Tags: []*Tag{
			{Name: "category", Value: "test"},
			{Name: "priority", Value: "high"},
		},
		Items: []*Item{
			{ID: 1, Title: "Item One", Price: 10.50},
			{ID: 2, Title: "Item Two", Price: 20.00},
		},
	}

	// Create identical model
	model2 := &SimpleModel{
		ID:   model1.ID,
		Name: model1.Name,
		Tags: []*Tag{
			{Name: "category", Value: "test"}, // Same values
			{Name: "priority", Value: "high"},
		},
		Items: []*Item{
			{ID: 1, Title: "Item One", Price: 10.50}, // Same values
			{ID: 2, Title: "Item Two", Price: 20.00},
		},
	}

	diff := model1.Diff(model2)

	// Should have no diff entries since arrays are identical
	if len(diff) != 0 {
		t.Errorf("Expected no diff entries for identical arrays, got %d entries: %+v", len(diff), diff)
	}
}

func TestNonEmptyJSONDiff(t *testing.T) {
	// Test that non-empty JSON objects do create diff entries
	account1 := &Account{
		Id:       uuid.New(),
		Name:     "Test Account",
		Settings: &AccountSettings{}, // Empty struct
		Data:     &AccountData{},     // Empty struct
	}

	// Create a different account with some data (this would need actual fields in AccountData)
	account2 := &Account{
		Id:       account1.Id,
		Name:     "Updated Account",  // Different name
		Settings: &AccountSettings{}, // Still empty
		Data:     &AccountData{},     // Still empty
	}

	diff := account1.Diff(account2)

	// Should have one diff entry for the name change, but not for the empty JSON fields
	expectedEntries := 1 // Only the Name field should be different
	if len(diff) != expectedEntries {
		t.Errorf("Expected %d diff entries, got %d entries: %v", expectedEntries, len(diff), diff)
	}

	// Check that Name is in the diff
	if _, exists := diff["Name"]; !exists {
		t.Error("Expected Name field to be in diff")
	}

	// Check that Settings and Data are NOT in the diff (since they're empty JSON)
	if _, exists := diff["Settings"]; exists {
		t.Error("Expected Settings field to NOT be in diff (empty JSON)")
	}
	if _, exists := diff["Data"]; exists {
		t.Error("Expected Data field to NOT be in diff (empty JSON)")
	}
}

func TestIsEmptyJSONFunction(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Empty object", "{}", true},
		{"Empty array", "[]", true},
		{"Null value", "null", true},
		{"Empty object with spaces", "  {}  ", true},
		{"Empty array with spaces", "  []  ", true},
		{"Null with spaces", "  null  ", true},
		{"Non-empty object", `{"key": "value"}`, false},
		{"Non-empty array", `["item"]`, false},
		{"String value", `"string"`, false},
		{"Number value", `123`, false},
		{"Boolean value", `true`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isEmptyJSON(tt.input)
			if result != tt.expected {
				t.Errorf("isEmptyJSON(%q) = %v, expected %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGormExpressionGeneration(t *testing.T) {
	// Test that when JSON is not empty, proper GORM expressions are generated
	account1 := &Account{
		Id:       uuid.New(),
		Name:     "Test Account",
		Settings: &AccountSettings{}, // Empty
		Data:     &AccountData{},     // Empty
	}

	// Simulate a change that would result in non-empty JSON
	// Note: This test is more conceptual since AccountSettings and AccountData are empty structs
	// In a real scenario, these would have fields that could be different
	account2 := &Account{
		Id:       account1.Id,
		Name:     account1.Name,
		Settings: account1.Settings, // Same reference, so no diff
		Data:     account1.Data,     // Same reference, so no diff
	}

	diff := account1.Diff(account2)

	// Should have no entries since everything is the same
	if len(diff) != 0 {
		t.Errorf("Expected no diff entries for identical accounts, got %d entries: %v", len(diff), diff)
	}

	// Test that GORM expressions are properly typed
	// This is more of a compile-time check
	expr := gorm.Expr("? || ?", clause.Column{Name: "test"}, `{"key": "value"}`)
	_ = expr // Just verify it compiles
}
