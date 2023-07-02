package errors_test

import (
	"fmt"
	"testing"

	"github.com/Zamony/go/errors"
)

const errBar = errors.ConstantError("bar")

func newFooError() error {
	return errors.New("foo")
}

func TestLog(t *testing.T) {
	testCases := []struct {
		Error       error
		WantMessage string
		WantStack   string
	}{
		{errors.New("new"), "new", "errors_test.TestLog:22/testing.tRunner/runtime.goexit"},
		{errors.Newf("newf: %d", 1), "newf: 1", "errors_test.TestLog:23/testing.tRunner/runtime.goexit"},
		{errors.Wrapf(errors.New("inner"), "wrap%d", 1), "wrap1: inner", "errors_test.TestLog:24/testing.tRunner/runtime.goexit"},
		{errors.Wrapf(newFooError(), "wrap%d", 1), "wrap1: foo", "errors_test.newFooError:13/TestLog/testing.tRunner/runtime.goexit"},
		{errors.Wrapf(errBar, "wrap%d", 1), "wrap1: bar", "errors_test.TestLog:26/testing.tRunner/runtime.goexit"},
		{errors.Wrap(errBar), "bar", "errors_test.TestLog:27/testing.tRunner/runtime.goexit"},
		{errors.Wrap(newFooError()), "foo", "errors_test.newFooError:13/TestLog/testing.tRunner/runtime.goexit"},
	}

	for i, tc := range testCases {
		if msg := tc.Error.Error(); msg != tc.WantMessage {
			t.Errorf("Test case #%d. Message mismatch: want %q != %q", i, tc.WantMessage, msg)
		}
		if stack := errors.Stacktrace(tc.Error); stack != tc.WantStack {
			t.Errorf("Test case #%d. Stack mismatch: want %q != %q", i, tc.WantStack, stack)
		}
	}
}

func TestBasicErrorFormat(t *testing.T) {
	const (
		errMsg   = "errFoo"
		errStack = "errors_test.TestBasicErrorFormat:47/testing.tRunner/runtime.goexit"
	)

	err := errors.New(errMsg)
	if got := fmt.Sprintf("%s", err); errMsg != got {
		t.Errorf("Mismatched %%s verb: got %q, want %q", got, errMsg)
	}
	if got := fmt.Sprintf("%q", err); errMsg != got {
		t.Errorf("Mismatched %%q verb: got %q, want %q", got, errMsg)
	}
	if got := fmt.Sprintf("%v", err); errMsg != got {
		t.Errorf("Mismatched %%v verb: got %q, want %q", got, errMsg)
	}

	want := fmt.Sprintf("%s: %s", errMsg, errStack)
	if got := fmt.Sprintf("%+v", err); want != got {
		t.Errorf("Mismatched %%+v verb: got %q, want %q", got, want)
	}
}

func TestBasicErrorWrap(t *testing.T) {
	if has := errors.Is(errors.Wrap(errBar), errBar); !has {
		t.Error("Wrap doesn't chain errors")
	}
	if has := errors.Is(errors.Wrapf(errBar, "wrap%d", 1), errBar); !has {
		t.Error("Wrapf doesn't chain errors")
	}

	var target errors.ConstantError
	if has := errors.As(errors.Wrapf(errBar, "wrap%d", 1), &target); !has {
		t.Error("Target error was not found in the chain")
	} else if target != errBar {
		t.Error("Target error doesn't match original error")
	}

	if errors.Wrap(nil) != nil {
		t.Error("Wrap of nil is not nil")
	}

	if errors.Wrapf(nil, "msg") != nil {
		t.Error("Wrapf of nil is not nil")
	}
}
