package models

import (
	"bytes"
	"github.com/bytedance/sonic"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"reflect"
	"strings"
)

// isEmptyJSON checks if a JSON string represents an empty object or array
func isEmptyJSON(jsonStr string) bool {
	trimmed := strings.TrimSpace(jsonStr)
	return trimmed == "{}" || trimmed == "[]" || trimmed == "null"
}

// Diff compares this AccountSettings instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *AccountSettings) Diff(old *AccountSettings) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	return diff
}

// Diff compares this AccountData instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *AccountData) Diff(old *AccountData) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	return diff
}

// Diff compares this Account instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *Account) Diff(old *Account) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Id

	// UUID comparison

	// Direct UUID comparison
	if new.Id != old.Id {
		diff["Id"] = new.Id
	}

	// Compare Name

	// Simple type comparison
	if new.Name != old.Name {
		diff["Name"] = new.Name
	}

	// Compare Settings

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if new.Settings == nil && old.Settings != nil {
		// new is nil, old is not nil - set to null
		diff["Settings"] = nil
	} else if new.Settings != nil && old.Settings == nil {
		// new is not nil, old is nil - use entire new
		jsonValue, err := sonic.Marshal(new.Settings)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
		} else if err != nil {
			diff["Settings"] = new.Settings
		}
	} else if new.Settings != nil && old.Settings != nil {
		// Both are not nil - use attribute-by-attribute diff
		SettingsDiff := new.Settings.Diff(old.Settings)
		if len(SettingsDiff) > 0 {
			jsonValue, err := sonic.Marshal(SettingsDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Settings"] = new.Settings
			}
		}
	}

	// Compare Data

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if new.Data == nil && old.Data != nil {
		// new is nil, old is not nil - set to null
		diff["Data"] = nil
	} else if new.Data != nil && old.Data == nil {
		// new is not nil, old is nil - use entire new
		jsonValue, err := sonic.Marshal(new.Data)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Data"] = gorm.Expr("? || ?", clause.Column{Name: "data"}, string(jsonValue))
		} else if err != nil {
			diff["Data"] = new.Data
		}
	} else if new.Data != nil && old.Data != nil {
		// Both are not nil - use attribute-by-attribute diff
		DataDiff := new.Data.Diff(old.Data)
		if len(DataDiff) > 0 {
			jsonValue, err := sonic.Marshal(DataDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Data"] = gorm.Expr("? || ?", clause.Column{Name: "data"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Data"] = new.Data
			}
		}
	}

	// Compare IsActive

	// Simple type comparison
	if new.IsActive != old.IsActive {
		diff["IsActive"] = new.IsActive
	}

	// Compare CorrelationId

	// Simple type comparison
	if new.CorrelationId != old.CorrelationId {
		diff["CorrelationId"] = new.CorrelationId
	}

	// Compare WebhookUrl

	// Simple type comparison
	if new.WebhookUrl != old.WebhookUrl {
		diff["WebhookUrl"] = new.WebhookUrl
	}

	// Compare CreatedAt

	// Time comparison

	// Direct time comparison
	if !new.CreatedAt.Equal(old.CreatedAt) {
		diff["CreatedAt"] = new.CreatedAt

	}

	// Compare UpdatedAt

	// Time comparison

	// Direct time comparison
	if !new.UpdatedAt.Equal(old.UpdatedAt) {
		diff["UpdatedAt"] = new.UpdatedAt

	}

	// Compare DeletedAt

	// GORM DeletedAt comparison
	if new.DeletedAt != old.DeletedAt {
		diff["DeletedAt"] = new.DeletedAt
	}

	// Compare Services

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.Services, old.Services) {
		diff["Services"] = new.Services
	}

	return diff
}

