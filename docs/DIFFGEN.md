# DiffGen - Efficient Diff Function Generator

DiffGen automatically generates efficient diff functions for Go structs that detect and return only the changed fields between two instances.

## Overview

DiffGen creates functions that compare two struct instances and return a map containing only the fields that have changed, with their new values. This is perfect for GORM's selective update functionality.

## Features

- **Selective Updates**: Only changed fields are included in the diff
- **Nested Struct Support**: Recursive diffing for complex structures
- **Type Safety**: Generated functions are fully type-safe
- **Performance**: No reflection overhead in generated code
- **GORM Integration**: Perfect for `db.Model().Updates(diff)` operations

## Usage

### Basic Usage

```go
package main

import "github.com/ikateclab/gorm-tracked-updates/pkg/diffgen"

func main() {
    // Create generator
    generator := diffgen.New()

    // Parse struct definitions
    err := generator.ParseFile("structs.go")
    if err != nil {
        panic(err)
    }

    // Generate diff functions
    code, err := generator.GenerateCode()
    if err != nil {
        panic(err)
    }

    // Write to file
    err = generator.WriteToFile("generated_diff.go")
    if err != nil {
        panic(err)
    }
}
```

### Generated Functions

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

DiffGen generates:

```go
func DiffPerson(a, b Person) map[string]interface{} {
    diff := make(map[string]interface{})

    // Simple type comparison
    if a.Name != b.Name {
        diff["Name"] = b.Name
    }
    if a.Age != b.Age {
        diff["Age"] = b.Age
    }

    // Struct type comparison
    if !reflect.DeepEqual(a.Address, b.Address) {
        nestedDiff := DiffAddress(a.Address, b.Address)
        if len(nestedDiff) > 0 {
            diff["Address"] = nestedDiff
        }
    }

    // Slice comparison
    if !reflect.DeepEqual(a.Contacts, b.Contacts) {
        diff["Contacts"] = b.Contacts
    }

    // Pointer comparison
    if !reflect.DeepEqual(a.Manager, b.Manager) {
        if a.Manager == nil || b.Manager == nil {
            diff["Manager"] = b.Manager
        } else {
            nestedDiff := DiffPerson(*a.Manager, *b.Manager)
            if len(nestedDiff) > 0 {
                diff["Manager"] = nestedDiff
            }
        }
    }

    // Map comparison
    if !reflect.DeepEqual(a.Metadata, b.Metadata) {
        diff["Metadata"] = b.Metadata
    }

    return diff
}
```

## Field Type Handling

### Simple Types
- **Types**: `string`, `int`, `bool`, `float64`, etc.
- **Strategy**: Direct comparison with `!=`
- **Performance**: Fastest possible comparison

### Struct Types
- **Types**: Custom struct fields
- **Strategy**: Recursive diffing using generated functions
- **Output**: Nested diff maps for changed fields only

### Pointer Types
- **Types**: `*Person`, `*Address`, etc.
- **Strategy**: Nil-safe comparison with recursive diffing
- **Handling**: Proper nil pointer management

### Slice Types
- **Types**: `[]Contact`, `[]*Person`, etc.
- **Strategy**: Deep equality check, full replacement on change
- **Note**: Element-by-element diffing not implemented (complex)

### Map Types
- **Types**: `map[string]interface{}`, etc.
- **Strategy**: Deep equality check, full replacement on change
- **Performance**: Efficient for most use cases

### Interface Types
- **Types**: `interface{}`, custom interfaces
- **Strategy**: Deep equality check with reflection
- **Safety**: Handles unknown types safely

## GORM Integration

Perfect for selective database updates:

```go
// Before modification
original := user

// After modification
user.Name = "New Name"
user.Email = "new@example.com"

// Generate diff
diff := DiffUser(original, user)
// Result: {"Name": "New Name", "Email": "new@example.com"}

// GORM selective update
db.Model(&user).Updates(diff)
// SQL: UPDATE users SET name = 'New Name', email = 'new@example.com' WHERE id = ?
```

## Advanced Examples

### Nested Struct Changes

```go
person1 := Person{
    Name: "John",
    Address: Address{City: "NYC", State: "NY"},
}

person2 := Person{
    Name: "John",
    Address: Address{City: "LA", State: "CA"},
}

diff := DiffPerson(person1, person2)
// Result: {
//   "Address": {
//     "City": "LA",
//     "State": "CA"
//   }
// }
```

### Pointer Changes

```go
person1 := Person{
    Name: "John",
    Manager: &Person{Name: "Jane", Age: 45},
}

person2 := Person{
    Name: "John",
    Manager: &Person{Name: "Jane", Age: 46},
}

diff := DiffPerson(person1, person2)
// Result: {
//   "Manager": {
//     "Age": 46
//   }
// }
```

## Performance Characteristics

- **Simple Fields**: O(1) direct comparison
- **Nested Structs**: O(n) where n is number of fields
- **Slices/Maps**: O(n) deep equality check
- **Memory**: Minimal allocations, only for changed fields

## Best Practices

1. **Use with Cloning**: Combine with CloneGen for safe workflows
2. **Handle Nil Pointers**: Generated code handles nil safely
3. **Batch Updates**: Collect multiple diffs for batch operations
4. **Validation**: Validate diff results before database operations

## Limitations

1. **Slice Diffing**: No element-by-element diffing (full replacement)
2. **Map Diffing**: No key-by-key diffing (full replacement)
3. **Circular References**: Not handled (would cause infinite recursion)
4. **Private Fields**: Only exported fields are processed

## Testing

Generated functions can be tested like any Go code:

```go
func TestDiffPerson(t *testing.T) {
    person1 := Person{Name: "John", Age: 30}
    person2 := Person{Name: "John", Age: 31}

    diff := DiffPerson(person1, person2)

    if diff["Age"] != 31 {
        t.Errorf("Expected Age diff to be 31")
    }
    if _, exists := diff["Name"]; exists {
        t.Errorf("Name should not be in diff (unchanged)")
    }
}
```

## Error Handling

The generator handles various edge cases:

- **Empty Structs**: Generates valid functions that return empty maps
- **No Changes**: Returns empty map when structs are identical
- **Nil Pointers**: Safe comparison without panics
- **Type Mismatches**: Compile-time safety prevents runtime errors

## Integration Examples

### With GORM Hooks

```go
func (u *User) BeforeUpdate(tx *gorm.DB) error {
    if original, ok := tx.Statement.Context.Value("original").(*User); ok {
        diff := DiffUser(*original, *u)
        // Log changes, validate, etc.
        log.Printf("User changes: %+v", diff)
    }
    return nil
}
```

### With Audit Logging

```go
func AuditChanges(original, updated User) {
    diff := DiffUser(original, updated)
    for field, newValue := range diff {
        auditLog := AuditLog{
            Field:    field,
            OldValue: getFieldValue(original, field),
            NewValue: newValue,
        }
        db.Create(&auditLog)
    }
}
```

DiffGen provides a robust, efficient solution for change detection in Go applications, particularly when working with GORM and database operations.
