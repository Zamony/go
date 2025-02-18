package optional

import (
	"bytes"
	"encoding/json"
	"errors"
	"testing"
)

// custom type used to check marshaling/unmarshaling for a non-primitive type.
type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func TestSomeAndNone(t *testing.T) {
	// Test with an int value
	optInt := Some(42)
	if value, ok := optInt.Get(); !ok || value != 42 {
		t.Errorf("Expected Some(42) to return (42, true), got (%v, %v)", value, ok)
	}

	noneInt := None[int]()
	if value, ok := noneInt.Get(); ok {
		t.Errorf("Expected None to return (zero value, false), got (%v, %v)", value, ok)
	}

	// Test with a string value
	optStr := Some("hello")
	if value, ok := optStr.Get(); !ok || value != "hello" {
		t.Errorf("Expected Some(\"hello\") to return (\"hello\", true), got (%v, %v)", value, ok)
	}

	noneStr := None[string]()
	if value, ok := noneStr.Get(); ok {
		t.Errorf("Expected None to return (zero value, false), got (%v, %v)", value, ok)
	}
}

func TestGetOrElse(t *testing.T) {
	opt := Some(100)
	if v := opt.GetOrElse(0); v != 100 {
		t.Errorf("Expected GetOrElse on Some to return contained value (100), got %d", v)
	}

	none := None[int]()
	if v := none.GetOrElse(50); v != 50 {
		t.Errorf("Expected GetOrElse on None to return default (50), got %d", v)
	}
}

func TestMarshalJSON(t *testing.T) {
	// Test marshaling an Optional with a value
	opt := Some("test")
	data, err := json.Marshal(&opt)
	if err != nil {
		t.Fatalf("Unexpected error marshaling valid Optional: %v", err)
	}
	expected, _ := json.Marshal("test")
	if !bytes.Equal(data, expected) {
		t.Errorf("Expected JSON %s, got %s", expected, data)
	}

	// Test marshaling an empty Optional (None)
	none := None[string]()
	data, err = json.Marshal(&none)
	if err != nil {
		t.Fatalf("Unexpected error marshaling None Optional: %v", err)
	}
	if !bytes.Equal(data, []byte("null")) {
		t.Errorf("Expected JSON null, got %s", data)
	}
}

func TestUnmarshalJSON(t *testing.T) {
	// Test unmarshaling a valid value into the Optional
	var opt Optional[int]
	err := json.Unmarshal([]byte("123"), &opt)
	if err != nil {
		t.Fatalf("Unexpected error unmarshaling value into Optional: %v", err)
	}
	value, ok := opt.Get()
	if !ok || value != 123 {
		t.Errorf("Expected Optional to contain 123, got (%v, %v)", value, ok)
	}

	// Test unmarshaling JSON null; should result in a None.
	err = json.Unmarshal([]byte("null"), &opt)
	if err != nil {
		t.Fatalf("Unexpected error unmarshaling null into Optional: %v", err)
	}
	_, ok = opt.Get()
	if ok {
		t.Error("Expected Optional to be None after unmarshaling null")
	}

	// Test with non-primitive type (struct)
	var optPerson Optional[Person]
	personJSON := `{"name": "Alice", "age": 30}`
	err = json.Unmarshal([]byte(personJSON), &optPerson)
	if err != nil {
		t.Fatalf("Unexpected error unmarshaling person JSON: %v", err)
	}
	p, ok := optPerson.Get()
	if !ok {
		t.Fatal("Expected Optional person to be present")
	}
	if p.Name != "Alice" || p.Age != 30 {
		t.Errorf("Unexpected person value: got %+v", p)
	}
}

func TestUnmarshalJSONError(t *testing.T) {
	// Provide malformed JSON to cause error in unmarshaling.
	var opt Optional[int]
	invalidJSON := []byte(`"not an int"`)
	err := json.Unmarshal(invalidJSON, &opt)
	if err == nil {
		t.Error("Expected error when unmarshaling invalid JSON for int, got nil")
	}
	// Verify that Optional remains in None state
	_, ok := opt.Get()
	if ok {
		t.Error("Expected Optional to be None after failed unmarshaling")
	}
}

func TestMultipleOperations(t *testing.T) {
	// Test a sequence of operations: set, marshal, unmarshal, then get value.
	original := Some(Person{Name: "Bob", Age: 45})
	data, err := json.Marshal(&original)
	if err != nil {
		t.Fatalf("Unexpected error marshaling person: %v", err)
	}

	var decoded Optional[Person]
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatalf("Unexpected error unmarshaling person: %v", err)
	}

	p, ok := decoded.Get()
	if !ok {
		t.Error("Expected Optional person to be present after unmarshaling")
	}
	if p.Name != "Bob" || p.Age != 45 {
		t.Errorf("After decoding, expected person Bob (45), got %+v", p)
	}
}

func TestMarshalThenUnmarshalNone(t *testing.T) {
	// Ensure that marshaling None produces null and unmarshaling null yields None.
	none := None[float64]()
	data, err := json.Marshal(&none)
	if err != nil {
		t.Fatalf("Error marshaling None Optional: %v", err)
	}
	if !bytes.Equal(data, []byte("null")) {
		t.Errorf("Expected marshaled data to be null, got %s", data)
	}

	var decoded Optional[float64]
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Errorf("Error unmarshaling JSON null: %v", err)
	}
	_, ok := decoded.Get()
	if ok {
		t.Error("Expected Optional to be None after unmarshaling null")
	}
}

func TestGetAndGetOrElseForComplexType(t *testing.T) {
	// Test with a slice as the contained type.
	opt := Some([]int{1, 2, 3})
	val, ok := opt.Get()
	if !ok {
		t.Error("Expected Optional slice to be present")
	}
	if len(val) != 3 {
		t.Errorf("Expected slice length 3, got %d", len(val))
	}

	none := None[[]int]()
	result := none.GetOrElse([]int{4, 5})
	if len(result) != 2 || result[0] != 4 || result[1] != 5 {
		t.Errorf("Expected default slice [4,5], got %v", result)
	}
}

func TestEdgeCaseEmptyJSON(t *testing.T) {
	// In some cases, an empty byte slice should generate an error.
	var opt Optional[string]
	err := json.Unmarshal([]byte(""), &opt)
	if err == nil {
		t.Error("Expected error when unmarshaling empty JSON, got nil")
	}
}

// Testing error propagation is a bit tricky since our MarshalJSON implementation simply
// marshals the contained value. We can simulate an error by creating a type that returns error on marshaling.
type errorMarshaler struct{}

// Implement json.Marshaler for errorMarshaler that always errors.
func (e errorMarshaler) MarshalJSON() ([]byte, error) {
	return nil, errors.New("marshal error")
}

func TestMarshalJSONErrorPropagation(t *testing.T) {
	opt := Some(errorMarshaler{})
	_, err := json.Marshal(&opt)
	if err == nil {
		t.Error("Expected error during marshaling of errorMarshaler, got nil")
	}
}