// Diff compares this ServerPod instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *ServerPod) Diff(old *ServerPod) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Id

	// UUID comparison

	// Direct UUID comparison
	if new.Id != old.Id {
		diff["Id"] = new.Id
	}

	// Compare Name

	// Simple type comparison
	if new.Name != old.Name {
		diff["Name"] = new.Name
	}

	// Compare Address

	// Simple type comparison
	if new.Address != old.Address {
		diff["Address"] = new.Address
	}

	// Compare Version

	// Simple type comparison
	if new.Version != old.Version {
		diff["Version"] = new.Version
	}

	// Compare Settings

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// Use bytes.Equal for datatypes.JSON ([]byte underlying type)
	if !bytes.Equal([]byte(new.Settings), []byte(old.Settings)) {
		jsonValue, err := sonic.Marshal(new.Settings)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
		} else if err != nil {
			// Fallback to regular assignment if JSON marshaling fails
			diff["Settings"] = new.Settings
		}
		// Skip adding to diff if JSON is empty (no-op update)
	}

	// Compare LastPingAt

	// Time comparison

	// Direct time comparison
	if !new.LastPingAt.Equal(old.LastPingAt) {
		diff["LastPingAt"] = new.LastPingAt

	}

	// Compare CreatedAt

	// Time comparison

	// Direct time comparison
	if !new.CreatedAt.Equal(old.CreatedAt) {
		diff["CreatedAt"] = new.CreatedAt

	}

	// Compare UpdatedAt

	// Time comparison

	// Direct time comparison
	if !new.UpdatedAt.Equal(old.UpdatedAt) {
		diff["UpdatedAt"] = new.UpdatedAt

	}

	// Compare DeletedAt

	// GORM DeletedAt comparison
	if new.DeletedAt != old.DeletedAt {
		diff["DeletedAt"] = new.DeletedAt
	}

	// Compare ServerPodTypeId

	// UUID comparison

	// Direct UUID comparison
	if new.ServerPodTypeId != old.ServerPodTypeId {
		diff["ServerPodTypeId"] = new.ServerPodTypeId
	}

	// Compare ServerPodType

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(new.ServerPodType, old.ServerPodType) {
		diff["ServerPodType"] = new.ServerPodType
	}

	return diff
}

// Diff compares this ServiceVersion instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *ServiceVersion) Diff(old *ServiceVersion) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare WppConnectVersion

	// Simple type comparison
	if new.WppConnectVersion != old.WppConnectVersion {
		diff["WppConnectVersion"] = new.WppConnectVersion
	}

	// Compare WaVersion

	// Simple type comparison
	if new.WaVersion != old.WaVersion {
		diff["WaVersion"] = new.WaVersion
	}

	return diff
}

// Diff compares this ServerPodType instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *ServerPodType) Diff(old *ServerPodType) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Id

	// UUID comparison

	// Direct UUID comparison
	if new.Id != old.Id {
		diff["Id"] = new.Id
	}

	// Compare Name

	// Simple type comparison
	if new.Name != old.Name {
		diff["Name"] = new.Name
	}

	// Compare Version

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if new.Version == nil && old.Version != nil {
		// new is nil, old is not nil - set to null
		diff["Version"] = nil
	} else if new.Version != nil && old.Version == nil {
		// new is not nil, old is nil - use entire new
		jsonValue, err := sonic.Marshal(new.Version)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Version"] = gorm.Expr("? || ?", clause.Column{Name: "version"}, string(jsonValue))
		} else if err != nil {
			diff["Version"] = new.Version
		}
	} else if new.Version != nil && old.Version != nil {
		// Both are not nil - use attribute-by-attribute diff
		VersionDiff := new.Version.Diff(old.Version)
		if len(VersionDiff) > 0 {
			jsonValue, err := sonic.Marshal(VersionDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Version"] = gorm.Expr("? || ?", clause.Column{Name: "version"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Version"] = new.Version
			}
		}
	}

	// Compare AutoScalable

	// Simple type comparison
	if new.AutoScalable != old.AutoScalable {
		diff["AutoScalable"] = new.AutoScalable
	}

	// Compare Cloud

	// Simple type comparison
	if new.Cloud != old.Cloud {
		diff["Cloud"] = new.Cloud
	}

	// Compare ServerSize

	// Simple type comparison
	if new.ServerSize != old.ServerSize {
		diff["ServerSize"] = new.ServerSize
	}

	// Compare MaxPerPod

	// Simple type comparison
	if new.MaxPerPod != old.MaxPerPod {
		diff["MaxPerPod"] = new.MaxPerPod
	}

	// Compare Min

	// Simple type comparison
	if new.Min != old.Min {
		diff["Min"] = new.Min
	}

	// Compare DesiredAvailable

	// Simple type comparison
	if new.DesiredAvailable != old.DesiredAvailable {
		diff["DesiredAvailable"] = new.DesiredAvailable
	}

	// Compare StartPriority

	// Simple type comparison
	if new.StartPriority != old.StartPriority {
		diff["StartPriority"] = new.StartPriority
	}

	// Compare AccountIdWhitelist

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - custom slice types with jsonb storage (not comparable with !=)
	if !reflect.DeepEqual(new.AccountIdWhitelist, old.AccountIdWhitelist) {
		jsonValue, err := sonic.Marshal(new.AccountIdWhitelist)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["AccountIdWhitelist"] = gorm.Expr("? || ?", clause.Column{Name: "account_id_whitelist"}, string(jsonValue))
		} else if err != nil {
			// Fallback to regular assignment if JSON marshaling fails
			diff["AccountIdWhitelist"] = new.AccountIdWhitelist
		}
		// Skip adding to diff if JSON is empty (no-op update)
	}

	// Compare ServiceIdWhitelist

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - custom slice types with jsonb storage (not comparable with !=)
	if !reflect.DeepEqual(new.ServiceIdWhitelist, old.ServiceIdWhitelist) {
		jsonValue, err := sonic.Marshal(new.ServiceIdWhitelist)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["ServiceIdWhitelist"] = gorm.Expr("? || ?", clause.Column{Name: "service_id_whitelist"}, string(jsonValue))
		} else if err != nil {
			// Fallback to regular assignment if JSON marshaling fails
			diff["ServiceIdWhitelist"] = new.ServiceIdWhitelist
		}
		// Skip adding to diff if JSON is empty (no-op update)
	}

	// Compare CreatedAt

	// Time comparison

	// Direct time comparison
	if !new.CreatedAt.Equal(old.CreatedAt) {
		diff["CreatedAt"] = new.CreatedAt

	}

	// Compare UpdatedAt

	// Time comparison

	// Direct time comparison
	if !new.UpdatedAt.Equal(old.UpdatedAt) {
		diff["UpdatedAt"] = new.UpdatedAt

	}

	// Compare DeletedAt

	// GORM DeletedAt comparison
	if new.DeletedAt != old.DeletedAt {
		diff["DeletedAt"] = new.DeletedAt
	}

	return diff
}

