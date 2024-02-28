package errors_test

import (
	"fmt"
	"testing"

	"github.com/Zamony/go/errors"
)

func TestCombine(t *testing.T) {
	if err := errors.Combine(nil, nil, nil); err != nil {
		t.Errorf("Combined nil errors are not nil: got %v", err)
	}

	var (
		errFoo = errors.SentinelError("errFoo")
		errBar = errors.SentinelError("errBar")
	)
	err := errors.Combine(errFoo, nil, errBar)
	if !errors.Is(err, errFoo) {
		t.Errorf("Not a foo error")
	}
	if !errors.Is(err, errBar) {
		t.Errorf("Not a bar error")
	}

	errStack := errors.New("error with stack")
	err = errors.Combine(errFoo, errStack)
	got, want := fmt.Sprintf("%v", err), `errFoo; error with stack`
	if got != want {
		t.Errorf("Combined %%v mismatch: got %q, want %q", got, want)
	}

	got, want = fmt.Sprintf("%+v", err), `errFoo; error with stack: errors_test.TestCombine:27/testing.tRunner/runtime.goexit`
	if got != want {
		t.Errorf("Combined %%+v mismatch: got %q, want %q", got, want)
	}
}
