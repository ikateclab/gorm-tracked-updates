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

// Diff compares this AccountSettings instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *AccountSettings) Diff(b *AccountSettings) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	return diff
}

// Diff compares this AccountData instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *AccountData) Diff(b *AccountData) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	return diff
}

// Diff compares this Account instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *Account) Diff(b *Account) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Id

	// UUID comparison

	// Direct UUID comparison
	if a.Id != b.Id {
		diff["Id"] = b.Id
	}

	// Compare Name

	// Simple type comparison
	if a.Name != b.Name {
		diff["Name"] = b.Name
	}

	// Compare Settings

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if a.Settings == nil && b.Settings != nil {
		// a is nil, b is not nil - use entire b
		jsonValue, err := sonic.Marshal(b.Settings)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
		} else if err != nil {
			diff["Settings"] = b.Settings
		}
	} else if a.Settings != nil && b.Settings == nil {
		// a is not nil, b is nil - set to null
		diff["Settings"] = nil
	} else if a.Settings != nil && b.Settings != nil {
		// Both are not nil - use attribute-by-attribute diff
		SettingsDiff := a.Settings.Diff(b.Settings)
		if len(SettingsDiff) > 0 {
			jsonValue, err := sonic.Marshal(SettingsDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Settings"] = b.Settings
			}
		}
	}

	// Compare Data

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if a.Data == nil && b.Data != nil {
		// a is nil, b is not nil - use entire b
		jsonValue, err := sonic.Marshal(b.Data)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Data"] = gorm.Expr("? || ?", clause.Column{Name: "data"}, string(jsonValue))
		} else if err != nil {
			diff["Data"] = b.Data
		}
	} else if a.Data != nil && b.Data == nil {
		// a is not nil, b is nil - set to null
		diff["Data"] = nil
	} else if a.Data != nil && b.Data != nil {
		// Both are not nil - use attribute-by-attribute diff
		DataDiff := a.Data.Diff(b.Data)
		if len(DataDiff) > 0 {
			jsonValue, err := sonic.Marshal(DataDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Data"] = gorm.Expr("? || ?", clause.Column{Name: "data"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Data"] = b.Data
			}
		}
	}

	// Compare IsActive

	// Simple type comparison
	if a.IsActive != b.IsActive {
		diff["IsActive"] = b.IsActive
	}

	// Compare CorrelationId

	// Simple type comparison
	if a.CorrelationId != b.CorrelationId {
		diff["CorrelationId"] = b.CorrelationId
	}

	// Compare WebhookUrl

	// Simple type comparison
	if a.WebhookUrl != b.WebhookUrl {
		diff["WebhookUrl"] = b.WebhookUrl
	}

	// Compare CreatedAt

	// Time comparison

	// Direct time comparison
	if !a.CreatedAt.Equal(b.CreatedAt) {
		diff["CreatedAt"] = b.CreatedAt

	}

	// Compare UpdatedAt

	// Time comparison

	// Direct time comparison
	if !a.UpdatedAt.Equal(b.UpdatedAt) {
		diff["UpdatedAt"] = b.UpdatedAt

	}

	// Compare DeletedAt

	// GORM DeletedAt comparison
	if a.DeletedAt != b.DeletedAt {
		diff["DeletedAt"] = b.DeletedAt
	}

	// Compare Services

	// Complex type comparison (slice, map, interface, etc.)
	if !reflect.DeepEqual(a.Services, b.Services) {
		diff["Services"] = b.Services
	}

	return diff
}

