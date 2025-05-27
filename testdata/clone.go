package structs

import (
	"reflect"
)

// Clone creates a deep copy of the Address struct
func (original Address) Clone() Address {
	clone := Address{}

	// Clone Street

	// Simple type - direct assignment
	clone.Street = original.Street

	// Clone City

	// Simple type - direct assignment
	clone.City = original.City

	// Clone State

	// Simple type - direct assignment
	clone.State = original.State

	// Clone ZipCode

	// Simple type - direct assignment
	clone.ZipCode = original.ZipCode

	// Clone Country

	// Simple type - direct assignment
	clone.Country = original.Country

	return clone
}

// Clone creates a deep copy of the Contact struct
func (original Contact) Clone() Contact {
	clone := Contact{}

	// Clone Type

	// Simple type - direct assignment
	clone.Type = original.Type

	// Clone Value

	// Simple type - direct assignment
	clone.Value = original.Value

	return clone
}

// Clone creates a deep copy of the Person struct
func (original Person) Clone() Person {
	clone := Person{}

	// Clone Name

	// Simple type - direct assignment
	clone.Name = original.Name

	// Clone Age

	// Simple type - direct assignment
	clone.Age = original.Age

	// Clone Address

	// Struct type - recursive clone
	clone.Address = original.Address.Clone()

	// Clone Contacts

	// Slice - create new slice and clone elements
	if original.Contacts != nil {
		clone.Contacts = make([]Contact, len(original.Contacts))

		for i, item := range original.Contacts {
			clone.Contacts[i] = item.Clone()
		}

	}

	// Clone Manager

	// Pointer to struct - create new instance and clone
	if original.Manager != nil {
		clonedManager := original.Manager.Clone()
		clone.Manager = &clonedManager
	}

	// Clone Metadata

	// Map - create new map and copy key-value pairs
	if original.Metadata != nil {
		clone.Metadata = make(map[string]interface{})
		for k, v := range original.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}

// Clone creates a deep copy of the Company struct
func (original Company) Clone() Company {
	clone := Company{}

	// Clone Name

	// Simple type - direct assignment
	clone.Name = original.Name

	// Clone Address

	// Struct type - recursive clone
	clone.Address = original.Address.Clone()

	// Clone Employees

	// Slice - create new slice and clone elements
	if original.Employees != nil {
		clone.Employees = make([]Person, len(original.Employees))

		for i, item := range original.Employees {
			clone.Employees[i] = item.Clone()
		}

	}

	// Clone CEO

	// Pointer to struct - create new instance and clone
	if original.CEO != nil {
		clonedCEO := original.CEO.Clone()
		clone.CEO = &clonedCEO
	}

	// Clone Founded

	// Simple type - direct assignment
	clone.Founded = original.Founded

	// Clone Active

	// Simple type - direct assignment
	clone.Active = original.Active

	return clone
}

// Clone creates a deep copy of the Project struct
func (original Project) Clone() Project {
	clone := Project{}

	// Clone Name

	// Simple type - direct assignment
	clone.Name = original.Name

	// Clone Description

	// Simple type - direct assignment
	clone.Description = original.Description

	// Clone TeamLead

	// Pointer to struct - create new instance and clone
	if original.TeamLead != nil {
		clonedTeamLead := original.TeamLead.Clone()
		clone.TeamLead = &clonedTeamLead
	}

	// Clone Members

	// Slice - create new slice and clone elements
	if original.Members != nil {
		clone.Members = make([]*Person, len(original.Members))

		for i, item := range original.Members {
			clone.Members[i] = item.Clone()
		}

	}

	// Clone Company

	// Pointer to struct - create new instance and clone
	if original.Company != nil {
		clonedCompany := original.Company.Clone()
		clone.Company = &clonedCompany
	}

	// Clone Budget

	// Simple type - direct assignment
	clone.Budget = original.Budget

	// Clone Tags

	// Slice - create new slice and clone elements
	if original.Tags != nil {
		clone.Tags = make([]string, len(original.Tags))

		copy(clone.Tags, original.Tags)

	}

	// Clone Properties

	// Map - create new map and copy key-value pairs
	if original.Properties != nil {
		clone.Properties = make(map[string]string)
		for k, v := range original.Properties {
			clone.Properties[k] = v
		}
	}

	return clone
}
