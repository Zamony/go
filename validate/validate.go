// Package validate provides a flexible validation framework for Go.
// It allows defining validation rules for various types and composing them together.
package validate

import (
	"cmp"
	"errors"
	"fmt"
	"regexp"
	"unicode/utf8"
)

// ValidationError represents a validation failure for a specific field.
// It implements the error interface.
type ValidationError struct {
	Field string // The name of the field that failed validation
	Msg   string // The validation error message
}

// Error implements the error interface for ValidationError.
func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Msg)
}

// Func is a generic validation function type that validates a value of type T.
// It takes the field name and value as parameters and returns an error if validation fails.
type Func[T any] func(name string, value T) error

// Validator is a function type that performs validation and returns an error if it fails.
type Validator func() error

// Value creates a Validator for a specific value with the given validation functions.
// name is the field name, value is the value to validate, and validateFuncs are the validation rules to apply.
func Value[T any](name string, value T, validateFuncs ...Func[T]) Validator {
	return func() error {
		for _, validate := range validateFuncs {
			if err := validate(name, value); err != nil {
				return err
			}
		}
		return nil
	}
}

// All runs all validators and collects all validation errors (if any).
// Returns a joined error containing all validation failures.
func All(validators ...Validator) error {
	errs := make([]error, len(validators))
	for i, validate := range validators {
		if err := validate(); err != nil {
			errs[i] = err
		}
	}
	return errors.Join(errs...)
}

// Chain runs validators one after another, stopping at the first validation error.
// Returns the first error encountered or nil if all validations pass.
func Chain(validators ...Validator) error {
	for _, validate := range validators {
		if err := validate(); err != nil {
			return err
		}
	}
	return nil
}

// MinLength returns a validation function that checks if a string has at least count runes.
func MinLength(count int) Func[string] {
	return func(name string, value string) error {
		if utf8.RuneCountInString(value) < count {
			return ValidationError{
				Field: name,
				Msg:   fmt.Sprintf("must be at least %d characters long", count),
			}
		}
		return nil
	}
}

// MaxLength returns a validation function that checks if a string has no more than count runes.
func MaxLength(count int) Func[string] {
	return func(name string, value string) error {
		if utf8.RuneCountInString(value) > count {
			return ValidationError{
				Field: name,
				Msg:   fmt.Sprintf("must be no more than %d characters long", count),
			}
		}
		return nil
	}
}

// RegExp returns a validation function that checks if a string matches the given regular expression.
func RegExp(re *regexp.Regexp) Func[string] {
	return func(name string, value string) error {
		if !re.MatchString(value) {
			return ValidationError{
				Field: name,
				Msg:   fmt.Sprintf("must match pattern %s", re),
			}
		}
		return nil
	}
}

// OneOf returns a validation function that checks if a value is one of the allowed values.
func OneOf[T comparable](allowed ...T) Func[T] {
	return func(name string, value T) error {
		for _, allowedValue := range allowed {
			if value == allowedValue {
				return nil
			}
		}
		return ValidationError{
			Field: name,
			Msg:   fmt.Sprintf("must be one of %v", allowed),
		}
	}
}

// Min returns a validation function that checks if a value is at least the given minimum.
// Works with any ordered type (implementing cmp.Ordered).
func Min[T cmp.Ordered](minimum T) Func[T] {
	return func(name string, value T) error {
		if value < minimum {
			return ValidationError{
				Field: name,
				Msg:   fmt.Sprintf("must be at least %v", minimum),
			}
		}
		return nil
	}
}

// Max returns a validation function that checks if a value is no more than the given maximum.
// Works with any ordered type (implementing cmp.Ordered).
func Max[T cmp.Ordered](maximum T) Func[T] {
	return func(name string, value T) error {
		if value > maximum {
			return ValidationError{
				Field: name,
				Msg:   fmt.Sprintf("must be no more than %v", maximum),
			}
		}
		return nil
	}
}

// MinItems returns a validation function that checks if a slice has at least count items.
func MinItems[T any](count int) Func[[]T] {
	return func(name string, value []T) error {
		if len(value) < count {
			return ValidationError{
				Field: name,
				Msg:   fmt.Sprintf("must contain at least %d items", count),
			}
		}
		return nil
	}
}

// MaxItems returns a validation function that checks if a slice has no more than count items.
func MaxItems[T any](count int) Func[[]T] {
	return func(name string, value []T) error {
		if len(value) > count {
			return ValidationError{
				Field: name,
				Msg:   fmt.Sprintf("must contain no more than %d items", count),
			}
		}
		return nil
	}
}

// UniqueItems returns a validation function that checks if all items in a slice are unique.
func UniqueItems[T comparable]() Func[[]T] {
	return func(name string, value []T) error {
		if len(value) <= 1 {
			return nil
		}
		seen := make(map[T]bool, len(value))
		for _, item := range value {
			if seen[item] {
				return ValidationError{
					Field: name,
					Msg:   "must contain unique items, duplicate found",
				}
			}
			seen[item] = true
		}
		return nil
	}
}
