package toast

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

// Equal checks if two values are equal.
// If they are not equal, it logs the difference in the error log,
// marks the function as having failed and continues its execution.
func Equal(t *testing.T, x, y any, opts ...cmp.Option) {
	t.Helper()
	if diff := cmp.Diff(x, y, opts...); diff != "" {
		t.Error(diff)
	}
}

// MustEqual checks if two values are equal.
// If they are not equal, it logs the difference in the error log,
// marks the function as having failed and stops its execution.
func MustEqual(t *testing.T, x, y any, opts ...cmp.Option) {
	t.Helper()
	if diff := cmp.Diff(x, y, opts...); diff != "" {
		t.Fatal(diff)
	}
}