// Diff compares this ServerPod instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *ServerPod) Diff(b *ServerPod) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Id

	// UUID comparison

	// Direct UUID comparison
	if a.Id != b.Id {
		diff["Id"] = b.Id
	}

	// Compare Name

	// Simple type comparison
	if a.Name != b.Name {
		diff["Name"] = b.Name
	}

	// Compare Address

	// Simple type comparison
	if a.Address != b.Address {
		diff["Address"] = b.Address
	}

	// Compare Version

	// Simple type comparison
	if a.Version != b.Version {
		diff["Version"] = b.Version
	}

	// Compare Settings

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// Use bytes.Equal for datatypes.JSON ([]byte underlying type)
	if !bytes.Equal([]byte(a.Settings), []byte(b.Settings)) {
		jsonValue, err := sonic.Marshal(b.Settings)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
		} else if err != nil {
			// Fallback to regular assignment if JSON marshaling fails
			diff["Settings"] = b.Settings
		}
		// Skip adding to diff if JSON is empty (no-op update)
	}

	// Compare LastPingAt

	// Time comparison

	// Direct time comparison
	if !a.LastPingAt.Equal(b.LastPingAt) {
		diff["LastPingAt"] = b.LastPingAt

	}

	// Compare CreatedAt

	// Time comparison

	// Direct time comparison
	if !a.CreatedAt.Equal(b.CreatedAt) {
		diff["CreatedAt"] = b.CreatedAt

	}

	// Compare UpdatedAt

	// Time comparison

	// Direct time comparison
	if !a.UpdatedAt.Equal(b.UpdatedAt) {
		diff["UpdatedAt"] = b.UpdatedAt

	}

	// Compare DeletedAt

	// GORM DeletedAt comparison
	if a.DeletedAt != b.DeletedAt {
		diff["DeletedAt"] = b.DeletedAt
	}

	// Compare ServerPodTypeId

	// UUID comparison

	// Direct UUID comparison
	if a.ServerPodTypeId != b.ServerPodTypeId {
		diff["ServerPodTypeId"] = b.ServerPodTypeId
	}

	// Compare ServerPodType

	// Struct type comparison - call Diff method directly
	nestedDiff := a.ServerPodType.Diff(&b.ServerPodType)
	if len(nestedDiff) > 0 {
		diff["ServerPodType"] = nestedDiff
	}

	return diff
}

// Diff compares this ServiceVersion instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *ServiceVersion) Diff(b *ServiceVersion) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare WppConnectVersion

	// Simple type comparison
	if a.WppConnectVersion != b.WppConnectVersion {
		diff["WppConnectVersion"] = b.WppConnectVersion
	}

	// Compare WaVersion

	// Simple type comparison
	if a.WaVersion != b.WaVersion {
		diff["WaVersion"] = b.WaVersion
	}

	return diff
}

// Diff compares this ServerPodType instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *ServerPodType) Diff(b *ServerPodType) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Id

	// UUID comparison

	// Direct UUID comparison
	if a.Id != b.Id {
		diff["Id"] = b.Id
	}

	// Compare Name

	// Simple type comparison
	if a.Name != b.Name {
		diff["Name"] = b.Name
	}

	// Compare Version

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if a.Version == nil && b.Version != nil {
		// a is nil, b is not nil - use entire b
		jsonValue, err := sonic.Marshal(b.Version)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Version"] = gorm.Expr("? || ?", clause.Column{Name: "version"}, string(jsonValue))
		} else if err != nil {
			diff["Version"] = b.Version
		}
	} else if a.Version != nil && b.Version == nil {
		// a is not nil, b is nil - set to null
		diff["Version"] = nil
	} else if a.Version != nil && b.Version != nil {
		// Both are not nil - use attribute-by-attribute diff
		VersionDiff := a.Version.Diff(b.Version)
		if len(VersionDiff) > 0 {
			jsonValue, err := sonic.Marshal(VersionDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Version"] = gorm.Expr("? || ?", clause.Column{Name: "version"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Version"] = b.Version
			}
		}
	}

	// Compare AutoScalable

	// Simple type comparison
	if a.AutoScalable != b.AutoScalable {
		diff["AutoScalable"] = b.AutoScalable
	}

	// Compare Cloud

	// Simple type comparison
	if a.Cloud != b.Cloud {
		diff["Cloud"] = b.Cloud
	}

	// Compare ServerSize

	// Simple type comparison
	if a.ServerSize != b.ServerSize {
		diff["ServerSize"] = b.ServerSize
	}

	// Compare MaxPerPod

	// Simple type comparison
	if a.MaxPerPod != b.MaxPerPod {
		diff["MaxPerPod"] = b.MaxPerPod
	}

	// Compare Min

	// Simple type comparison
	if a.Min != b.Min {
		diff["Min"] = b.Min
	}

	// Compare DesiredAvailable

	// Simple type comparison
	if a.DesiredAvailable != b.DesiredAvailable {
		diff["DesiredAvailable"] = b.DesiredAvailable
	}

	// Compare StartPriority

	// Simple type comparison
	if a.StartPriority != b.StartPriority {
		diff["StartPriority"] = b.StartPriority
	}

	// Compare AccountIdWhitelist

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - custom slice types with jsonb storage (not comparable with !=)
	if !reflect.DeepEqual(a.AccountIdWhitelist, b.AccountIdWhitelist) {
		jsonValue, err := sonic.Marshal(b.AccountIdWhitelist)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["AccountIdWhitelist"] = gorm.Expr("? || ?", clause.Column{Name: "account_id_whitelist"}, string(jsonValue))
		} else if err != nil {
			// Fallback to regular assignment if JSON marshaling fails
			diff["AccountIdWhitelist"] = b.AccountIdWhitelist
		}
		// Skip adding to diff if JSON is empty (no-op update)
	}

	// Compare ServiceIdWhitelist

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - custom slice types with jsonb storage (not comparable with !=)
	if !reflect.DeepEqual(a.ServiceIdWhitelist, b.ServiceIdWhitelist) {
		jsonValue, err := sonic.Marshal(b.ServiceIdWhitelist)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["ServiceIdWhitelist"] = gorm.Expr("? || ?", clause.Column{Name: "service_id_whitelist"}, string(jsonValue))
		} else if err != nil {
			// Fallback to regular assignment if JSON marshaling fails
			diff["ServiceIdWhitelist"] = b.ServiceIdWhitelist
		}
		// Skip adding to diff if JSON is empty (no-op update)
	}

	// Compare CreatedAt

	// Time comparison

	// Direct time comparison
	if !a.CreatedAt.Equal(b.CreatedAt) {
		diff["CreatedAt"] = b.CreatedAt

	}

	// Compare UpdatedAt

	// Time comparison

	// Direct time comparison
	if !a.UpdatedAt.Equal(b.UpdatedAt) {
		diff["UpdatedAt"] = b.UpdatedAt

	}

	// Compare DeletedAt

	// GORM DeletedAt comparison
	if a.DeletedAt != b.DeletedAt {
		diff["DeletedAt"] = b.DeletedAt
	}

	return diff
}

