package errors_test

import (
	"testing"

	"github.com/Zamony/go/errors"
)

type ConstantError string

func (e ConstantError) Error() string { return string(e) }

func TestConstantErrors(t *testing.T) {
	const (
		errInternal = ConstantError("internal")
		errDefault  = errors.ConstantError("default")
	)

	errInternalWrap := errors.Wrapf(errors.Wrapf(errInternal, "foo"), "bar")
	errDefaultWrap := errors.Wrapf(errors.Wrapf(errDefault, "foo"), "bar")

	if errors.Is(errInternalWrap, errDefault) {
		t.Error("Wrapped internal error is default")
	}

	if !errors.Is(errDefaultWrap, errDefault) {
		t.Error("Wrapped default error is not default")
	}

	var internalTarget ConstantError
	if errors.As(errDefaultWrap, &internalTarget) {
		t.Error("Wrapped default error is assignable to internal")
	}

	var defaultTarget errors.ConstantError
	if !errors.As(errDefaultWrap, &defaultTarget) {
		t.Error("Wrapped default error is not assignable to default")
	}

	if defaultTarget != errDefault {
		t.Errorf("Default target '%v' is not equal to '%v'", defaultTarget, errDefault)
	}
}
