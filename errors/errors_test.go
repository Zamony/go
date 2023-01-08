package errors

import (
	"testing"

	"github.com/Zamony/go/errors/stackerr"
)

func newFooError() error {
	return New("foo")
}

const errBar = ConstantError("bar")

func TestLog(t *testing.T) {
	testCases := []struct {
		Error       error
		WantMessage string
		WantStack   string
	}{
		{New("new"), "new", "errors.TestLog:21/testing.tRunner/runtime.goexit"},
		{Newf("newf: %d", 1), "newf: 1", "errors.TestLog:22/testing.tRunner/runtime.goexit"},
		{Wrap(New("inner"), "wrap"), "wrap: inner", "errors.TestLog:23/testing.tRunner/runtime.goexit"},
		{Wrapf(New("inner"), "wrap%d", 1), "wrap1: inner", "errors.TestLog:24/testing.tRunner/runtime.goexit"},
		{Wrap(newFooError(), "wrap"), "wrap: foo", "errors.newFooError:10/TestLog/testing.tRunner/runtime.goexit"},
		{Wrapf(newFooError(), "wrap%d", 1), "wrap1: foo", "errors.newFooError:10/TestLog/testing.tRunner/runtime.goexit"},
		{Wrap(errBar, "wrap"), "wrap: bar", "errors.TestLog:27/testing.tRunner/runtime.goexit"},
		{Wrapf(errBar, "wrap%d", 1), "wrap1: bar", "errors.TestLog:28/testing.tRunner/runtime.goexit"},
		{From(errBar), "bar", "errors.TestLog:29/testing.tRunner/runtime.goexit"},
		{From(newFooError()), "foo", "errors.newFooError:10/TestLog/testing.tRunner/runtime.goexit"},
	}

	for i, tc := range testCases {
		if msg := tc.Error.Error(); msg != tc.WantMessage {
			t.Fatalf("Test case #%d. Message mismatch: want %q != %q", i, tc.WantMessage, msg)
		}
		if stack := stackerr.MarshalCompact(tc.Error); stack != tc.WantStack {
			t.Fatalf("Test case #%d. Stack mismatch: want %q != %q", i, tc.WantStack, stack)
		}
	}
}

func TestWrap(t *testing.T) {
	if has := Is(Wrap(errBar, "wrap"), errBar); !has {
		t.Fatal("Wrap doesn't chain error")
	}
	if has := Is(Wrapf(errBar, "wrap%d", 1), errBar); !has {
		t.Fatal("Wrapf doesn't chain error")
	}
}
