package errors_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/Zamony/go/errors"
)

func newError(text string) error {
	return errors.New(text)
}

func TestMessageAndStack(t *testing.T) {
	testCases := []struct {
		TestName string
		Error    error
		Message  string
		Stack    string
		Verbose  string
	}{
		{
			TestName: "Error with stack",
			Error:    errors.New("new"),
			Message:  "new",
			Stack:    "errors_test.TestMessageAndStack:25/testing.tRunner/runtime.goexit",
			Verbose:  "new: errors_test.TestMessageAndStack:25/testing.tRunner/runtime.goexit",
		},
		{
			TestName: "Formatted error with stack",
			Error:    errors.Newf("%s", "newf"),
			Message:  "newf",
			Stack:    "errors_test.TestMessageAndStack:32/testing.tRunner/runtime.goexit",
			Verbose:  "newf: errors_test.TestMessageAndStack:32/testing.tRunner/runtime.goexit",
		},
		{
			TestName: "Sentinel error",
			Error:    errors.SentinelError("sentinel"),
			Message:  "sentinel",
			Stack:    "",
			Verbose:  "sentinel",
		},
		{
			TestName: "Wrapped sentinel error",
			Error:    errors.Wrap(errors.SentinelError("sentinel wrap")),
			Message:  "sentinel wrap",
			Stack:    "errors_test.TestMessageAndStack:46/testing.tRunner/runtime.goexit",
			Verbose:  "sentinel wrap: errors_test.TestMessageAndStack:46/testing.tRunner/runtime.goexit",
		},
		{
			TestName: "Wrapped formatted sentinel error",
			Error:    errors.Wrapf(errors.SentinelError("sentinel"), "%s", "wrapf"),
			Message:  "wrapf: sentinel",
			Stack:    "errors_test.TestMessageAndStack:53/testing.tRunner/runtime.goexit",
			Verbose:  "wrapf: sentinel: errors_test.TestMessageAndStack:53/testing.tRunner/runtime.goexit",
		},
		{
			TestName: "Wrapped error with stack",
			Error:    errors.Wrap(newError("new")),
			Message:  "new",
			Stack:    "errors_test.newError:12/TestMessageAndStack/testing.tRunner/runtime.goexit",
			Verbose:  "new: errors_test.newError:12/TestMessageAndStack/testing.tRunner/runtime.goexit",
		},
		{
			TestName: "Wrapped formatted error with stack",
			Error:    errors.Wrapf(newError("new"), "%s", "wrapf"),
			Message:  "wrapf: new",
			Stack:    "errors_test.newError:12/TestMessageAndStack/testing.tRunner/runtime.goexit",
			Verbose:  "wrapf: new: errors_test.newError:12/TestMessageAndStack/testing.tRunner/runtime.goexit",
		},
		{
			TestName: "Joined sentinel errors",
			Error:    errors.Join(errors.SentinelError("foo"), errors.SentinelError("bar"), errors.SentinelError("baz")),
			Message:  `["foo" "bar" "baz"]`,
			Stack:    "errors_test.TestMessageAndStack:74/testing.tRunner/runtime.goexit",
			Verbose:  `["foo" "bar" "baz"]`,
		},
		{
			TestName: "Joined errors with stack",
			Error:    errors.Join(errors.New("foo"), newError("bar")),
			Message:  `["foo" "bar"]`,
			Stack:    "errors_test.TestMessageAndStack:81/testing.tRunner/runtime.goexit",
			Verbose:  `["foo: errors_test.TestMessageAndStack:81/testing.tRunner/runtime.goexit" "bar: errors_test.newError:12/TestMessageAndStack/testing.tRunner/runtime.goexit"]`,
		},
		{
			TestName: "Joined nested errors",
			Error:    errors.Join(errors.SentinelError("foo"), errors.Join(errors.SentinelError("bar"), errors.SentinelError("baz"))),
			Message:  `["foo" "[\"bar\" \"baz\"]"]`,
			Stack:    "errors_test.TestMessageAndStack:88/testing.tRunner/runtime.goexit",
			Verbose:  `["foo" "[\"bar\" \"baz\"]"]`,
		},
		{
			TestName: "Join single stacked error and nils",
			Error:    errors.Join(nil, newError("foo"), nil),
			Message:  "foo",
			Stack:    "errors_test.newError:12/TestMessageAndStack/testing.tRunner/runtime.goexit",
			Verbose:  "foo: errors_test.newError:12/TestMessageAndStack/testing.tRunner/runtime.goexit",
		},
		{
			TestName: "Join single sentinel error and nils",
			Error:    errors.Join(nil, nil, errors.SentinelError("foo")),
			Message:  "foo",
			Stack:    "errors_test.TestMessageAndStack:102/testing.tRunner/runtime.goexit",
			Verbose:  "foo: errors_test.TestMessageAndStack:102/testing.tRunner/runtime.goexit",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			equal(t, errors.StackTrace(tc.Error), tc.Stack)
			equal(t, fmt.Sprintf("%s", tc.Error), tc.Message)
			equal(t, fmt.Sprintf("%+v", tc.Error), tc.Verbose)
		})
	}
}

