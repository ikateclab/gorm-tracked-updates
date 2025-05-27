package models

import (
	"encoding/json"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
)

// Diff compares this Address instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a Address) Diff(b Address) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare ID

	// Simple type comparison
	if a.ID != b.ID {
		diff["ID"] = b.ID
	}

	// Compare UserID

	// Simple type comparison
	if a.UserID != b.UserID {
		diff["UserID"] = b.UserID
	}

	// Compare Type

	// Simple type comparison
	if a.Type != b.Type {
		diff["Type"] = b.Type
	}

	// Compare Street

	// Simple type comparison
	if a.Street != b.Street {
		diff["Street"] = b.Street
	}

	// Compare City

	// Simple type comparison
	if a.City != b.City {
		diff["City"] = b.City
	}

	// Compare State

	// Simple type comparison
	if a.State != b.State {
		diff["State"] = b.State
	}

	// Compare ZipCode

	// Simple type comparison
	if a.ZipCode != b.ZipCode {
		diff["ZipCode"] = b.ZipCode
	}

	// Compare Country

	// Simple type comparison
	if a.Country != b.Country {
		diff["Country"] = b.Country
	}

	// Compare Primary

	// Simple type comparison
	if a.Primary != b.Primary {
		diff["Primary"] = b.Primary
	}

	return diff
}

// Diff compares this Order instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a Order) Diff(b Order) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare ID

	// Simple type comparison
	if a.ID != b.ID {
		diff["ID"] = b.ID
	}

	// Compare UserID

	// Simple type comparison
	if a.UserID != b.UserID {
		diff["UserID"] = b.UserID
	}

	// Compare User

	// Pointer to struct comparison
	if !reflect.DeepEqual(a.User, b.User) {
		if a.User == nil || b.User == nil {
			diff["User"] = b.User
		} else {
			nestedDiff := (*a.User).Diff(*b.User)
			if len(nestedDiff) > 0 {
				diff["User"] = nestedDiff
			}
		}
	}

	// Compare Items

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.Items, b.Items) {
		diff["Items"] = b.Items
	}

	// Compare Total

	// Simple type comparison
	if a.Total != b.Total {
		diff["Total"] = b.Total
	}

	// Compare Status

	// Simple type comparison
	if a.Status != b.Status {
		diff["Status"] = b.Status
	}

	// Compare ShippingAddress

	// Struct type comparison
	if !reflect.DeepEqual(a.ShippingAddress, b.ShippingAddress) {
		nestedDiff := a.ShippingAddress.Diff(b.ShippingAddress)
		if len(nestedDiff) > 0 {
			diff["ShippingAddress"] = nestedDiff
		}
	}

	// Compare BillingAddress

	// Struct type comparison
	if !reflect.DeepEqual(a.BillingAddress, b.BillingAddress) {
		nestedDiff := a.BillingAddress.Diff(b.BillingAddress)
		if len(nestedDiff) > 0 {
			diff["BillingAddress"] = nestedDiff
		}
	}

	// Compare CreatedAt

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.CreatedAt, b.CreatedAt) {
		diff["CreatedAt"] = b.CreatedAt
	}

	// Compare UpdatedAt

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.UpdatedAt, b.UpdatedAt) {
		diff["UpdatedAt"] = b.UpdatedAt
	}

	return diff
}

// Diff compares this OrderItem instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a OrderItem) Diff(b OrderItem) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare ID

	// Simple type comparison
	if a.ID != b.ID {
		diff["ID"] = b.ID
	}

	// Compare OrderID

	// Simple type comparison
	if a.OrderID != b.OrderID {
		diff["OrderID"] = b.OrderID
	}

	// Compare ProductID

	// Simple type comparison
	if a.ProductID != b.ProductID {
		diff["ProductID"] = b.ProductID
	}

	// Compare Quantity

	// Simple type comparison
	if a.Quantity != b.Quantity {
		diff["Quantity"] = b.Quantity
	}

	// Compare Price

	// Simple type comparison
	if a.Price != b.Price {
		diff["Price"] = b.Price
	}

	// Compare Total

	// Simple type comparison
	if a.Total != b.Total {
		diff["Total"] = b.Total
	}

	return diff
}

// Diff compares this User instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a User) Diff(b User) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare ID

	// Simple type comparison
	if a.ID != b.ID {
		diff["ID"] = b.ID
	}

	// Compare Name

	// Simple type comparison
	if a.Name != b.Name {
		diff["Name"] = b.Name
	}

	// Compare Email

	// Simple type comparison
	if a.Email != b.Email {
		diff["Email"] = b.Email
	}

	// Compare Age

	// Simple type comparison
	if a.Age != b.Age {
		diff["Age"] = b.Age
	}

	// Compare Profile

	// Struct type comparison
	if !reflect.DeepEqual(a.Profile, b.Profile) {
		nestedDiff := a.Profile.Diff(b.Profile)
		if len(nestedDiff) > 0 {
			diff["Profile"] = nestedDiff
		}
	}

	// Compare Addresses

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.Addresses, b.Addresses) {
		diff["Addresses"] = b.Addresses
	}

	// Compare CreatedAt

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.CreatedAt, b.CreatedAt) {
		diff["CreatedAt"] = b.CreatedAt
	}

	// Compare UpdatedAt

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.UpdatedAt, b.UpdatedAt) {
		diff["UpdatedAt"] = b.UpdatedAt
	}

	return diff
}

// Diff compares this Profile instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a Profile) Diff(b Profile) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare Bio

	// Simple type comparison
	if a.Bio != b.Bio {
		diff["Bio"] = b.Bio
	}

	// Compare Avatar

	// Simple type comparison
	if a.Avatar != b.Avatar {
		diff["Avatar"] = b.Avatar
	}

	// Compare Verified

	// Simple type comparison
	if a.Verified != b.Verified {
		diff["Verified"] = b.Verified
	}

	// Compare Settings

	// JSON field comparison - use GORM JSON merge expression
	if !reflect.DeepEqual(a.Settings, b.Settings) {
		jsonValue, err := json.Marshal(b.Settings)
		if err == nil {
			diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "Settings"}, string(jsonValue))
		} else {
			// Fallback to regular assignment if JSON marshaling fails
			diff["Settings"] = b.Settings
		}
	}

	// Compare Metadata

	// JSON field comparison - use GORM JSON merge expression
	if !reflect.DeepEqual(a.Metadata, b.Metadata) {
		jsonValue, err := json.Marshal(b.Metadata)
		if err == nil {
			diff["Metadata"] = gorm.Expr("? || ?", clause.Column{Name: "Metadata"}, string(jsonValue))
		} else {
			// Fallback to regular assignment if JSON marshaling fails
			diff["Metadata"] = b.Metadata
		}
	}

	return diff
}
