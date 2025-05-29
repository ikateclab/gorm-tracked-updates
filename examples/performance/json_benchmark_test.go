package performance

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/bytedance/sonic"
)

// Test data structures for JSON benchmarking
type BenchmarkData struct {
	ID       string                 `json:"id"`
	Name     string                 `json:"name"`
	Age      int                    `json:"age"`
	Email    string                 `json:"email"`
	Active   bool                   `json:"active"`
	Tags     []string               `json:"tags"`
	Metadata map[string]interface{} `json:"metadata"`
	Address  BenchmarkAddress       `json:"address"`
	Contacts []BenchmarkContact     `json:"contacts"`
}

type BenchmarkAddress struct {
	Street   string `json:"street"`
	City     string `json:"city"`
	State    string `json:"state"`
	ZipCode  string `json:"zip_code"`
	Country  string `json:"country"`
	Location struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
}

type BenchmarkContact struct {
	Type     string `json:"type"`
	Value    string `json:"value"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

// Create complex test data
func createBenchmarkData() BenchmarkData {
	return BenchmarkData{
		ID:     "user-12345",
		Name:   "John Doe",
		Age:    30,
		Email:  "john.doe@example.com",
		Active: true,
		Tags:   []string{"developer", "golang", "backend", "senior", "team-lead"},
		Metadata: map[string]interface{}{
			"role":        "Senior Developer",
			"department":  "Engineering",
			"team":        "Backend",
			"level":       "L5",
			"salary":      120000,
			"remote":      true,
			"skills":      []string{"Go", "Python", "Docker", "Kubernetes"},
			"performance": map[string]interface{}{
				"rating": 4.8,
				"goals":  []string{"improve performance", "mentor juniors"},
			},
		},
		Address: BenchmarkAddress{
			Street:  "123 Main Street, Apt 4B",
			City:    "San Francisco",
			State:   "CA",
			ZipCode: "94105",
			Country: "USA",
			Location: struct {
				Latitude  float64 `json:"latitude"`
				Longitude float64 `json:"longitude"`
			}{
				Latitude:  37.7749,
				Longitude: -122.4194,
			},
		},
		Contacts: []BenchmarkContact{
			{Type: "email", Value: "john.doe@example.com", Primary: true, Verified: true},
			{Type: "phone", Value: "+1-555-123-4567", Primary: true, Verified: true},
			{Type: "mobile", Value: "+1-555-987-6543", Primary: false, Verified: true},
			{Type: "slack", Value: "@johndoe", Primary: false, Verified: false},
			{Type: "linkedin", Value: "linkedin.com/in/johndoe", Primary: false, Verified: true},
		},
	}
}

// Benchmark Marshal operations
func BenchmarkMarshalNativeJSON(b *testing.B) {
	data := createBenchmarkData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMarshalSonic(b *testing.B) {
	data := createBenchmarkData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := sonic.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Unmarshal operations
func BenchmarkUnmarshalNativeJSON(b *testing.B) {
	data := createBenchmarkData()
	jsonData, _ := json.Marshal(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result BenchmarkData
		err := json.Unmarshal(jsonData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkUnmarshalSonic(b *testing.B) {
	data := createBenchmarkData()
	jsonData, _ := sonic.Marshal(data)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		var result BenchmarkData
		err := sonic.Unmarshal(jsonData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark Marshal + Unmarshal round trip
func BenchmarkRoundTripNativeJSON(b *testing.B) {
	data := createBenchmarkData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsonData, err := json.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
		var result BenchmarkData
		err = json.Unmarshal(jsonData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRoundTripSonic(b *testing.B) {
	data := createBenchmarkData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		jsonData, err := sonic.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
		var result BenchmarkData
		err = sonic.Unmarshal(jsonData, &result)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Benchmark small data structures (like what might be in diff operations)
func BenchmarkSmallDataMarshalNativeJSON(b *testing.B) {
	data := BenchmarkContact{
		Type:     "email",
		Value:    "test@example.com",
		Primary:  true,
		Verified: true,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := json.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkSmallDataMarshalSonic(b *testing.B) {
	data := BenchmarkContact{
		Type:     "email",
		Value:    "test@example.com",
		Primary:  true,
		Verified: true,
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := sonic.Marshal(data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Realistic benchmark simulating your actual usage pattern
func BenchmarkRealisticUsageNativeJSON(b *testing.B) {
	// Simulate JsonbStringSlice database operations (unmarshal heavy)
	// Plus some diff operations (marshal)
	slice := []string{"tag1", "tag2", "tag3", "tag4", "tag5"}
	jsonData, _ := json.Marshal(slice)

	contact := BenchmarkContact{Type: "email", Value: "test@example.com", Primary: true}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate database read (unmarshal) - happens frequently
		var readSlice []string
		json.Unmarshal(jsonData, &readSlice)

		// Simulate diff operation (marshal) - happens less frequently
		if i%3 == 0 { // Every 3rd iteration
			json.Marshal(contact)
		}
	}
}

func BenchmarkRealisticUsageSonic(b *testing.B) {
	// Simulate JsonbStringSlice database operations (unmarshal heavy)
	// Plus some diff operations (marshal)
	slice := []string{"tag1", "tag2", "tag3", "tag4", "tag5"}
	jsonData, _ := sonic.Marshal(slice)

	contact := BenchmarkContact{Type: "email", Value: "test@example.com", Primary: true}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate database read (unmarshal) - happens frequently
		var readSlice []string
		sonic.Unmarshal(jsonData, &readSlice)

		// Simulate diff operation (marshal) - happens less frequently
		if i%3 == 0 { // Every 3rd iteration
			sonic.Marshal(contact)
		}
	}
}

// Hybrid approach: Native JSON for marshal, Sonic for unmarshal
func BenchmarkRealisticUsageHybrid(b *testing.B) {
	// Simulate JsonbStringSlice database operations (unmarshal heavy)
	// Plus some diff operations (marshal)
	slice := []string{"tag1", "tag2", "tag3", "tag4", "tag5"}
	jsonData, _ := json.Marshal(slice) // Use native for marshal

	contact := BenchmarkContact{Type: "email", Value: "test@example.com", Primary: true}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Simulate database read (unmarshal) - happens frequently
		var readSlice []string
		sonic.Unmarshal(jsonData, &readSlice) // Use Sonic for unmarshal

		// Simulate diff operation (marshal) - happens less frequently
		if i%3 == 0 { // Every 3rd iteration
			json.Marshal(contact) // Use native for marshal
		}
	}
}

// Benchmark pure marshal operations comparison
func BenchmarkMarshalComparison(b *testing.B) {
	contact := BenchmarkContact{Type: "email", Value: "test@example.com", Primary: true}

	b.Run("Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(contact)
		}
	})

	b.Run("Sonic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sonic.Marshal(contact)
		}
	})
}

// Benchmark pure unmarshal operations comparison
func BenchmarkUnmarshalComparison(b *testing.B) {
	contact := BenchmarkContact{Type: "email", Value: "test@example.com", Primary: true}
	jsonData, _ := json.Marshal(contact)

	b.Run("Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var result BenchmarkContact
			json.Unmarshal(jsonData, &result)
		}
	})

	b.Run("Sonic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			var result BenchmarkContact
			sonic.Unmarshal(jsonData, &result)
		}
	})
}

// More detailed marshal benchmarks with different data sizes
func BenchmarkDetailedMarshalComparison(b *testing.B) {
	// Very small data
	smallData := struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}{ID: 1, Name: "test"}

	// Medium data
	mediumData := BenchmarkContact{
		Type:     "email",
		Value:    "test@example.com",
		Primary:  true,
		Verified: true,
	}

	// Large data
	largeData := createBenchmarkData()

	b.Run("Small/Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(smallData)
		}
	})

	b.Run("Small/Sonic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sonic.Marshal(smallData)
		}
	})

	b.Run("Medium/Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(mediumData)
		}
	})

	b.Run("Medium/Sonic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sonic.Marshal(mediumData)
		}
	})

	b.Run("Large/Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(largeData)
		}
	})

	b.Run("Large/Sonic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sonic.Marshal(largeData)
		}
	})
}

// Test with different Go versions and build flags
func BenchmarkMarshalWithBuildTags(b *testing.B) {
	data := BenchmarkContact{
		Type:     "email",
		Value:    "test@example.com",
		Primary:  true,
		Verified: true,
	}

	b.Run("Native_Optimized", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("Sonic_Optimized", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_, err := sonic.Marshal(data)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// Test string-heavy vs number-heavy data
func BenchmarkMarshalDataTypes(b *testing.B) {
	stringHeavy := struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Tags        string `json:"tags"`
		Notes       string `json:"notes"`
	}{
		Name:        "This is a long string with lots of text content",
		Description: "Another very long string with detailed description content",
		Category:    "Category with more text",
		Tags:        "tag1,tag2,tag3,tag4,tag5,tag6,tag7,tag8,tag9,tag10",
		Notes:       "Additional notes with even more text content here",
	}

	numberHeavy := struct {
		ID       int     `json:"id"`
		Count    int     `json:"count"`
		Price    float64 `json:"price"`
		Discount float64 `json:"discount"`
		Tax      float64 `json:"tax"`
		Total    float64 `json:"total"`
		Quantity int     `json:"quantity"`
		Weight   float64 `json:"weight"`
	}{
		ID: 12345, Count: 100, Price: 99.99, Discount: 10.0,
		Tax: 8.25, Total: 98.24, Quantity: 5, Weight: 2.5,
	}

	b.Run("StringHeavy/Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(stringHeavy)
		}
	})

	b.Run("StringHeavy/Sonic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sonic.Marshal(stringHeavy)
		}
	})

	b.Run("NumberHeavy/Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			json.Marshal(numberHeavy)
		}
	})

	b.Run("NumberHeavy/Sonic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			sonic.Marshal(numberHeavy)
		}
	})
}

// Test with larger, more complex data like Sonic's official benchmarks
func BenchmarkOfficialSonicStyle(b *testing.B) {
	// Create data similar to Sonic's official benchmarks
	largeComplexData := struct {
		Users []struct {
			ID       int                    `json:"id"`
			Name     string                 `json:"name"`
			Email    string                 `json:"email"`
			Active   bool                   `json:"active"`
			Profile  map[string]interface{} `json:"profile"`
			Settings struct {
				Theme       string   `json:"theme"`
				Language    string   `json:"language"`
				Permissions []string `json:"permissions"`
			} `json:"settings"`
		} `json:"users"`
		Metadata map[string]interface{} `json:"metadata"`
		Config   struct {
			Version     string  `json:"version"`
			Environment string  `json:"environment"`
			Debug       bool    `json:"debug"`
			Timeout     float64 `json:"timeout"`
		} `json:"config"`
	}{
		Users: make([]struct {
			ID       int                    `json:"id"`
			Name     string                 `json:"name"`
			Email    string                 `json:"email"`
			Active   bool                   `json:"active"`
			Profile  map[string]interface{} `json:"profile"`
			Settings struct {
				Theme       string   `json:"theme"`
				Language    string   `json:"language"`
				Permissions []string `json:"permissions"`
			} `json:"settings"`
		}, 100),
		Metadata: map[string]interface{}{
			"total_users": 100,
			"created_at":  "2024-01-01T00:00:00Z",
			"version":     "1.0.0",
			"features":    []string{"auth", "api", "dashboard"},
		},
		Config: struct {
			Version     string  `json:"version"`
			Environment string  `json:"environment"`
			Debug       bool    `json:"debug"`
			Timeout     float64 `json:"timeout"`
		}{
			Version:     "2.1.0",
			Environment: "production",
			Debug:       false,
			Timeout:     30.5,
		},
	}

	// Fill users array
	for i := 0; i < 100; i++ {
		largeComplexData.Users[i] = struct {
			ID       int                    `json:"id"`
			Name     string                 `json:"name"`
			Email    string                 `json:"email"`
			Active   bool                   `json:"active"`
			Profile  map[string]interface{} `json:"profile"`
			Settings struct {
				Theme       string   `json:"theme"`
				Language    string   `json:"language"`
				Permissions []string `json:"permissions"`
			} `json:"settings"`
		}{
			ID:     i + 1,
			Name:   fmt.Sprintf("User %d", i+1),
			Email:  fmt.Sprintf("user%d@example.com", i+1),
			Active: i%2 == 0,
			Profile: map[string]interface{}{
				"age":      20 + (i % 50),
				"country":  "US",
				"timezone": "UTC",
				"avatar":   fmt.Sprintf("https://example.com/avatar/%d.jpg", i+1),
			},
			Settings: struct {
				Theme       string   `json:"theme"`
				Language    string   `json:"language"`
				Permissions []string `json:"permissions"`
			}{
				Theme:       "dark",
				Language:    "en",
				Permissions: []string{"read", "write", "admin"},
			},
		}
	}

	b.Run("LargeComplex/Native", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(largeComplexData)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("LargeComplex/Sonic", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, err := sonic.Marshal(largeComplexData)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}