func TestErrorIs(t *testing.T) {
	var (
		errStack    = errors.New("stackful")
		errSentinel = errors.SentinelError("sentinel")
		errJoin     = errors.Join(errStack, errSentinel)
	)
	testCases := []struct {
		TestName string
		Error    error
		Target   error
		Result   bool
	}{
		{
			TestName: "Stackful error itself",
			Error:    errStack,
			Target:   errStack,
			Result:   true,
		},
		{
			TestName: "Different stackful errors",
			Error:    errStack,
			Target:   errors.New("foo"),
			Result:   false,
		},
		{
			TestName: "Sentinel error itself",
			Error:    errSentinel,
			Target:   errSentinel,
			Result:   true,
		},
		{
			TestName: "Different sentinel errors",
			Error:    errSentinel,
			Target:   errors.SentinelError("sentinel"),
			Result:   false,
		},
		{
			TestName: "Joined error itself",
			Error:    errJoin,
			Target:   errJoin,
			Result:   true,
		},
		{
			TestName: "Different joined errors",
			Error:    errJoin,
			Target:   errors.Join(errStack, errSentinel),
			Result:   false,
		},
		{
			TestName: "Joined stackful error",
			Error:    errJoin,
			Target:   errStack,
			Result:   true,
		},
		{
			TestName: "Joined sentinel error",
			Error:    errJoin,
			Target:   errSentinel,
			Result:   true,
		},
		{
			TestName: "Wrapped stackful error",
			Error:    errors.Wrap(errStack),
			Target:   errStack,
			Result:   true,
		},
		{
			TestName: "Wrapped joined error",
			Error:    errors.Wrap(errJoin),
			Target:   errJoin,
			Result:   true,
		},
		{
			TestName: "Wrapped sentinel error",
			Error:    errors.Wrap(errSentinel),
			Target:   errSentinel,
			Result:   true,
		},
		{
			TestName: "Wrapfed sentinel error",
			Error:    errors.Wrapf(errSentinel, "%s", "wrap"),
			Target:   errSentinel,
			Result:   true,
		},
		{
			TestName: "Wrapped joined sentinel error",
			Error:    errors.Wrap(errors.Join(errSentinel, errStack)),
			Target:   errSentinel,
			Result:   true,
		},
		{
			TestName: "Wrapped nested sentinel error",
			Error:    errors.Wrapf(errors.Join(errJoin, errors.Wrap(errors.Join(errSentinel, errStack))), ""),
			Target:   errSentinel,
			Result:   true,
		},
		{
			TestName: "Wrapped nil error",
			Error:    errors.Wrap(nil),
			Target:   nil,
			Result:   true,
		},
		{
			TestName: "Wrapfed nil error",
			Error:    errors.Wrapf(nil, "%s", "nil"),
			Target:   nil,
			Result:   true,
		},
		{
			TestName: "Empty join",
			Error:    errors.Join(),
			Target:   nil,
			Result:   true,
		},
		{
			TestName: "Join of nil",
			Error:    errors.Join(nil),
			Target:   nil,
			Result:   true,
		},
		{
			TestName: "Join of multiple nils",
			Error:    errors.Join(nil, nil),
			Target:   nil,
			Result:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			equal(t, errors.Is(tc.Error, tc.Target), tc.Result)
		})
	}
}

type CustomError struct {
	Message string
}

func (c CustomError) Error() string { return c.Message }

func TestErrorAs(t *testing.T) {
	var (
		errStack    = errors.New("stackful")
		errSentinel = errors.SentinelError("sentinel")
		errJoin     = errors.Join(errStack, errSentinel)
		stackTarget interface{ StackTrace() errors.StackFrames }
	)

	testCases := []struct {
		TestName string
		Error    error
		Target   any
		Result   error
	}{
		{
			TestName: "Stackful error itself",
			Error:    errStack,
			Target:   &stackTarget,
			Result:   errStack,
		},
		{
			TestName: "Joined error itself",
			Error:    errJoin,
			Target:   &stackTarget,
			Result:   errJoin,
		},
		{
			TestName: "Wrapped custom error",
			Error:    errors.Wrap(CustomError{"custom wrap"}),
			Target:   &CustomError{},
			Result:   CustomError{"custom wrap"},
		},
		{
			TestName: "Wrapfed custom error",
			Error:    errors.Wrapf(CustomError{"custom wrapf"}, "wrapf"),
			Target:   &CustomError{},
			Result:   CustomError{"custom wrapf"},
		},
		{
			TestName: "Joined custom error",
			Error:    errors.Join(errSentinel, CustomError{"custom join"}),
			Target:   &CustomError{},
			Result:   CustomError{"custom join"},
		},
		{
			TestName: "No joined custom error",
			Error:    errors.Join(errSentinel, errStack),
			Target:   &CustomError{},
			Result:   nil,
		},
		{
			TestName: "No wrapped custom error",
			Error:    errors.Wrap(errSentinel),
			Target:   &CustomError{},
			Result:   nil,
		},
		{
			TestName: "No wrapfed custom error",
			Error:    errors.Wrapf(errSentinel, "wrapf"),
			Target:   &CustomError{},
			Result:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.TestName, func(t *testing.T) {
			found := tc.Result != nil
			equal(t, errors.As(tc.Error, tc.Target), found)
			if found {
				equal(t, unptr(tc.Target), tc.Result)
			}
		})
	}
}

func unptr(target any) any {
	if v := reflect.ValueOf(target); v.Kind() == reflect.Ptr {
		return v.Elem().Interface()
	}
	return target
}

func equal(t *testing.T, got, want any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Not equal (-want, +got):\n- %+v\n+ %+v\n", want, got)
	}
}
