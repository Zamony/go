package errors_test

import (
	"testing"

	"github.com/Zamony/go/errors"
)

func TestStacktrace(t *testing.T) {
	got := errors.Stacktrace(errors.ConstantError("myerror"))
	if got != "" {
		t.Errorf("Stacktrace is not empty for constant error: %q", got)
	}
}