// Diff compares this ServiceDataStatus instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *ServiceDataStatus) Diff(b *ServiceDataStatus) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare IsSyncing

	// Simple type comparison
	if a.IsSyncing != b.IsSyncing {
		diff["isSyncing"] = b.IsSyncing
	}

	// Compare IsConnected

	// Simple type comparison
	if a.IsConnected != b.IsConnected {
		diff["isConnected"] = b.IsConnected
	}

	// Compare IsStarting

	// Simple type comparison
	if a.IsStarting != b.IsStarting {
		diff["isStarting"] = b.IsStarting
	}

	// Compare IsStarted

	// Simple type comparison
	if a.IsStarted != b.IsStarted {
		diff["isStarted"] = b.IsStarted
	}

	// Compare IsConflicted

	// Simple type comparison
	if a.IsConflicted != b.IsConflicted {
		diff["isConflicted"] = b.IsConflicted
	}

	// Compare IsLoading

	// Simple type comparison
	if a.IsLoading != b.IsLoading {
		diff["isLoading"] = b.IsLoading
	}

	// Compare IsOnChatPage

	// Simple type comparison
	if a.IsOnChatPage != b.IsOnChatPage {
		diff["isOnChatPage"] = b.IsOnChatPage
	}

	// Compare EnteredQrCodePageAt

	// Simple type comparison
	if a.EnteredQrCodePageAt != b.EnteredQrCodePageAt {
		diff["enteredQrCodePageAt"] = b.EnteredQrCodePageAt
	}

	// Compare DisconnectedAt

	// Simple type comparison
	if a.DisconnectedAt != b.DisconnectedAt {
		diff["disconnectedAt"] = b.DisconnectedAt
	}

	// Compare IsOnQrPage

	// Simple type comparison
	if a.IsOnQrPage != b.IsOnQrPage {
		diff["isOnQrPage"] = b.IsOnQrPage
	}

	// Compare IsQrCodeExpired

	// Simple type comparison
	if a.IsQrCodeExpired != b.IsQrCodeExpired {
		diff["isQrCodeExpired"] = b.IsQrCodeExpired
	}

	// Compare IsWebConnected

	// Simple type comparison
	if a.IsWebConnected != b.IsWebConnected {
		diff["isWebConnected"] = b.IsWebConnected
	}

	// Compare IsWebSyncing

	// Simple type comparison
	if a.IsWebSyncing != b.IsWebSyncing {
		diff["isWebSyncing"] = b.IsWebSyncing
	}

	// Compare Mode

	// Simple type comparison
	if a.Mode != b.Mode {
		diff["mode"] = b.Mode
	}

	// Compare MyId

	// Simple type comparison
	if a.MyId != b.MyId {
		diff["myId"] = b.MyId
	}

	// Compare MyName

	// Simple type comparison
	if a.MyName != b.MyName {
		diff["myName"] = b.MyName
	}

	// Compare MyNumber

	// Simple type comparison
	if a.MyNumber != b.MyNumber {
		diff["myNumber"] = b.MyNumber
	}

	// Compare QrCodeExpiresAt

	// Simple type comparison
	if a.QrCodeExpiresAt != b.QrCodeExpiresAt {
		diff["qrCodeExpiresAt"] = b.QrCodeExpiresAt
	}

	// Compare QrCodeUrl

	// Simple type comparison
	if a.QrCodeUrl != b.QrCodeUrl {
		diff["qrCodeUrl"] = b.QrCodeUrl
	}

	// Compare State

	// Simple type comparison
	if a.State != b.State {
		diff["state"] = b.State
	}

	// Compare WaVersion

	// Simple type comparison
	if a.WaVersion != b.WaVersion {
		diff["waVersion"] = b.WaVersion
	}

	return diff
}

