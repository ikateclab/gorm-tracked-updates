package diffgen

import (
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Test models for nested JSON structures
// ServiceDataStatus represents nested status within ServiceData
// @jsonb
type TestServiceDataStatus struct {
	IsSyncing           bool   `json:"isSyncing,omitempty"`
	IsConnected         bool   `json:"isConnected,omitempty"`
	IsStarting          bool   `json:"isStarting,omitempty"`
	IsStarted           bool   `json:"isStarted,omitempty"`
	IsConflicted        bool   `json:"isConflicted,omitempty"`
	IsLoading           bool   `json:"isLoading,omitempty"`
	IsOnChatPage        bool   `json:"isOnChatPage,omitempty"`
	EnteredQrCodePageAt string `json:"enteredQrCodePageAt,omitempty"`
	DisconnectedAt      string `json:"disconnectedAt,omitempty"`
	IsOnQrPage          bool   `json:"isOnQrPage,omitempty"`
	IsQrCodeExpired     bool   `json:"isQrCodeExpired,omitempty"`
	IsWebConnected      bool   `json:"isWebConnected,omitempty"`
	IsWebSyncing        bool   `json:"isWebSyncing,omitempty"`
	Mode                string `json:"mode,omitempty"`
	MyId                string `json:"myId,omitempty"`
	MyName              string `json:"myName,omitempty"`
	MyNumber            string `json:"myNumber,omitempty"`
	QrCodeExpiresAt     string `json:"qrCodeExpiresAt,omitempty"`
	QrCodeUrl           string `json:"qrCodeUrl,omitempty"`
	State               string `json:"state,omitempty"`
	WaVersion           string `json:"waVersion,omitempty"`
}

// ServiceData represents service data stored in JSONB
// @jsonb
type TestServiceData struct {
	MyId                 string                 `json:"myId,omitempty"`
	LastSyncAt           *time.Time             `json:"lastSyncAt,omitempty"`
	LastMessageTimestamp *time.Time             `json:"lastMessageTimestamp,omitempty"`
	SyncCount            int                    `json:"syncCount,omitempty"`
	SyncFlowDone         bool                   `json:"syncFlowDone,omitempty"`
	Status               TestServiceDataStatus  `json:"status,omitempty"`
	StatusTimestamp      *time.Time             `json:"statusTimestamp,omitempty"`
}

// ServiceSettings represents service settings stored in JSONB
// @jsonb
type TestServiceSettings struct {
	KeepOnline        bool   `json:"keepOnline,omitempty"`
	WppConnectVersion string `json:"wppConnectVersion,omitempty"`
	WaVersion         string `json:"waVersion,omitempty"`
}

type TestService struct {
	Id          uuid.UUID            `gorm:"type:uuid;default:gen_random_uuid();primary_key"`
	Name        string               `gorm:"type:string;not null"`
	Data        *TestServiceData     `gorm:"type:jsonb;not null;default:'{}';serializer:json"`
	Settings    *TestServiceSettings `gorm:"type:jsonb;not null;default:'{}';serializer:json"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

func TestNestedJSONStructGeneration(t *testing.T) {
	generator := New()

	// Create test structs
	structs := []StructInfo{
		{
			Name:    "TestServiceDataStatus",
			Package: "diffgen",
			IsJSONB: true,
			Fields: []StructField{
				{Name: "IsSyncing", Type: "bool", FieldType: FieldTypeSimple, DiffKey: "isSyncing"},
				{Name: "IsConnected", Type: "bool", FieldType: FieldTypeSimple, DiffKey: "isConnected"},
				{Name: "Mode", Type: "string", FieldType: FieldTypeSimple, DiffKey: "mode"},
			},
		},
		{
			Name:    "TestServiceData",
			Package: "diffgen",
			IsJSONB: true,
			Fields: []StructField{
				{Name: "MyId", Type: "string", FieldType: FieldTypeSimple, DiffKey: "myId"},
				{Name: "SyncCount", Type: "int", FieldType: FieldTypeSimple, DiffKey: "syncCount"},
				{Name: "Status", Type: "TestServiceDataStatus", FieldType: FieldTypeStruct, DiffKey: "status"},
			},
		},
		{
			Name:    "TestService",
			Package: "diffgen",
			Fields: []StructField{
				{Name: "Id", Type: "uuid.UUID", FieldType: FieldTypeUUID, DiffKey: "Id"},
				{Name: "Name", Type: "string", FieldType: FieldTypeSimple, DiffKey: "Name"},
				{Name: "Data", Type: "*TestServiceData", FieldType: FieldTypeJSON, Tag: `gorm:"type:jsonb;serializer:json"`, DiffKey: "Data"},
			},
		},
	}

	generator.Structs = structs
	generator.JSONBStructs = map[string]bool{
		"TestServiceDataStatus": true,
		"TestServiceData":       true,
	}

	code, err := generator.GenerateCode()
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Test 1: Verify nested struct uses regular struct comparison, not JSON
	if !strings.Contains(code, "nestedDiff := new.Status.Diff(&old.Status)") {
		t.Error("Expected nested struct to use regular struct diff method")
	}

	// Test 2: Verify nested struct does NOT use gorm.Expr
	statusDiffSection := extractFunctionCode(code, "TestServiceData", "Status")
	if strings.Contains(statusDiffSection, "gorm.Expr") {
		t.Error("Nested struct should not use gorm.Expr - this would create nested gorm.Expr calls")
	}

	// Test 3: Verify root JSONB field DOES use gorm.Expr
	dataDiffSection := extractFunctionCode(code, "TestService", "Data")
	if !strings.Contains(dataDiffSection, "gorm.Expr") {
		t.Error("Root JSONB field should use gorm.Expr")
	}

	// Test 4: Verify proper JSON marshaling structure
	// The test uses uppercase field names since DiffKey computation happens later
	if !strings.Contains(code, `diff["Status"] = nestedDiff`) {
		t.Error("Expected nested diff to be assigned as plain value, not wrapped in gorm.Expr")
	}
}

func TestNestedJSONFieldTypeDetection(t *testing.T) {
	generator := New()
	generator.JSONBStructs = map[string]bool{
		"TestServiceDataStatus": true,
		"TestServiceData":       true,
	}

	tests := []struct {
		name        string
		fieldType   string
		fieldTag    string
		expected    FieldType
		description string
	}{
		{
			name:        "Root JSONB field with gorm tag",
			fieldType:   "*TestServiceData",
			fieldTag:    `gorm:"type:jsonb;serializer:json"`,
			expected:    FieldTypeJSON,
			description: "Should be JSON type because it has database JSON tags",
		},
		{
			name:        "Nested JSONB struct without gorm tag",
			fieldType:   "TestServiceDataStatus",
			fieldTag:    `json:"status,omitempty"`,
			expected:    FieldTypeStruct,
			description: "Should be Struct type to avoid nested gorm.Expr",
		},
		{
			name:        "Nested JSONB struct pointer without gorm tag",
			fieldType:   "*TestServiceDataStatus",
			fieldTag:    `json:"status,omitempty"`,
			expected:    FieldTypeStructPtr,
			description: "Should be StructPtr type to avoid nested gorm.Expr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := generator.determineFieldType(nil, tt.fieldType, tt.fieldTag)
			if result != tt.expected {
				t.Errorf("Field type detection failed for %s: expected %v, got %v. %s",
					tt.name, tt.expected, result, tt.description)
			}
		})
	}
}

func TestNestedJSONIntegration(t *testing.T) {
	// Test the actual generated diff methods work correctly
	generator := New()

	// Parse the test structs
	err := generator.ParseFile("nested_json_test.go")
	if err != nil {
		t.Fatalf("Failed to parse test file: %v", err)
	}

	// Generate code
	code, err := generator.GenerateCode()
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Verify the generated code structure
	t.Run("Generated code structure", func(t *testing.T) {
		// Should have diff functions for all structs
		expectedFunctions := []string{
			"func (new *TestServiceDataStatus) Diff(",
			"func (new *TestServiceData) Diff(",
			"func (new *TestService) Diff(",
		}

		for _, expected := range expectedFunctions {
			if !strings.Contains(code, expected) {
				t.Errorf("Generated code missing expected function: %s", expected)
			}
		}
	})

	t.Run("Nested struct handling", func(t *testing.T) {
		// Nested struct should use struct comparison
		if !strings.Contains(code, "nestedDiff := new.Status.Diff(&old.Status)") {
			t.Error("Nested struct should use struct diff method")
		}

		// Should assign nested diff as plain value
		if !strings.Contains(code, `diff["status"] = nestedDiff`) {
			t.Error("Nested diff should be assigned as plain value")
		}
	})

	t.Run("Root JSONB field handling", func(t *testing.T) {
		// Root JSONB field should use gorm.Expr
		if !strings.Contains(code, `gorm.Expr("? || ?", clause.Column{Name: "data"}`) {
			t.Error("Root JSONB field should use gorm.Expr with proper column name")
		}

		// Should marshal nested diff to JSON
		if !strings.Contains(code, "sonic.Marshal(DataDiff)") {
			t.Error("Root JSONB field should marshal nested diff to JSON")
		}
	})
}

func TestPreventNestedGormExpr(t *testing.T) {
	// This test specifically verifies that we don't generate nested gorm.Expr calls
	generator := New()
	generator.JSONBStructs = map[string]bool{
		"TestServiceDataStatus": true,
		"TestServiceData":       true,
	}

	// Create a scenario that would previously cause nested gorm.Expr
	structs := []StructInfo{
		{
			Name:    "TestServiceDataStatus",
			Package: "diffgen",
			IsJSONB: true,
			Fields: []StructField{
				{Name: "IsConnected", Type: "bool", FieldType: FieldTypeSimple, DiffKey: "isConnected"},
			},
		},
		{
			Name:    "TestServiceData",
			Package: "diffgen",
			IsJSONB: true,
			Fields: []StructField{
				{Name: "Status", Type: "TestServiceDataStatus", FieldType: FieldTypeStruct, DiffKey: "status"},
			},
		},
		{
			Name:    "TestService",
			Package: "diffgen",
			Fields: []StructField{
				{Name: "Data", Type: "*TestServiceData", FieldType: FieldTypeJSON, Tag: `gorm:"type:jsonb;serializer:json"`, DiffKey: "Data"},
			},
		},
	}

	generator.Structs = structs
	code, err := generator.GenerateCode()
	if err != nil {
		t.Fatalf("Failed to generate code: %v", err)
	}

	// Count gorm.Expr occurrences in the nested structure
	statusFuncCode := extractFunctionCode(code, "TestServiceDataStatus", "IsConnected")
	dataFuncCode := extractFunctionCode(code, "TestServiceData", "Status")
	serviceFuncCode := extractFunctionCode(code, "TestService", "Data")

	// Nested struct functions should NOT contain gorm.Expr
	if strings.Contains(statusFuncCode, "gorm.Expr") {
		t.Error("TestServiceDataStatus.Diff should not contain gorm.Expr")
	}

	if strings.Contains(dataFuncCode, "gorm.Expr") {
		t.Error("TestServiceData.Diff Status field should not contain gorm.Expr")
	}

	// Only the root JSONB field should contain gorm.Expr
	if !strings.Contains(serviceFuncCode, "gorm.Expr") {
		t.Error("TestService.Diff Data field should contain gorm.Expr")
	}

	// Verify there's only one level of gorm.Expr (no nesting)
	gormExprCount := strings.Count(serviceFuncCode, "gorm.Expr")
	if gormExprCount != 2 { // One for nil->value case, one for value->value case
		t.Errorf("Expected exactly 2 gorm.Expr calls in root JSONB field, got %d", gormExprCount)
	}
}

// Helper function to extract specific field comparison code from generated diff function
func extractFunctionCode(code, structName, fieldName string) string {
	lines := strings.Split(code, "\n")
	inFunction := false
	inFieldComparison := false
	var result []string

	functionStart := "func (new *" + structName + ") Diff("
	fieldComment := "// Compare " + fieldName

	for _, line := range lines {
		if strings.Contains(line, functionStart) {
			inFunction = true
			continue
		}

		if inFunction {
			if strings.Contains(line, fieldComment) {
				inFieldComparison = true
			}

			if inFieldComparison {
				result = append(result, line)

				// Stop when we hit the next field comparison or end of function
				if strings.HasPrefix(strings.TrimSpace(line), "// Compare ") && !strings.Contains(line, fieldComment) {
					break
				}
				if strings.Contains(line, "return diff") {
					break
				}
			}
		}
	}

	return strings.Join(result, "\n")
}
