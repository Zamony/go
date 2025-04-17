
[![Go Reference](https://pkg.go.dev/badge/github.com/Zamony/go/validate.svg)](https://pkg.go.dev/github.com/Zamony/go/validate)

A lightweight, type-safe validation library for Go with zero dependencies.

## Features

- Type-safe validation functions using generics
- Composable validator chains
- Built-in common validators:
  - String length (min/max)
  - Numeric ranges (min/max)
  - Regular expressions
  - Enums (one-of)
  - Slice validations (length, uniqueness)
- Custom error messages with field names
- Two validation modes:
  - `All()` - collect all validation errors
  - `Chain()` - fail on first error

## Installation

```bash
go get github.com/Zamony/go/validate
```

## Usage

```go
import "github.com/Zamony/go/validate"

type User struct {
    Username string
    Age      int
    Email    string
    Tags     []string
}

func ValidateUser(u User) error {
    emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
    
    return validate.All(
        validate.Value("username", u.Username,
            validate.MinLength(3),
            validate.MaxLength(20),
        ),
        validate.Value("age", u.Age,
            validate.Min(18),
            validate.Max(120),
        ),
        validate.Value("email", u.Email,
            validate.MinLength(5),
            validate.RegExp(emailRegex),
        ),
        validate.Value("tags", u.Tags,
            validate.MaxItems(5),
            validate.UniqueItems[string](),
        ),
    )
}
```

## Custom Validators

Create your own validators by implementing the Func[T] signature:

```go
func IsUpperCase(name string, value string) error {
    if value != strings.ToUpper(value) {
        return validate.ValidationError{
            Field: name,
            Msg:   "must be uppercase",
        }
    }
    return nil
}
```