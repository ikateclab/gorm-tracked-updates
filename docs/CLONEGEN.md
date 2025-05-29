# CloneGen - Performant Clone Method Generator

CloneGen automatically generates high-performance deep clone methods for Go structs, ensuring complete memory independence between original and cloned instances.

## Overview

CloneGen creates methods that perform deep copying of struct instances, handling all field types optimally without reflection overhead. Perfect for creating backups before modifications or ensuring data isolation.

## Features

- **Deep Cloning**: Complete independence between original and clone
- **Type Safety**: Generated methods are fully type-safe
- **Performance**: 3-23x faster than reflection/JSON alternatives
- **Memory Safety**: Proper handling of pointers, slices, and maps
- **Nil Safety**: Correct handling of nil pointers and empty collections

## Usage

### Basic Usage

```go
package main

import "github.com/ikateclab/gorm-tracked-updates/pkg/clonegen"

func main() {
    // Create generator
    generator := clonegen.New()

    // Parse struct definitions
    err := generator.ParseFile("structs.go")
    if err != nil {
        panic(err)
    }

    // Generate clone methods
    code, err := generator.GenerateCode()
    if err != nil {
        panic(err)
    }

    // Write to file
    err = generator.WriteToFile("generated_clone.go")
    if err != nil {
        panic(err)
    }
}
```

### Generated Methods

For a struct like:

```go
type Person struct {
    Name     string
    Age      int
    Address  Address
    Contacts []Contact
    Manager  *Person
    Metadata map[string]interface{}
}
```

CloneGen generates:

```go
func (original Person) ClonePerson() Person {
    clone := Person{}

    // Simple type - direct assignment
    clone.Name = original.Name
    clone.Age = original.Age

    // Struct type - recursive clone
    clone.Address = original.Address.CloneAddress()

    // Slice - create new slice and clone elements
    if original.Contacts != nil {
        clone.Contacts = make([]Contact, len(original.Contacts))
        for i, item := range original.Contacts {
            clone.Contacts[i] = item.CloneContact()
        }
    }

    // Pointer to struct - create new instance and clone
    if original.Manager != nil {
        clonedManager := original.Manager.ClonePerson()
        clone.Manager = &clonedManager
    }

    // Map - create new map and copy key-value pairs
    if original.Metadata != nil {
        clone.Metadata = make(map[string]interface{})
        for k, v := range original.Metadata {
            clone.Metadata[k] = v
        }
    }

    return clone
}
```

## Field Type Handling

### Simple Types
- **Types**: `string`, `int`, `bool`, `float64`, etc.
- **Strategy**: Direct assignment (copy by value)
- **Performance**: Fastest possible copying

### Struct Types
- **Types**: Custom struct fields
- **Strategy**: Recursive cloning using generated methods
- **Independence**: Complete memory separation

### Pointer Types
- **Types**: `*Person`, `*Address`, etc.
- **Strategy**: Create new instance, clone pointed-to value
- **Safety**: Handles nil pointers correctly

### Slice Types
- **Types**: `[]Contact`, `[]*Person`, etc.
- **Strategy**: Create new slice, clone each element
- **Optimization**: Different strategies for struct vs primitive elements

### Map Types
- **Types**: `map[string]interface{}`, etc.
- **Strategy**: Create new map, copy key-value pairs
- **Note**: Values copied by reference for complex types

### Interface Types
- **Types**: `interface{}`, custom interfaces
- **Strategy**: Reflection-based copying for safety
- **Performance**: Slower but safe for unknown types

## Performance Comparison

Benchmark results (10,000 iterations):

| Method | Time | Relative Performance |
|--------|------|---------------------|
| Generated Clone | 323.7 ns/op | 1.0x (baseline) |
| Reflection Clone | 1207 ns/op | 3.7x slower |
| JSON Clone | 7450 ns/op | 23x slower |

## Usage Examples

### Basic Cloning

```go
original := Person{
    Name: "John Doe",
    Age:  30,
}

cloned := original.ClonePerson()
cloned.Age = 31

// original.Age is still 30
// cloned.Age is now 31
```

### Deep Independence

