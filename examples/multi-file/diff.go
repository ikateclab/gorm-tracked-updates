package multifile

import (
	"reflect"
)

// Diff compares this Address instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *Address) Diff(old *Address) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Street

	// Simple type comparison
	if new.Street != old.Street {
		diff["Street"] = new.Street
	}

	// Compare City

	// Simple type comparison
	if new.City != old.City {
		diff["City"] = new.City
	}

	// Compare State

	// Simple type comparison
	if new.State != old.State {
		diff["State"] = new.State
	}

	// Compare ZipCode

	// Simple type comparison
	if new.ZipCode != old.ZipCode {
		diff["ZipCode"] = new.ZipCode
	}

	// Compare Country

	// Simple type comparison
	if new.Country != old.Country {
		diff["Country"] = new.Country
	}

	return diff
}

// Diff compares this Company instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *Company) Diff(old *Company) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Name

	// Simple type comparison
	if new.Name != old.Name {
		diff["Name"] = new.Name
	}

	// Compare Address

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.Address, old.Address) {
		diff["Address"] = new.Address
	}

	// Compare Employees

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.Employees, old.Employees) {
		diff["Employees"] = new.Employees
	}

	// Compare CEO

	// Comparable type comparison
	if new.CEO != old.CEO {
		diff["CEO"] = new.CEO
	}

	// Compare Founded

	// Simple type comparison
	if new.Founded != old.Founded {
		diff["Founded"] = new.Founded
	}

	// Compare Active

	// Simple type comparison
	if new.Active != old.Active {
		diff["Active"] = new.Active
	}

	return diff
}

// Diff compares this Project instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *Project) Diff(old *Project) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Name

	// Simple type comparison
	if new.Name != old.Name {
		diff["Name"] = new.Name
	}

	// Compare Description

	// Simple type comparison
	if new.Description != old.Description {
		diff["Description"] = new.Description
	}

	// Compare TeamLead

	// Comparable type comparison
	if new.TeamLead != old.TeamLead {
		diff["TeamLead"] = new.TeamLead
	}

	// Compare Members

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.Members, old.Members) {
		diff["Members"] = new.Members
	}

	// Compare Company

	// Comparable type comparison
	if new.Company != old.Company {
		diff["Company"] = new.Company
	}

	// Compare Budget

	// Simple type comparison
	if new.Budget != old.Budget {
		diff["Budget"] = new.Budget
	}

	// Compare Tags

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.Tags, old.Tags) {
		diff["Tags"] = new.Tags
	}

	// Compare Properties

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.Properties, old.Properties) {
		diff["Properties"] = new.Properties
	}

	return diff
}

// Diff compares this Contact instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *Contact) Diff(old *Contact) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Type

	// Simple type comparison
	if new.Type != old.Type {
		diff["Type"] = new.Type
	}

	// Compare Value

	// Simple type comparison
	if new.Value != old.Value {
		diff["Value"] = new.Value
	}

	return diff
}

// Diff compares this Person instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *Person) Diff(old *Person) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Name

	// Simple type comparison
	if new.Name != old.Name {
		diff["Name"] = new.Name
	}

	// Compare Age

	// Simple type comparison
	if new.Age != old.Age {
		diff["Age"] = new.Age
	}

	// Compare Address

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.Address, old.Address) {
		diff["Address"] = new.Address
	}

	// Compare Contacts

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.Contacts, old.Contacts) {
		diff["Contacts"] = new.Contacts
	}

	// Compare Manager

	// Comparable type comparison
	if new.Manager != old.Manager {
		diff["Manager"] = new.Manager
	}

	// Compare Metadata

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.Metadata, old.Metadata) {
		diff["Metadata"] = new.Metadata
	}

	return diff
}
