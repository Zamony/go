package validate_test

import (
	"regexp"
	"testing"

	"github.com/Zamony/go/validate"
)

func TestValidationError_Error(t *testing.T) {
	err := validate.ValidationError{
		Field: "username",
		Msg:   "must be at least 5 characters long",
	}
	expected := "username: must be at least 5 characters long"
	if actual := err.Error(); actual != expected {
		t.Errorf("expected %q, got %q", expected, actual)
	}
}

func TestValue(t *testing.T) {
	t.Run("valid value passes all validators", func(t *testing.T) {
		validator := validate.Value("age", 25, validate.Min(18), validate.Max(30))
		if err := validator(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid value fails validation", func(t *testing.T) {
		validator := validate.Value("age", 15, validate.Min(18))
		err := validator()
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		validationErr, ok := err.(validate.ValidationError)
		if !ok {
			t.Fatalf("expected ValidationError, got %T", err)
		}
		if validationErr.Field != "age" {
			t.Errorf("expected field 'age', got %q", validationErr.Field)
		}
	})
}

func TestAll(t *testing.T) {
	t.Run("all validators pass", func(t *testing.T) {
		err := validate.All(
			validate.Value("username", "user123", validate.MinLength(5)),
			validate.Value("age", 25, validate.Min(18)),
		)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("multiple validation errors", func(t *testing.T) {
		err := validate.All(
			validate.Value("username", "usr", validate.MinLength(5)),
			validate.Value("age", 15, validate.Min(18)),
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if errors := err.(interface{ Unwrap() []error }).Unwrap(); len(errors) != 2 {
			t.Errorf("expected 2 errors, got %d", len(errors))
		}
	})
}

func TestChain(t *testing.T) {
	t.Run("stops on first error", func(t *testing.T) {
		var called bool
		customValidator := func() error {
			called = true
			return nil
		}

		err := validate.Chain(
			validate.Value("username", "usr", validate.MinLength(5)),
			customValidator,
		)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if called {
			t.Error("subsequent validator should not be called after error")
		}
	})
}

func TestMinLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		min       int
		expectErr bool
	}{
		{"valid", "hello", 5, false},
		{"too short", "hi", 5, true},
		{"unicode", "привет", 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.MinLength(tt.min)("field", tt.value)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMaxLength(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		max       int
		expectErr bool
	}{
		{"valid", "hello", 5, false},
		{"too long", "hello world", 5, true},
		{"unicode", "привет", 6, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.MaxLength(tt.max)("field", tt.value)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestRegExp(t *testing.T) {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	validator := validate.RegExp(emailRegex)

	tests := []struct {
		name      string
		email     string
		expectErr bool
	}{
		{"valid", "test@example.com", false},
		{"invalid", "not-an-email", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator("email", tt.email)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestOneOf(t *testing.T) {
	validator := validate.OneOf("red", "green", "blue")

	tests := []struct {
		name      string
		color     string
		expectErr bool
	}{
		{"valid", "green", false},
		{"invalid", "yellow", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator("color", tt.color)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMin(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		min       int
		expectErr bool
	}{
		{"valid", 10, 5, false},
		{"equal", 5, 5, false},
		{"invalid", 3, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Min(tt.min)("number", tt.value)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMax(t *testing.T) {
	tests := []struct {
		name      string
		value     int
		max       int
		expectErr bool
	}{
		{"valid", 5, 10, false},
		{"equal", 5, 5, false},
		{"invalid", 8, 5, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Max(tt.max)("number", tt.value)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMinItems(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		min       int
		expectErr bool
	}{
		{"valid", []int{1, 2, 3}, 2, false},
		{"equal", []int{1, 2}, 2, false},
		{"invalid", []int{1}, 2, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.MinItems[int](tt.min)("items", tt.items)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestMaxItems(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		max       int
		expectErr bool
	}{
		{"valid", []int{1, 2}, 3, false},
		{"equal", []int{1, 2, 3}, 3, false},
		{"invalid", []int{1, 2, 3, 4}, 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.MaxItems[int](tt.max)("items", tt.items)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestUniqueItems(t *testing.T) {
	tests := []struct {
		name      string
		items     []int
		expectErr bool
	}{
		{"valid", []int{1, 2, 3}, false},
		{"empty", []int{}, false},
		{"single", []int{1}, false},
		{"duplicates", []int{1, 2, 2}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.UniqueItems[int]()("items", tt.items)
			if tt.expectErr && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