```go
original := Person{
    Name: "John",
    Manager: &Person{Name: "Jane", Age: 45},
}

cloned := original.ClonePerson()
cloned.Manager.Age = 46

// original.Manager.Age is still 45
// cloned.Manager.Age is now 46
// Different memory addresses: original.Manager != cloned.Manager
```

### Slice Independence

```go
original := Person{
    Contacts: []Contact{
        {Type: "email", Value: "john@example.com"},
    },
}

cloned := original.ClonePerson()
cloned.Contacts[0].Value = "john@newexample.com"

// original.Contacts[0].Value is still "john@example.com"
// cloned.Contacts[0].Value is now "john@newexample.com"
```

## Integration with DiffGen

Perfect workflow for tracked updates:

```go
// 1. Clone before modifications
backup := user.ClonePerson()

// 2. Make changes
user.Name = "New Name"
user.Email = "new@example.com"

// 3. Generate diff
changes := DiffPerson(backup, user)

// 4. GORM selective update
db.Model(&user).Updates(changes)
```

## Memory Safety

CloneGen ensures complete memory independence:

```go
original := Person{
    Metadata: map[string]interface{}{
        "settings": map[string]string{"theme": "dark"},
    },
}

cloned := original.ClonePerson()

// Maps are independent
cloned.Metadata["new_key"] = "new_value"
// original.Metadata does not contain "new_key"

// But nested values may share references (for interface{} types)
// This is a limitation of the current implementation
```

## Performance Characteristics

- **Simple Fields**: O(1) direct assignment
- **Nested Structs**: O(n) where n is total number of fields
- **Slices**: O(n) where n is slice length
- **Maps**: O(n) where n is map size
- **Memory**: Allocates only necessary memory for new instances

## Best Practices

1. **Use Before Modifications**: Clone before making changes
2. **Combine with DiffGen**: Use together for optimal workflows
3. **Handle Large Structures**: Be aware of memory usage for large structs
4. **Test Independence**: Verify clones are truly independent

## Limitations

1. **Interface Values**: May share references for complex interface{} values
2. **Circular References**: Not handled (would cause infinite recursion)
3. **Private Fields**: Only exported fields are cloned
4. **Function Fields**: Function values are copied by reference

## Testing Clone Methods

```go
func TestClonePerson(t *testing.T) {
    original := Person{
        Name: "John",
        Age:  30,
        Manager: &Person{Name: "Jane", Age: 45},
    }

    cloned := original.ClonePerson()

    // Test equality
    if !reflect.DeepEqual(original, cloned) {
        t.Error("Clone should be equal to original")
    }

    // Test independence
    cloned.Manager.Age = 46
    if original.Manager.Age == cloned.Manager.Age {
        t.Error("Clone should be independent of original")
    }

    // Test different references
    if original.Manager == cloned.Manager {
        t.Error("Pointers should be different")
    }
}
```

## Error Handling

Generated methods handle edge cases safely:

- **Nil Pointers**: No panic, nil values preserved
- **Empty Slices**: Handled correctly (nil vs empty)
- **Empty Maps**: Proper initialization
- **Zero Values**: Correctly copied

## Advanced Use Cases

### State Management

```go
type GameState struct {
    Players []Player
    Board   Board
    Turn    int
}

// Save state before move
savedState := gameState.CloneGameState()

// Make move
gameState.MakeMove(move)

// Rollback if invalid
if !gameState.IsValid() {
    gameState = savedState
}
```

### Testing

```go
func TestUserUpdate(t *testing.T) {
    // Create test data
    user := createTestUser()
    original := user.CloneUser()

    // Perform operation
    updateUser(&user, request)

    // Verify changes
    diff := DiffUser(original, user)
    assert.Contains(t, diff, "Name")
    assert.NotContains(t, diff, "ID")
}
```

### Caching

```go
type Cache struct {
    data map[string]Person
}

func (c *Cache) Get(key string) Person {
    if person, exists := c.data[key]; exists {
        return person.ClonePerson() // Return independent copy
    }
    return Person{}
}
```

CloneGen provides a robust, high-performance solution for deep copying in Go applications, ensuring memory safety and independence while maintaining excellent performance characteristics.