// Diff compares this ServiceDataStatus instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *ServiceDataStatus) Diff(old *ServiceDataStatus) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare IsSyncing

	// Simple type comparison
	if new.IsSyncing != old.IsSyncing {
		diff["isSyncing"] = new.IsSyncing
	}

	// Compare IsConnected

	// Simple type comparison
	if new.IsConnected != old.IsConnected {
		diff["isConnected"] = new.IsConnected
	}

	// Compare IsStarting

	// Simple type comparison
	if new.IsStarting != old.IsStarting {
		diff["isStarting"] = new.IsStarting
	}

	// Compare IsStarted

	// Simple type comparison
	if new.IsStarted != old.IsStarted {
		diff["isStarted"] = new.IsStarted
	}

	// Compare IsConflicted

	// Simple type comparison
	if new.IsConflicted != old.IsConflicted {
		diff["isConflicted"] = new.IsConflicted
	}

	// Compare IsLoading

	// Simple type comparison
	if new.IsLoading != old.IsLoading {
		diff["isLoading"] = new.IsLoading
	}

	// Compare IsOnChatPage

	// Simple type comparison
	if new.IsOnChatPage != old.IsOnChatPage {
		diff["isOnChatPage"] = new.IsOnChatPage
	}

	// Compare EnteredQrCodePageAt

	// Simple type comparison
	if new.EnteredQrCodePageAt != old.EnteredQrCodePageAt {
		diff["enteredQrCodePageAt"] = new.EnteredQrCodePageAt
	}

	// Compare DisconnectedAt

	// Simple type comparison
	if new.DisconnectedAt != old.DisconnectedAt {
		diff["disconnectedAt"] = new.DisconnectedAt
	}

	// Compare IsOnQrPage

	// Simple type comparison
	if new.IsOnQrPage != old.IsOnQrPage {
		diff["isOnQrPage"] = new.IsOnQrPage
	}

	// Compare IsQrCodeExpired

	// Simple type comparison
	if new.IsQrCodeExpired != old.IsQrCodeExpired {
		diff["isQrCodeExpired"] = new.IsQrCodeExpired
	}

	// Compare IsWebConnected

	// Simple type comparison
	if new.IsWebConnected != old.IsWebConnected {
		diff["isWebConnected"] = new.IsWebConnected
	}

	// Compare IsWebSyncing

	// Simple type comparison
	if new.IsWebSyncing != old.IsWebSyncing {
		diff["isWebSyncing"] = new.IsWebSyncing
	}

	// Compare Mode

	// Simple type comparison
	if new.Mode != old.Mode {
		diff["mode"] = new.Mode
	}

	// Compare MyId

	// Simple type comparison
	if new.MyId != old.MyId {
		diff["myId"] = new.MyId
	}

	// Compare MyName

	// Simple type comparison
	if new.MyName != old.MyName {
		diff["myName"] = new.MyName
	}

	// Compare MyNumber

	// Simple type comparison
	if new.MyNumber != old.MyNumber {
		diff["myNumber"] = new.MyNumber
	}

	// Compare QrCodeExpiresAt

	// Simple type comparison
	if new.QrCodeExpiresAt != old.QrCodeExpiresAt {
		diff["qrCodeExpiresAt"] = new.QrCodeExpiresAt
	}

	// Compare QrCodeUrl

	// Simple type comparison
	if new.QrCodeUrl != old.QrCodeUrl {
		diff["qrCodeUrl"] = new.QrCodeUrl
	}

	// Compare State

	// Simple type comparison
	if new.State != old.State {
		diff["state"] = new.State
	}

	// Compare WaVersion

	// Simple type comparison
	if new.WaVersion != old.WaVersion {
		diff["waVersion"] = new.WaVersion
	}

	return diff
}

