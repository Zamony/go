// Package optional provides a type for representing optional values.
package optional

import (
	"bytes"
	"encoding/json"
)

// Optional is a type that can hold a value of any type, or no value at all.
type Optional[T any] struct {
	value   T
	present bool
}

// Some creates an Optional instance that contains a value.
func Some[T any](value T) Optional[T] {
	return Optional[T]{
		value:   value,
		present: true,
	}
}

// None creates an empty Optional instance.
func None[T any]() Optional[T] {
	return Optional[T]{}
}

// Get returns the value and a boolean indicating if a value is present.
func (o *Optional[T]) Get() (value T, ok bool) {
	return o.value, o.present
}

// GetOrElse returns the value if present, or the provided default value if empty.
func (o *Optional[T]) GetOrElse(defaultValue T) T {
	if o.present {
		return o.value
	}
	return defaultValue
}

// MarshalJSON implements the json.Marshaler interface.
// If the Optional has no value, it returns JSON null; otherwise, it marshals the contained value.
func (o Optional[T]) MarshalJSON() ([]byte, error) {
	if !o.present {
		return []byte("null"), nil
	}
	return json.Marshal(o.value)
}

var nullBytes = []byte("null")

// UnmarshalJSON implements the json.Unmarshaler interface.
// If the data is JSON null, the Optional remains empty. Otherwise, it decodes the data into the value.
func (o *Optional[T]) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, nullBytes) {
		*o = None[T]()
		return nil
	}
	var v T
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}
	*o = Some(v)
	return nil
}
