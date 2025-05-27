package multifile

// Company represents a company with employees
type Company struct {
	Name      string
	Address   Address   // From address.go
	Employees []Person  // From person.go
	CEO       *Person   // From person.go
	Founded   int
	Active    bool
}

// Project represents a project with team members
type Project struct {
	Name        string
	Description string
	TeamLead    *Person           // From person.go
	Members     []*Person         // From person.go
	Company     *Company          // Self-reference
	Budget      float64
	Tags        []string
	Properties  map[string]string
}