// Diff compares this ServiceData instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *ServiceData) Diff(old *ServiceData) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare MyId

	// Simple type comparison
	if new.MyId != old.MyId {
		diff["myId"] = new.MyId
	}

	// Compare LastSyncAt

	// Time comparison

	// Pointer to time comparison
	if (new.LastSyncAt == nil) != (old.LastSyncAt == nil) || (new.LastSyncAt != nil && !new.LastSyncAt.Equal(*old.LastSyncAt)) {
		diff["lastSyncAt"] = new.LastSyncAt
	}

	// Compare LastMessageTimestamp

	// Time comparison

	// Pointer to time comparison
	if (new.LastMessageTimestamp == nil) != (old.LastMessageTimestamp == nil) || (new.LastMessageTimestamp != nil && !new.LastMessageTimestamp.Equal(*old.LastMessageTimestamp)) {
		diff["lastMessageTimestamp"] = new.LastMessageTimestamp
	}

	// Compare SyncCount

	// Simple type comparison
	if new.SyncCount != old.SyncCount {
		diff["syncCount"] = new.SyncCount
	}

	// Compare SyncFlowDone

	// Simple type comparison
	if new.SyncFlowDone != old.SyncFlowDone {
		diff["syncFlowDone"] = new.SyncFlowDone
	}

	// Compare Status

	// Struct type comparison - call Diff method directly
	nestedDiff := new.Status.Diff(&old.Status)
	if len(nestedDiff) > 0 {
		diff["status"] = nestedDiff
	}

	// Compare StatusTimestamp

	// Time comparison

	// Pointer to time comparison
	if (new.StatusTimestamp == nil) != (old.StatusTimestamp == nil) || (new.StatusTimestamp != nil && !new.StatusTimestamp.Equal(*old.StatusTimestamp)) {
		diff["statusTimestamp"] = new.StatusTimestamp
	}

	return diff
}

// Diff compares this ServiceSettings instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *ServiceSettings) Diff(old *ServiceSettings) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare KeepOnline

	// Simple type comparison
	if new.KeepOnline != old.KeepOnline {
		diff["keepOnline"] = new.KeepOnline
	}

	// Compare WppConnectVersion

	// Simple type comparison
	if new.WppConnectVersion != old.WppConnectVersion {
		diff["wppConnectVersion"] = new.WppConnectVersion
	}

	// Compare WaVersion

	// Simple type comparison
	if new.WaVersion != old.WaVersion {
		diff["waVersion"] = new.WaVersion
	}

	return diff
}

