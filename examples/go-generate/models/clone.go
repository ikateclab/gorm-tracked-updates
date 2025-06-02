package models

import (
	"gorm.io/datatypes"
)

// Clone creates a deep copy of the AccountSettings struct
func (original *AccountSettings) Clone() *AccountSettings {
	if original == nil {
		return nil
	}
	// Create new instance - all fields are simple types
	clone := *original
	return &clone
}

// Clone creates a deep copy of the AccountData struct
func (original *AccountData) Clone() *AccountData {
	if original == nil {
		return nil
	}
	// Create new instance - all fields are simple types
	clone := *original
	return &clone
}

// Clone creates a deep copy of the Account struct
func (original *Account) Clone() *Account {
	if original == nil {
		return nil
	}
	// Create new instance and copy all simple fields
	clone := *original

	// Only handle the fields that need deep cloning

	// TODO: Id (uuid.UUID) may need manual deep copy handling

	if original.Settings != nil {
		clone.Settings = original.Settings.Clone()
	}

	if original.Data != nil {
		clone.Data = original.Data.Clone()
	}

	// TODO: DeletedAt (gorm.DeletedAt) may need manual deep copy handling

	if original.Services != nil {
		clone.Services = make([]*Service, len(original.Services))

		for i, item := range original.Services {
			if item != nil {
				clone.Services[i] = item.Clone()
			}
		}

	}

	return &clone
}

// Clone creates a deep copy of the ServerPod struct
func (original *ServerPod) Clone() *ServerPod {
	if original == nil {
		return nil
	}
	// Create new instance and copy all simple fields
	clone := *original

	// Only handle the fields that need deep cloning

	// TODO: Id (uuid.UUID) may need manual deep copy handling

	if original.Settings != nil {
		clone.Settings = make(datatypes.JSON, len(original.Settings))

		copy(clone.Settings, original.Settings)

	}

	// TODO: DeletedAt (gorm.DeletedAt) may need manual deep copy handling

	// TODO: ServerPodTypeId (uuid.UUID) may need manual deep copy handling

	clone.ServerPodType = *(&original.ServerPodType).Clone()

	return &clone
}

// Clone creates a deep copy of the ServiceVersion struct
func (original *ServiceVersion) Clone() *ServiceVersion {
	if original == nil {
		return nil
	}
	// Create new instance - all fields are simple types
	clone := *original
	return &clone
}

// Clone creates a deep copy of the ServerPodType struct
func (original *ServerPodType) Clone() *ServerPodType {
	if original == nil {
		return nil
	}
	// Create new instance and copy all simple fields
	clone := *original

	// Only handle the fields that need deep cloning

	// TODO: Id (uuid.UUID) may need manual deep copy handling

	if original.Version != nil {
		clone.Version = original.Version.Clone()
	}

	// TODO: AccountIdWhitelist (JsonbStringSlice) may need manual deep copy handling

	// TODO: ServiceIdWhitelist (JsonbStringSlice) may need manual deep copy handling

	// TODO: DeletedAt (gorm.DeletedAt) may need manual deep copy handling

	return &clone
}

// Clone creates a deep copy of the ServiceDataStatus struct
func (original *ServiceDataStatus) Clone() *ServiceDataStatus {
	if original == nil {
		return nil
	}
	// Create new instance - all fields are simple types
	clone := *original
	return &clone
}

// Clone creates a deep copy of the ServiceData struct
func (original *ServiceData) Clone() *ServiceData {
	if original == nil {
		return nil
	}
	// Create new instance and copy all simple fields
	clone := *original

	// Only handle the fields that need deep cloning

	clone.Status = *(&original.Status).Clone()

	return &clone
}

// Clone creates a deep copy of the ServiceSettings struct
func (original *ServiceSettings) Clone() *ServiceSettings {
	if original == nil {
		return nil
	}
	// Create new instance - all fields are simple types
	clone := *original
	return &clone
}

// Clone creates a deep copy of the Service struct
func (original *Service) Clone() *Service {
	if original == nil {
		return nil
	}
	// Create new instance and copy all simple fields
	clone := *original

	// Only handle the fields that need deep cloning

	// TODO: Id (uuid.UUID) may need manual deep copy handling

	if original.Data != nil {
		clone.Data = original.Data.Clone()
	}

	if original.Settings != nil {
		clone.Settings = original.Settings.Clone()
	}

	// TODO: DeletedAt (gorm.DeletedAt) may need manual deep copy handling

	// TODO: AccountId (uuid.UUID) may need manual deep copy handling

	// TODO: ServerPodId (*uuid.UUID) may need manual deep copy handling

	if original.Account != nil {
		clone.Account = original.Account.Clone()
	}

	if original.ServerPod != nil {
		clone.ServerPod = original.ServerPod.Clone()
	}

	return &clone
}