// Diff compares this ServiceData instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *ServiceData) Diff(b *ServiceData) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare MyId

	// Simple type comparison
	if a.MyId != b.MyId {
		diff["myId"] = b.MyId
	}

	// Compare LastSyncAt

	// Time comparison

	// Pointer to time comparison
	if (a.LastSyncAt == nil) != (b.LastSyncAt == nil) || (a.LastSyncAt != nil && !a.LastSyncAt.Equal(*b.LastSyncAt)) {
		diff["lastSyncAt"] = b.LastSyncAt
	}

	// Compare LastMessageTimestamp

	// Time comparison

	// Pointer to time comparison
	if (a.LastMessageTimestamp == nil) != (b.LastMessageTimestamp == nil) || (a.LastMessageTimestamp != nil && !a.LastMessageTimestamp.Equal(*b.LastMessageTimestamp)) {
		diff["lastMessageTimestamp"] = b.LastMessageTimestamp
	}

	// Compare SyncCount

	// Simple type comparison
	if a.SyncCount != b.SyncCount {
		diff["syncCount"] = b.SyncCount
	}

	// Compare SyncFlowDone

	// Simple type comparison
	if a.SyncFlowDone != b.SyncFlowDone {
		diff["syncFlowDone"] = b.SyncFlowDone
	}

	// Compare Status

	// Struct type comparison - call Diff method directly
	nestedDiff := a.Status.Diff(&b.Status)
	if len(nestedDiff) > 0 {
		diff["status"] = nestedDiff
	}

	// Compare StatusTimestamp

	// Time comparison

	// Pointer to time comparison
	if (a.StatusTimestamp == nil) != (b.StatusTimestamp == nil) || (a.StatusTimestamp != nil && !a.StatusTimestamp.Equal(*b.StatusTimestamp)) {
		diff["statusTimestamp"] = b.StatusTimestamp
	}

	return diff
}

// Diff compares this ServiceSettings instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *ServiceSettings) Diff(b *ServiceSettings) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare KeepOnline

	// Simple type comparison
	if a.KeepOnline != b.KeepOnline {
		diff["keepOnline"] = b.KeepOnline
	}

	// Compare WppConnectVersion

	// Simple type comparison
	if a.WppConnectVersion != b.WppConnectVersion {
		diff["wppConnectVersion"] = b.WppConnectVersion
	}

	// Compare WaVersion

	// Simple type comparison
	if a.WaVersion != b.WaVersion {
		diff["waVersion"] = b.WaVersion
	}

	return diff
}