// Diff compares this Service instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *Service) Diff(old *Service) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Id

	// UUID comparison

	// Direct UUID comparison
	if new.Id != old.Id {
		diff["Id"] = new.Id
	}

	// Compare Name

	// Simple type comparison
	if new.Name != old.Name {
		diff["Name"] = new.Name
	}

	// Compare Data

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if new.Data == nil && old.Data != nil {
		// new is nil, old is not nil - set to null
		diff["Data"] = nil
	} else if new.Data != nil && old.Data == nil {
		// new is not nil, old is nil - use entire new
		jsonValue, err := sonic.Marshal(new.Data)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Data"] = gorm.Expr("? || ?", clause.Column{Name: "data"}, string(jsonValue))
		} else if err != nil {
			diff["Data"] = new.Data
		}
	} else if new.Data != nil && old.Data != nil {
		// Both are not nil - use attribute-by-attribute diff
		DataDiff := new.Data.Diff(old.Data)
		if len(DataDiff) > 0 {
			jsonValue, err := sonic.Marshal(DataDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Data"] = gorm.Expr("? || ?", clause.Column{Name: "data"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Data"] = new.Data
			}
		}
	}

	// Compare Settings

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if new.Settings == nil && old.Settings != nil {
		// new is nil, old is not nil - set to null
		diff["Settings"] = nil
	} else if new.Settings != nil && old.Settings == nil {
		// new is not nil, old is nil - use entire new
		jsonValue, err := sonic.Marshal(new.Settings)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
		} else if err != nil {
			diff["Settings"] = new.Settings
		}
	} else if new.Settings != nil && old.Settings != nil {
		// Both are not nil - use attribute-by-attribute diff
		SettingsDiff := new.Settings.Diff(old.Settings)
		if len(SettingsDiff) > 0 {
			jsonValue, err := sonic.Marshal(SettingsDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Settings"] = new.Settings
			}
		}
	}

	// Compare CreatedAt

	// Time comparison

	// Direct time comparison
	if !new.CreatedAt.Equal(old.CreatedAt) {
		diff["CreatedAt"] = new.CreatedAt

	}

	// Compare UpdatedAt

	// Time comparison

	// Direct time comparison
	if !new.UpdatedAt.Equal(old.UpdatedAt) {
		diff["UpdatedAt"] = new.UpdatedAt

	}

	// Compare DeletedAt

	// GORM DeletedAt comparison
	if new.DeletedAt != old.DeletedAt {
		diff["DeletedAt"] = new.DeletedAt
	}

	// Compare AccountId

	// UUID comparison

	// Direct UUID comparison
	if new.AccountId != old.AccountId {
		diff["AccountId"] = new.AccountId
	}

	// Compare ServerPodId

	// UUID comparison

	// Pointer to UUID comparison
	if (new.ServerPodId == nil) != (old.ServerPodId == nil) || (new.ServerPodId != nil && *new.ServerPodId != *old.ServerPodId) {
		diff["ServerPodId"] = new.ServerPodId
	}

	// Compare Account

	// Comparable type comparison
	if new.Account != old.Account {
		diff["Account"] = new.Account
	}

	// Compare ServerPod

	// Comparable type comparison
	if new.ServerPod != old.ServerPod {
		diff["ServerPod"] = new.ServerPod
	}

	return diff
}

// Diff compares this Tag instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *Tag) Diff(old *Tag) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Name

	// Simple type comparison
	if new.Name != old.Name {
		diff["name"] = new.Name
	}

	// Compare Value

	// Simple type comparison
	if new.Value != old.Value {
		diff["value"] = new.Value
	}

	return diff
}

// Diff compares this Item instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *Item) Diff(old *Item) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare ID

	// Simple type comparison
	if new.ID != old.ID {
		diff["id"] = new.ID
	}

	// Compare Title

	// Simple type comparison
	if new.Title != old.Title {
		diff["title"] = new.Title
	}

	// Compare Price

	// Simple type comparison
	if new.Price != old.Price {
		diff["price"] = new.Price
	}

	return diff
}

// Diff compares this SimpleModel instance (new) with another (old) and returns a map of differences
// with only the new values for fields that have changed.
// Usage: newValues = new.Diff(old)
// Returns nil if either pointer is nil.
func (new *SimpleModel) Diff(old *SimpleModel) map[string]interface{} {
	// Handle nil pointers
	if new == nil || old == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare ID

	// UUID comparison

	// Direct UUID comparison
	if new.ID != old.ID {
		diff["ID"] = new.ID
	}

	// Compare Name

	// Simple type comparison
	if new.Name != old.Name {
		diff["Name"] = new.Name
	}

	// Compare Tags

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - custom slice types with jsonb storage (not comparable with !=)
	if !reflect.DeepEqual(new.Tags, old.Tags) {
		jsonValue, err := sonic.Marshal(new.Tags)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Tags"] = gorm.Expr("? || ?", clause.Column{Name: "tags"}, string(jsonValue))
		} else if err != nil {
			// Fallback to regular assignment if JSON marshaling fails
			diff["Tags"] = new.Tags
		}
		// Skip adding to diff if JSON is empty (no-op update)
	}

	// Compare Items

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - custom slice types with jsonb storage (not comparable with !=)
	if !reflect.DeepEqual(new.Items, old.Items) {
		jsonValue, err := sonic.Marshal(new.Items)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Items"] = gorm.Expr("? || ?", clause.Column{Name: "items"}, string(jsonValue))
		} else if err != nil {
			// Fallback to regular assignment if JSON marshaling fails
			diff["Items"] = new.Items
		}
		// Skip adding to diff if JSON is empty (no-op update)
	}

	return diff
}
