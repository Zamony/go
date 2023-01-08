package errors

import (
	stderrors "errors"
	"fmt"
	"runtime"
)

type ConstantError string

func (e ConstantError) Error() string { return string(e) }

type base struct {
	err    error
	frames *runtime.Frames
}

func (b *base) Error() string {
	return b.err.Error()
}

func (b *base) Unwrap() error {
	return b.err
}

func (b *base) StackFrames() *runtime.Frames {
	return b.frames
}

func New(text string) error {
	return newError(frames(), text)
}

func Newf(format string, a ...interface{}) error {
	return newError(frames(), fmt.Sprintf(format, a...))
}

func From(err error) error {
	if err == nil {
		return nil
	}
	if stackErr, ok := err.(stackFramer); ok {
		return &base{err, stackErr.StackFrames()}
	}
	return &base{err, frames()}
}

func Wrap(err error, msg string) error {
	if stackErr, ok := err.(stackFramer); ok {
		return wrapError(stackErr.StackFrames(), err, msg)
	}
	return wrapError(frames(), err, msg)
}

func Wrapf(err error, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	if stackErr, ok := err.(stackFramer); ok {
		return wrapError(stackErr.StackFrames(), err, msg)
	}
	return wrapError(frames(), err, msg)
}

func Is(err, target error) bool {
	return stderrors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return stderrors.As(err, target)
}

func newError(frames *runtime.Frames, text string) error {
	return &base{stderrors.New(text), frames}
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
		return &base{newErr, stackErr.StackFrames()}
	}

	return &base{newErr, frames}
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