// Diff compares this Service instance with another and returns a map of differences
// with only the new values for fields that have changed.
// Returns nil if either pointer is nil.
func (a *Service) Diff(b *Service) map[string]interface{} {
	// Handle nil pointers
	if a == nil || b == nil {
		return nil
	}

	diff := make(map[string]interface{})

	// Compare Id

	// UUID comparison

	// Direct UUID comparison
	if a.Id != b.Id {
		diff["Id"] = b.Id
	}

	// Compare Name

	// Simple type comparison
	if a.Name != b.Name {
		diff["Name"] = b.Name
	}

	// Compare Data

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if a.Data == nil && b.Data != nil {
		// a is nil, b is not nil - use entire b
		jsonValue, err := sonic.Marshal(b.Data)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Data"] = gorm.Expr("? || ?", clause.Column{Name: "data"}, string(jsonValue))
		} else if err != nil {
			diff["Data"] = b.Data
		}
	} else if a.Data != nil && b.Data == nil {
		// a is not nil, b is nil - set to null
		diff["Data"] = nil
	} else if a.Data != nil && b.Data != nil {
		// Both are not nil - use attribute-by-attribute diff
		DataDiff := a.Data.Diff(b.Data)
		if len(DataDiff) > 0 {
			jsonValue, err := sonic.Marshal(DataDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Data"] = gorm.Expr("? || ?", clause.Column{Name: "data"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Data"] = b.Data
			}
		}
	}

	// Compare Settings

	// JSON field comparison - handle both datatypes.JSON and struct types with jsonb storage

	// JSON field comparison - attribute-by-attribute diff for struct types

	// Handle pointer to struct
	if a.Settings == nil && b.Settings != nil {
		// a is nil, b is not nil - use entire b
		jsonValue, err := sonic.Marshal(b.Settings)
		if err == nil && !isEmptyJSON(string(jsonValue)) {
			diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
		} else if err != nil {
			diff["Settings"] = b.Settings
		}
	} else if a.Settings != nil && b.Settings == nil {
		// a is not nil, b is nil - set to null
		diff["Settings"] = nil
	} else if a.Settings != nil && b.Settings != nil {
		// Both are not nil - use attribute-by-attribute diff
		SettingsDiff := a.Settings.Diff(b.Settings)
		if len(SettingsDiff) > 0 {
			jsonValue, err := sonic.Marshal(SettingsDiff)
			if err == nil && !isEmptyJSON(string(jsonValue)) {
				diff["Settings"] = gorm.Expr("? || ?", clause.Column{Name: "settings"}, string(jsonValue))
			} else if err != nil {
				// Fallback to regular assignment if JSON marshaling fails
				diff["Settings"] = b.Settings
			}
		}
	}

	// Compare CreatedAt

	// Time comparison

	// Direct time comparison
	if !a.CreatedAt.Equal(b.CreatedAt) {
		diff["CreatedAt"] = b.CreatedAt

	}

	// Compare UpdatedAt

	// Time comparison

	// Direct time comparison
	if !a.UpdatedAt.Equal(b.UpdatedAt) {
		diff["UpdatedAt"] = b.UpdatedAt

	}

	// Compare DeletedAt

	// GORM DeletedAt comparison
	if a.DeletedAt != b.DeletedAt {
		diff["DeletedAt"] = b.DeletedAt
	}

	// Compare AccountId

	// UUID comparison

	// Direct UUID comparison
	if a.AccountId != b.AccountId {
		diff["AccountId"] = b.AccountId
	}

	// Compare ServerPodId

	// UUID comparison

	// Pointer to UUID comparison
	if (a.ServerPodId == nil) != (b.ServerPodId == nil) || (a.ServerPodId != nil && *a.ServerPodId != *b.ServerPodId) {
		diff["ServerPodId"] = b.ServerPodId
	}

	// Compare Account

	// Pointer to struct comparison
	if a.Account == nil || b.Account == nil {
		if a.Account != b.Account {
			diff["Account"] = b.Account
		}
	} else {
		nestedDiff := a.Account.Diff(b.Account)
		if len(nestedDiff) > 0 {
			diff["Account"] = nestedDiff
		}
	}

	// Compare ServerPod

	// Pointer to struct comparison
	if a.ServerPod == nil || b.ServerPod == nil {
		if a.ServerPod != b.ServerPod {
			diff["ServerPod"] = b.ServerPod
		}
	} else {
		nestedDiff := a.ServerPod.Diff(b.ServerPod)
		if len(nestedDiff) > 0 {
			diff["ServerPod"] = nestedDiff
		}
	}

	return diff
}
