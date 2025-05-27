package structs

// Example nested structs for demonstrating diff and clone generation
// These structs showcase various field types and relationships

// Address represents a physical address
type Address struct {
	Street  string
	City    string
	State   string
	ZipCode string
	Country string
}

// Contact represents a contact method
type Contact struct {
	Type  string // email, phone, etc.
	Value string
}

// Person represents a person with various field types
type Person struct {
	Name     string                 // Simple type
	Age      int                    // Simple type
	Address  Address                // Nested struct
	Contacts []Contact              // Slice of nested structs
	Manager  *Person                // Pointer to the same struct type
	Metadata map[string]interface{} // Map type
}

// Company represents a company with employees
type Company struct {
	Name      string
	Address   Address
	Employees []Person
	CEO       *Person
	Founded   int
	Active    bool
}

// Project represents a project with team members
type Project struct {
	Name        string
	Description string
	TeamLead    *Person
	Members     []*Person
	Company     *Company
	Budget      float64
	Tags        []string
	Properties  map[string]string
}
