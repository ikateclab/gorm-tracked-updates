package models

// Clone creates a deep copy of the Address struct
func (original Address) Clone() Address {
	clone := Address{}

	// Clone ID

	// Simple type - direct assignment
	clone.ID = original.ID

	// Clone UserID

	// Simple type - direct assignment
	clone.UserID = original.UserID

	// Clone Type

	// Simple type - direct assignment
	clone.Type = original.Type

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

	// Clone Primary

	// Simple type - direct assignment
	clone.Primary = original.Primary

	return clone
}

// Clone creates a deep copy of the Order struct
func (original Order) Clone() Order {
	clone := Order{}

	// Clone ID

	// Simple type - direct assignment
	clone.ID = original.ID

	// Clone UserID

	// Simple type - direct assignment
	clone.UserID = original.UserID

	// Clone User

	// Pointer to struct - create new instance and clone
	if original.User != nil {
		clonedUser := original.User.Clone()
		clone.User = &clonedUser
	}

	// Clone Items

	// Slice - create new slice and clone elements
	if original.Items != nil {
		clone.Items = make([]OrderItem, len(original.Items))

		for i, item := range original.Items {
			clone.Items[i] = item.Clone()
		}

	}

	// Clone Total

	// Simple type - direct assignment
	clone.Total = original.Total

	// Clone Status

	// Simple type - direct assignment
	clone.Status = original.Status

	// Clone ShippingAddress

	// Struct type - recursive clone
	clone.ShippingAddress = original.ShippingAddress.Clone()

	// Clone BillingAddress

	// Struct type - recursive clone
	clone.BillingAddress = original.BillingAddress.Clone()

	// Clone CreatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.CreatedAt = original.CreatedAt

	// Clone UpdatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.UpdatedAt = original.UpdatedAt

	return clone
}

// Clone creates a deep copy of the OrderItem struct
func (original OrderItem) Clone() OrderItem {
	clone := OrderItem{}

	// Clone ID

	// Simple type - direct assignment
	clone.ID = original.ID

	// Clone OrderID

	// Simple type - direct assignment
	clone.OrderID = original.OrderID

	// Clone ProductID

	// Simple type - direct assignment
	clone.ProductID = original.ProductID

	// Clone Quantity

	// Simple type - direct assignment
	clone.Quantity = original.Quantity

	// Clone Price

	// Simple type - direct assignment
	clone.Price = original.Price

	// Clone Total

	// Simple type - direct assignment
	clone.Total = original.Total

	return clone
}

// Clone creates a deep copy of the User struct
func (original User) Clone() User {
	clone := User{}

	// Clone ID

	// Simple type - direct assignment
	clone.ID = original.ID

	// Clone Name

	// Simple type - direct assignment
	clone.Name = original.Name

	// Clone Email

	// Simple type - direct assignment
	clone.Email = original.Email

	// Clone Age

	// Simple type - direct assignment
	clone.Age = original.Age

	// Clone Profile

	// Struct type - recursive clone
	clone.Profile = original.Profile.Clone()

	// Clone Addresses

	// Slice - create new slice and clone elements
	if original.Addresses != nil {
		clone.Addresses = make([]Address, len(original.Addresses))

		for i, item := range original.Addresses {
			clone.Addresses[i] = item.Clone()
		}

	}

	// Clone CreatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.CreatedAt = original.CreatedAt

	// Clone UpdatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.UpdatedAt = original.UpdatedAt

	return clone
}

// Clone creates a deep copy of the Profile struct
func (original Profile) Clone() Profile {
	clone := Profile{}

	// Clone Bio

	// Simple type - direct assignment
	clone.Bio = original.Bio

	// Clone Avatar

	// Simple type - direct assignment
	clone.Avatar = original.Avatar

	// Clone Verified

	// Simple type - direct assignment
	clone.Verified = original.Verified

	// Clone Settings

	// Map - create new map and copy key-value pairs
	if original.Settings != nil {
		clone.Settings = make(map[string]interface{})
		for k, v := range original.Settings {
			clone.Settings[k] = v
		}
	}

	// Clone Metadata

	// Map - create new map and copy key-value pairs
	if original.Metadata != nil {
		clone.Metadata = make(map[string]string)
		for k, v := range original.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}
