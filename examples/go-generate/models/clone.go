package models

// Clone creates a deep copy of the AccountSettings struct
func (original AccountSettings) Clone() AccountSettings {
	clone := AccountSettings{}

	return clone
}

// Clone creates a deep copy of the AccountData struct
func (original AccountData) Clone() AccountData {
	clone := AccountData{}

	return clone
}

// Clone creates a deep copy of the Account struct
func (original Account) Clone() Account {
	clone := Account{}

	// Clone Id

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.Id = original.Id

	// Clone Name

	// Simple type - direct assignment
	clone.Name = original.Name

	// Clone Settings

	// Pointer to struct - create new instance and clone
	if original.Settings != nil {
		clonedSettings := original.Settings.Clone()
		clone.Settings = &clonedSettings
	}

	// Clone Data

	// Pointer to struct - create new instance and clone
	if original.Data != nil {
		clonedData := original.Data.Clone()
		clone.Data = &clonedData
	}

	// Clone IsActive

	// Simple type - direct assignment
	clone.IsActive = original.IsActive

	// Clone CorrelationId

	// Simple type - direct assignment
	clone.CorrelationId = original.CorrelationId

	// Clone WebhookUrl

	// Simple type - direct assignment
	clone.WebhookUrl = original.WebhookUrl

	// Clone CreatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.CreatedAt = original.CreatedAt

	// Clone UpdatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.UpdatedAt = original.UpdatedAt

	// Clone DeletedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.DeletedAt = original.DeletedAt

	// Clone Services

	// Slice - create new slice and clone elements
	if original.Services != nil {
		clone.Services = make([]*Service, len(original.Services))

		for i, item := range original.Services {
			if item != nil {
				clonedItem := item.Clone()
				clone.Services[i] = &clonedItem
			}
		}

	}

	return clone
}

// Clone creates a deep copy of the ServerPod struct
func (original ServerPod) Clone() ServerPod {
	clone := ServerPod{}

	// Clone Id

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.Id = original.Id

	// Clone Name

	// Simple type - direct assignment
	clone.Name = original.Name

	// Clone Address

	// Simple type - direct assignment
	clone.Address = original.Address

	// Clone Version

	// Simple type - direct assignment
	clone.Version = original.Version

	// Clone Settings

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.Settings = original.Settings

	// Clone LastPingAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.LastPingAt = original.LastPingAt

	// Clone CreatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.CreatedAt = original.CreatedAt

	// Clone UpdatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.UpdatedAt = original.UpdatedAt

	// Clone DeletedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.DeletedAt = original.DeletedAt

	// Clone ServerPodTypeId

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.ServerPodTypeId = original.ServerPodTypeId

	// Clone ServerPodType

	// Struct type - recursive clone
	clone.ServerPodType = original.ServerPodType.Clone()

	return clone
}

// Clone creates a deep copy of the ServiceVersion struct
func (original ServiceVersion) Clone() ServiceVersion {
	clone := ServiceVersion{}

	// Clone WppConnectVersion

	// Simple type - direct assignment
	clone.WppConnectVersion = original.WppConnectVersion

	// Clone WaVersion

	// Simple type - direct assignment
	clone.WaVersion = original.WaVersion

	return clone
}

// Clone creates a deep copy of the ServerPodType struct
func (original ServerPodType) Clone() ServerPodType {
	clone := ServerPodType{}

	// Clone Id

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.Id = original.Id

	// Clone Name

	// Simple type - direct assignment
	clone.Name = original.Name

	// Clone Version

	// Pointer to struct - create new instance and clone
	if original.Version != nil {
		clonedVersion := original.Version.Clone()
		clone.Version = &clonedVersion
	}

	// Clone AutoScalable

	// Simple type - direct assignment
	clone.AutoScalable = original.AutoScalable

	// Clone Cloud

	// Simple type - direct assignment
	clone.Cloud = original.Cloud

	// Clone ServerSize

	// Simple type - direct assignment
	clone.ServerSize = original.ServerSize

	// Clone MaxPerPod

	// Simple type - direct assignment
	clone.MaxPerPod = original.MaxPerPod

	// Clone Min

	// Simple type - direct assignment
	clone.Min = original.Min

	// Clone DesiredAvailable

	// Simple type - direct assignment
	clone.DesiredAvailable = original.DesiredAvailable

	// Clone StartPriority

	// Simple type - direct assignment
	clone.StartPriority = original.StartPriority

	// Clone AccountIdWhitelist

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.AccountIdWhitelist = original.AccountIdWhitelist

	// Clone ServiceIdWhitelist

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.ServiceIdWhitelist = original.ServiceIdWhitelist

	// Clone CreatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.CreatedAt = original.CreatedAt

	// Clone UpdatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.UpdatedAt = original.UpdatedAt

	// Clone DeletedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.DeletedAt = original.DeletedAt

	return clone
}

