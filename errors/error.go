package errors

import (
	stderrors "errors"
	"fmt"
	"io"
	"runtime"
)

type baseError struct {
	err    error
	frames StackFrames
}

func (baseError) Is(error) bool {
	return false // don't use stackful errors as sentinel ones
}

func (b baseError) Unwrap() error {
	return b.err
}

type StackFrames []uintptr

func (b baseError) StackTrace() StackFrames {
	return copystack(b.frames)
}

func copystack(stack StackFrames) StackFrames {
	r := make(StackFrames, len(stack))
	copy(r, stack)
	return r
}

func (b baseError) Error() string {
	return b.err.Error()
}

func (b baseError) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			if trace := StackTrace(b); trace != "" {
				_, _ = fmt.Fprintf(state, "%s: %s", b.Error(), trace)
				return
			}
		}
		fallthrough
	case 's', 'q':
		io.WriteString(state, b.Error())
	}
}

// New creates an error with message and a stacktrace.
func New(text string) error {
	return baseError{stderrors.New(text), frames()}
}

// Newf creates an error with formatted message and a stacktrace.
func Newf(format string, a ...any) error {
	return baseError{stderrors.New(fmt.Sprintf(format, a...)), frames()}
}

const maxStackDepth = 32

func frames() []uintptr {
	const skip = 3 // +1 Callers, +1 frames, +1 New
	var pcs [maxStackDepth]uintptr
	depth := runtime.Callers(skip, pcs[:])
	return pcs[:depth]
}

type stacktracer interface {
	StackTrace() StackFrames
}

// Wrap adds a stacktrace to the error if it doesn't have one.
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(stacktracer); ok {
		return err
	}
	return baseError{err, frames()}
}

// Wrapf annotates an error with a message.
// Also adds a stacktrace to the error if it doesn't have one.
func Wrapf(err error, format string, a ...any) error {
	if err == nil {
		return nil
	}

	msg := fmt.Sprintf(format, a...)
	newErr := fmt.Errorf("%s: %w", msg, err)
	if b, ok := err.(stacktracer); ok {
		return baseError{newErr, b.StackTrace()}
	}
	return baseError{newErr, frames()}
}

type joinError struct {
	errors []error
	frames StackFrames
}

func (e joinError) StackTrace() StackFrames {
	return copystack(e.frames)
}

func (joinError) Is(error) bool {
	return false // it only wraps other errors
}

func (e joinError) Unwrap() []error {
	return e.errors
}

func (e joinError) toString(format func(error) string) string {
	errs := make([]string, len(e.errors))
	for i := range e.errors {
		errs[i] = format(e.errors[i])
	}

	return fmt.Sprintf("%q", errs)
}

func (e joinError) Error() string {
	return e.toString(func(err error) string {
		return err.Error()
	})
}

func (e joinError) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			io.WriteString(state, e.toString(func(err error) string {
				return fmt.Sprintf("%+v", err)
			}))
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(state, e.Error())
	}
}

// Join joins multiple errors into one.
// Nil errors are ignored.
func Join(errs ...error) error {
	// Could be Join(errA, errB error, errsTail ...error) error.
	// But it would make it impossible to unpack errors: Join(errors...).
	var stack StackFrames
	jerrs := make([]error, 0, len(errs))
	for _, err := range errs {
		if err == nil {
			continue
		}
		if len(stack) == 0 {
			if b, ok := err.(stacktracer); ok {
				stack = b.StackTrace()
			}
		}
		jerrs = append(jerrs, err)
	}

	if len(jerrs) == 0 {
		return nil
	}
	if len(stack) == 0 {
		stack = frames()
	}
	if len(jerrs) == 1 {
		return baseError{jerrs[0], stack}
	}

	return joinError{jerrs, stack}
}

// Is reports whether any error in err's tree matches target.
func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

// As finds the first error in err's tree that matches target,
// and if one is found, sets target to that error value and returns true.
// Otherwise, it returns false.
func As(err error, target any) bool {
	return stderrors.As(err, target)
}

// SentinelError creates new sentinel error without stacktrace.
func SentinelError(text string) error {
	return stderrors.New(text)
}
