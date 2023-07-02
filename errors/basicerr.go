package errors

import (
	stderrors "errors"
	"fmt"
	"io"
	"runtime"
)

type basicError struct {
	err    error
	frames *runtime.Frames
}

func (b *basicError) Unwrap() error {
	return b.err
}

func (b *basicError) StackFrames() *runtime.Frames {
	return b.frames
}

func (b *basicError) Error() string {
	return b.err.Error()
}

func (b *basicError) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			if trace := Stacktrace(b); trace != "" {
				_, _ = fmt.Fprintf(state, "%s: %s", b.Error(), trace)
				return
			}
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(state, b.Error())
	}
}

// New creates an error with message and a stacktrace.
func New(text string) error {
	return newError(frames(), text)
}

// Newf creates an error with formatted message and a stacktrace.
func Newf(format string, a ...any) error {
	return newError(frames(), fmt.Sprintf(format, a...))
}

// Wrap adds a stacktrace to the error if it doesn't have one.
func Wrap(err error) error {
	if err == nil {
		return nil
	}
	if stackErr, ok := err.(stackFramer); ok {
		return &basicError{err, stackErr.StackFrames()}
	}
	return &basicError{err, frames()}
}

// Wrap annotates an error with a message.
// Also adds a stacktrace to the error if it doesn't have one.
func Wrapf(err error, format string, a ...any) error {
	msg := fmt.Sprintf(format, a...)
	if stackErr, ok := err.(stackFramer); ok {
		return wrapError(stackErr.StackFrames(), err, msg)
	}
	return wrapError(frames(), err, msg)
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

func newError(frames *runtime.Frames, text string) error {
	return &basicError{stderrors.New(text), frames}
}

type stackFramer interface {
	StackFrames() *runtime.Frames
}

func wrapError(frames *runtime.Frames, err error, msg string) error {
	if err == nil {
		return nil
	}

	newErr := fmt.Errorf("%s: %w", msg, err)
	if stackErr, ok := err.(stackFramer); ok {
		return &basicError{newErr, stackErr.StackFrames()}
	}

	return &basicError{newErr, frames}
}

func frames() *runtime.Frames {
	const (
		maxDepth = 32
		skip     = 3 // +1 Callers, +1 frames, +1 New
	)
	var pcs [maxDepth]uintptr
	depth := runtime.Callers(skip, pcs[:])
	return runtime.CallersFrames(pcs[:depth])
}