// Clone creates a deep copy of the ServiceDataStatus struct
func (original ServiceDataStatus) Clone() ServiceDataStatus {
	clone := ServiceDataStatus{}

	// Clone IsSyncing

	// Simple type - direct assignment
	clone.IsSyncing = original.IsSyncing

	// Clone IsConnected

	// Simple type - direct assignment
	clone.IsConnected = original.IsConnected

	// Clone IsStarting

	// Simple type - direct assignment
	clone.IsStarting = original.IsStarting

	// Clone IsStarted

	// Simple type - direct assignment
	clone.IsStarted = original.IsStarted

	// Clone IsConflicted

	// Simple type - direct assignment
	clone.IsConflicted = original.IsConflicted

	// Clone IsLoading

	// Simple type - direct assignment
	clone.IsLoading = original.IsLoading

	// Clone IsOnChatPage

	// Simple type - direct assignment
	clone.IsOnChatPage = original.IsOnChatPage

	// Clone EnteredQrCodePageAt

	// Simple type - direct assignment
	clone.EnteredQrCodePageAt = original.EnteredQrCodePageAt

	// Clone DisconnectedAt

	// Simple type - direct assignment
	clone.DisconnectedAt = original.DisconnectedAt

	// Clone IsOnQrPage

	// Simple type - direct assignment
	clone.IsOnQrPage = original.IsOnQrPage

	// Clone IsQrCodeExpired

	// Simple type - direct assignment
	clone.IsQrCodeExpired = original.IsQrCodeExpired

	// Clone IsWebConnected

	// Simple type - direct assignment
	clone.IsWebConnected = original.IsWebConnected

	// Clone IsWebSyncing

	// Simple type - direct assignment
	clone.IsWebSyncing = original.IsWebSyncing

	// Clone Mode

	// Simple type - direct assignment
	clone.Mode = original.Mode

	// Clone MyId

	// Simple type - direct assignment
	clone.MyId = original.MyId

	// Clone MyName

	// Simple type - direct assignment
	clone.MyName = original.MyName

	// Clone MyNumber

	// Simple type - direct assignment
	clone.MyNumber = original.MyNumber

	// Clone QrCodeExpiresAt

	// Simple type - direct assignment
	clone.QrCodeExpiresAt = original.QrCodeExpiresAt

	// Clone QrCodeUrl

	// Simple type - direct assignment
	clone.QrCodeUrl = original.QrCodeUrl

	// Clone State

	// Simple type - direct assignment
	clone.State = original.State

	// Clone WaVersion

	// Simple type - direct assignment
	clone.WaVersion = original.WaVersion

	return clone
}

// Clone creates a deep copy of the ServiceData struct
func (original ServiceData) Clone() ServiceData {
	clone := ServiceData{}

	// Clone MyId

	// Simple type - direct assignment
	clone.MyId = original.MyId

	// Clone LastSyncAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.LastSyncAt = original.LastSyncAt

	// Clone LastMessageTimestamp

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.LastMessageTimestamp = original.LastMessageTimestamp

	// Clone SyncCount

	// Simple type - direct assignment
	clone.SyncCount = original.SyncCount

	// Clone SyncFlowDone

	// Simple type - direct assignment
	clone.SyncFlowDone = original.SyncFlowDone

	// Clone Status

	// Struct type - recursive clone
	clone.Status = original.Status.Clone()

	// Clone StatusTimestamp

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.StatusTimestamp = original.StatusTimestamp

	return clone
}

// Clone creates a deep copy of the ServiceSettings struct
func (original ServiceSettings) Clone() ServiceSettings {
	clone := ServiceSettings{}

	// Clone KeepOnline

	// Simple type - direct assignment
	clone.KeepOnline = original.KeepOnline

	// Clone WppConnectVersion

	// Simple type - direct assignment
	clone.WppConnectVersion = original.WppConnectVersion

	// Clone WaVersion

	// Simple type - direct assignment
	clone.WaVersion = original.WaVersion

	return clone
}

// Clone creates a deep copy of the Service struct
func (original Service) Clone() Service {
	clone := Service{}

	// Clone Id

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.Id = original.Id

	// Clone Name

	// Simple type - direct assignment
	clone.Name = original.Name

	// Clone Data

	// Pointer to struct - create new instance and clone
	if original.Data != nil {
		clonedData := original.Data.Clone()
		clone.Data = &clonedData
	}

	// Clone Settings

	// Pointer to struct - create new instance and clone
	if original.Settings != nil {
		clonedSettings := original.Settings.Clone()
		clone.Settings = &clonedSettings
	}

	// Clone CreatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.CreatedAt = original.CreatedAt

	// Clone UpdatedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.UpdatedAt = original.UpdatedAt

	// Clone DeletedAt

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.DeletedAt = original.DeletedAt

	// Clone AccountId

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.AccountId = original.AccountId

	// Clone ServerPodId

	// Complex type - direct assignment (may need manual handling for deep copy)
	clone.ServerPodId = original.ServerPodId

	// Clone Account

	// Pointer to struct - create new instance and clone
	if original.Account != nil {
		clonedAccount := original.Account.Clone()
		clone.Account = &clonedAccount
	}

	// Clone ServerPod

	// Pointer to struct - create new instance and clone
	if original.ServerPod != nil {
		clonedServerPod := original.ServerPod.Clone()
		clone.ServerPod = &clonedServerPod
	}

	return clone
}
