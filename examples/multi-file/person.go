package multifile

// Person represents a person with various field types
type Person struct {
	Name     string                 // Simple type
	Age      int                    // Simple type
	Address  Address                // Nested struct from address.go
	Contacts []Contact              // Slice of nested structs from contact.go
	Manager  *Person                // Pointer to the same struct type
	Metadata map[string]interface{} // Map type
}
