package structs

// Clone creates a deep copy of the Address struct
func (original *Address) Clone() *Address {
	if original == nil {
		return nil
	}
	// Create new instance - all fields are simple types
	clone := *original
	return &clone
}

// Clone creates a deep copy of the Contact struct
func (original *Contact) Clone() *Contact {
	if original == nil {
		return nil
	}
	// Create new instance - all fields are simple types
	clone := *original
	return &clone
}

// Clone creates a deep copy of the Person struct
func (original *Person) Clone() *Person {
	if original == nil {
		return nil
	}
	// Create new instance and copy all simple fields
	clone := *original

	// Only handle the fields that need deep cloning

	clone.Address = *(&original.Address).Clone()

	if original.Contacts != nil {
		clone.Contacts = make([]Contact, len(original.Contacts))

		for i, item := range original.Contacts {
			clone.Contacts[i] = item.Clone()
		}

	}

	if original.Manager != nil {
		clone.Manager = original.Manager.Clone()
	}

	if original.Metadata != nil {
		clone.Metadata = make(map[string]interface{})
		for k, v := range original.Metadata {

			clone.Metadata[k] = v

		}
	}

	return &clone
}

// Clone creates a deep copy of the Company struct
func (original *Company) Clone() *Company {
	if original == nil {
		return nil
	}
	// Create new instance and copy all simple fields
	clone := *original

	// Only handle the fields that need deep cloning

	clone.Address = *(&original.Address).Clone()

	if original.Employees != nil {
		clone.Employees = make([]Person, len(original.Employees))

		for i, item := range original.Employees {
			clone.Employees[i] = item.Clone()
		}

	}

	if original.CEO != nil {
		clone.CEO = original.CEO.Clone()
	}

	return &clone
}

// Clone creates a deep copy of the Project struct
func (original *Project) Clone() *Project {
	if original == nil {
		return nil
	}
	// Create new instance and copy all simple fields
	clone := *original

	// Only handle the fields that need deep cloning

	if original.TeamLead != nil {
		clone.TeamLead = original.TeamLead.Clone()
	}

	if original.Members != nil {
		clone.Members = make([]*Person, len(original.Members))

		for i, item := range original.Members {
			if item != nil {
				clone.Members[i] = item.Clone()
			}
		}

	}

	if original.Company != nil {
		clone.Company = original.Company.Clone()
	}

	if original.Tags != nil {
		clone.Tags = make([]string, len(original.Tags))

		copy(clone.Tags, original.Tags)

	}

	if original.Properties != nil {
		clone.Properties = make(map[string]string)
		for k, v := range original.Properties {

			clone.Properties[k] = v

		}
	}

	return &clone
}
