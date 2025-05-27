package structs

import (
	"reflect"
)

// Diff compares this Address instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a Address) Diff(b Address) map[string]interface{} {
	diff := make(map[string]interface{})

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

	return diff
}

// Diff compares this Contact instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a Contact) Diff(b Contact) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare Type

	// Simple type comparison
	if a.Type != b.Type {
		diff["Type"] = b.Type
	}

	// Compare Value

	// Simple type comparison
	if a.Value != b.Value {
		diff["Value"] = b.Value
	}

	return diff
}

// Diff compares this Person instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a Person) Diff(b Person) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare Name

	// Simple type comparison
	if a.Name != b.Name {
		diff["Name"] = b.Name
	}

	// Compare Age

	// Simple type comparison
	if a.Age != b.Age {
		diff["Age"] = b.Age
	}

	// Compare Address

	// Struct type comparison
	if !reflect.DeepEqual(a.Address, b.Address) {
		nestedDiff := a.Address.Diff(b.Address)
		if len(nestedDiff) > 0 {
			diff["Address"] = nestedDiff
		}
	}

	// Compare Contacts

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.Contacts, b.Contacts) {
		diff["Contacts"] = b.Contacts
	}

	// Compare Manager

	// Pointer to struct comparison
	if !reflect.DeepEqual(a.Manager, b.Manager) {
		if a.Manager == nil || b.Manager == nil {
			diff["Manager"] = b.Manager
		} else {
			nestedDiff := (*a.Manager).Diff(*b.Manager)
			if len(nestedDiff) > 0 {
				diff["Manager"] = nestedDiff
			}
		}
	}

	// Compare Metadata

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.Metadata, b.Metadata) {
		diff["Metadata"] = b.Metadata
	}

	return diff
}

// Diff compares this Company instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a Company) Diff(b Company) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare Name

	// Simple type comparison
	if a.Name != b.Name {
		diff["Name"] = b.Name
	}

	// Compare Address

	// Struct type comparison
	if !reflect.DeepEqual(a.Address, b.Address) {
		nestedDiff := a.Address.Diff(b.Address)
		if len(nestedDiff) > 0 {
			diff["Address"] = nestedDiff
		}
	}

	// Compare Employees

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.Employees, b.Employees) {
		diff["Employees"] = b.Employees
	}

	// Compare CEO

	// Pointer to struct comparison
	if !reflect.DeepEqual(a.CEO, b.CEO) {
		if a.CEO == nil || b.CEO == nil {
			diff["CEO"] = b.CEO
		} else {
			nestedDiff := (*a.CEO).Diff(*b.CEO)
			if len(nestedDiff) > 0 {
				diff["CEO"] = nestedDiff
			}
		}
	}

	// Compare Founded

	// Simple type comparison
	if a.Founded != b.Founded {
		diff["Founded"] = b.Founded
	}

	// Compare Active

	// Simple type comparison
	if a.Active != b.Active {
		diff["Active"] = b.Active
	}

	return diff
}

// Diff compares this Project instance with another and returns a map of differences
// with only the new values for fields that have changed
func (a Project) Diff(b Project) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare Name

	// Simple type comparison
	if a.Name != b.Name {
		diff["Name"] = b.Name
	}

	// Compare Description

	// Simple type comparison
	if a.Description != b.Description {
		diff["Description"] = b.Description
	}

	// Compare TeamLead

	// Pointer to struct comparison
	if !reflect.DeepEqual(a.TeamLead, b.TeamLead) {
		if a.TeamLead == nil || b.TeamLead == nil {
			diff["TeamLead"] = b.TeamLead
		} else {
			nestedDiff := (*a.TeamLead).Diff(*b.TeamLead)
			if len(nestedDiff) > 0 {
				diff["TeamLead"] = nestedDiff
			}
		}
	}

	// Compare Members

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.Members, b.Members) {
		diff["Members"] = b.Members
	}

	// Compare Company

	// Pointer to struct comparison
	if !reflect.DeepEqual(a.Company, b.Company) {
		if a.Company == nil || b.Company == nil {
			diff["Company"] = b.Company
		} else {
			nestedDiff := (*a.Company).Diff(*b.Company)
			if len(nestedDiff) > 0 {
				diff["Company"] = nestedDiff
			}
		}
	}

	// Compare Budget

	// Simple type comparison
	if a.Budget != b.Budget {
		diff["Budget"] = b.Budget
	}

	// Compare Tags

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.Tags, b.Tags) {
		diff["Tags"] = b.Tags
	}

	// Compare Properties

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.Properties, b.Properties) {
		diff["Properties"] = b.Properties
	}

	return diff
}
