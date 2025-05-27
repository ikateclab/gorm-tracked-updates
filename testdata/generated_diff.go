package structs

import (
	"reflect"
)

// DiffAddress compares two Address instances and returns a map of differences
// with only the new values for fields that have changed
func DiffAddress(a, b Address) map[string]interface{} {
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

// DiffContact compares two Contact instances and returns a map of differences
// with only the new values for fields that have changed
func DiffContact(a, b Contact) map[string]interface{} {
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

// DiffPerson compares two Person instances and returns a map of differences
// with only the new values for fields that have changed
func DiffPerson(a, b Person) map[string]interface{} {
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
		nestedDiff := DiffAddress(a.Address, b.Address)
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
			nestedDiff := DiffPerson(*a.Manager, *b.Manager)
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

// DiffCompany compares two Company instances and returns a map of differences
// with only the new values for fields that have changed
func DiffCompany(a, b Company) map[string]interface{} {
	diff := make(map[string]interface{})

	// Compare Name

	// Simple type comparison
	if a.Name != b.Name {
		diff["Name"] = b.Name
	}

	// Compare Address

	// Struct type comparison
	if !reflect.DeepEqual(a.Address, b.Address) {
		nestedDiff := DiffAddress(a.Address, b.Address)
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
			nestedDiff := DiffPerson(*a.CEO, *b.CEO)
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

// DiffProject compares two Project instances and returns a map of differences
// with only the new values for fields that have changed
func DiffProject(a, b Project) map[string]interface{} {
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
			nestedDiff := DiffPerson(*a.TeamLead, *b.TeamLead)
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
			nestedDiff := DiffCompany(*a.Company, *b.Company)
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
