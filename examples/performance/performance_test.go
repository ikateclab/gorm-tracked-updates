package main

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

// Performance test structures
type PerfAddress struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
}

type PerfContact struct {
	Type  string
	Value string
}

type PerfPerson struct {
	Name     string
	Age      int
	Address  PerfAddress
	Contacts []PerfContact
	Manager  *PerfPerson
	Metadata map[string]interface{}
}

// Generated clone method (simulated)
func (original PerfPerson) Clone() PerfPerson {
	clone := PerfPerson{}

	// Simple types
	clone.Name = original.Name
	clone.Age = original.Age

	// Struct type
	clone.Address = original.Address.Clone()

	// Slice
	if original.Contacts != nil {
		clone.Contacts = make([]PerfContact, len(original.Contacts))
		for i, item := range original.Contacts {
			clone.Contacts[i] = item.Clone()
		}
	}

	// Pointer
	if original.Manager != nil {
		clonedManager := original.Manager.Clone()
		clone.Manager = &clonedManager
	}

	// Map
	if original.Metadata != nil {
		clone.Metadata = make(map[string]interface{})
		for k, v := range original.Metadata {
			clone.Metadata[k] = v
		}
	}

	return clone
}

func (original PerfAddress) Clone() PerfAddress {
	return PerfAddress{
		Street:  original.Street,
		City:    original.City,
		State:   original.State,
		ZipCode: original.ZipCode,
		Country: original.Country,
	}
}

func (original PerfContact) Clone() PerfContact {
	return PerfContact{
		Type:  original.Type,
		Value: original.Value,
	}
}

// Reflection-based clone
func cloneWithReflection(original interface{}) interface{} {
	val := reflect.ValueOf(original)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	clone := reflect.New(val.Type()).Elem()
	copyWithReflection(clone, val)
	return clone.Interface()
}

func copyWithReflection(dst, src reflect.Value) {
	switch src.Kind() {
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			srcField := src.Field(i)
			dstField := dst.Field(i)
			if dstField.CanSet() {
				copyWithReflection(dstField, srcField)
			}
		}
	case reflect.Slice:
		if !src.IsNil() {
			newSlice := reflect.MakeSlice(src.Type(), src.Len(), src.Cap())
			for i := 0; i < src.Len(); i++ {
				copyWithReflection(newSlice.Index(i), src.Index(i))
			}
			dst.Set(newSlice)
		}
	case reflect.Map:
		if !src.IsNil() {
			newMap := reflect.MakeMap(src.Type())
			for _, key := range src.MapKeys() {
				newMap.SetMapIndex(key, src.MapIndex(key))
			}
			dst.Set(newMap)
		}
	case reflect.Ptr:
		if !src.IsNil() {
			newPtr := reflect.New(src.Type().Elem())
			copyWithReflection(newPtr.Elem(), src.Elem())
			dst.Set(newPtr)
		}
	default:
		dst.Set(src)
	}
}

// JSON-based clone
func cloneWithJSON(original interface{}) interface{} {
	data, _ := json.Marshal(original)
	clone := reflect.New(reflect.TypeOf(original)).Interface()
	json.Unmarshal(data, clone)
	return reflect.ValueOf(clone).Elem().Interface()
}

// Create test data
func createTestPerson() PerfPerson {
	return PerfPerson{
		Name: "John Doe",
		Age:  30,
		Address: PerfAddress{
			Street:  "123 Main St",
			City:    "Anytown",
			State:   "CA",
			ZipCode: "12345",
			Country: "USA",
		},
		Contacts: []PerfContact{
			{Type: "email", Value: "john@example.com"},
			{Type: "phone", Value: "555-1234"},
			{Type: "mobile", Value: "555-5678"},
		},
		Manager: &PerfPerson{
			Name: "Jane Doe",
			Age:  45,
			Address: PerfAddress{
				Street:  "789 Oak Dr",
				City:    "Managertown",
				State:   "CA",
				ZipCode: "54321",
				Country: "USA",
			},
			Contacts: []PerfContact{
				{Type: "email", Value: "jane@company.com"},
			},
			Metadata: map[string]interface{}{
				"role":       "Senior Manager",
				"department": "Engineering",
			},
		},
		Metadata: map[string]interface{}{
			"role":     "developer",
			"team":     "backend",
			"level":    "senior",
			"projects": []string{"project1", "project2"},
		},
	}
}

// Benchmark tests
func BenchmarkCloneGenerated(b *testing.B) {
	person := createTestPerson()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = person.Clone()
	}
}

func BenchmarkCloneReflection(b *testing.B) {
	person := createTestPerson()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = cloneWithReflection(person).(PerfPerson)
	}
}

func BenchmarkCloneJSON(b *testing.B) {
	person := createTestPerson()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = cloneWithJSON(person).(PerfPerson)
	}
}

// Correctness tests
func TestCloneCorrectness(t *testing.T) {
	original := createTestPerson()

	// Test generated clone
	t.Run("Generated clone correctness", func(t *testing.T) {
		cloned := original.Clone()

		// Verify equality
		if !reflect.DeepEqual(original, cloned) {
			t.Error("Generated clone should be equal to original")
		}

		// Verify independence
		cloned.Age = 31
		cloned.Address.City = "Newtown"
		cloned.Manager.Age = 46

		if original.Age == cloned.Age {
			t.Error("Simple field modification should not affect original")
		}
		if original.Address.City == cloned.Address.City {
			t.Error("Nested struct modification should not affect original")
		}
		if original.Manager.Age == cloned.Manager.Age {
			t.Error("Pointer target modification should not affect original")
		}
	})

	// Test reflection clone
	t.Run("Reflection clone correctness", func(t *testing.T) {
		cloned := cloneWithReflection(original).(PerfPerson)

		// Verify equality
		if !reflect.DeepEqual(original, cloned) {
			t.Error("Reflection clone should be equal to original")
		}
	})

	// Test JSON clone
	t.Run("JSON clone correctness", func(t *testing.T) {
		cloned := cloneWithJSON(original).(PerfPerson)

		// Verify equality (note: JSON clone may have type differences)
		if cloned.Name != original.Name || cloned.Age != original.Age {
			t.Error("JSON clone should preserve basic fields")
		}
	})
}

// Performance comparison test
func TestPerformanceComparison(t *testing.T) {
	person := createTestPerson()
	iterations := 10000

	// Generated clone
	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = person.Clone()
	}
	generatedTime := time.Since(start)

	// Reflection clone
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = cloneWithReflection(person)
	}
	reflectionTime := time.Since(start)

	// JSON clone
	start = time.Now()
	for i := 0; i < iterations; i++ {
		_ = cloneWithJSON(person)
	}
	jsonTime := time.Since(start)

	t.Logf("Performance comparison (%d iterations):", iterations)
	t.Logf("Generated clone: %v", generatedTime)
	t.Logf("Reflection clone: %v (%.1fx slower)", reflectionTime, float64(reflectionTime)/float64(generatedTime))
	t.Logf("JSON clone: %v (%.1fx slower)", jsonTime, float64(jsonTime)/float64(generatedTime))

	// Generated should be fastest
	if generatedTime > reflectionTime {
		t.Log("Warning: Generated clone is slower than reflection (unexpected)")
	}
	if generatedTime > jsonTime {
		t.Log("Warning: Generated clone is slower than JSON (unexpected)")
	}
}
