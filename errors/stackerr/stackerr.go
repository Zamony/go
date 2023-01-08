package stackerr

import (
	"fmt"
	"runtime"
	"strings"
)

type stackFramer interface {
	StackFrames() *runtime.Frames
}

func MarshalCompact(err error) string {
	sterr, ok := err.(stackFramer)
	if !ok {
		return ""
	}

	var frameset frameSet
	frames := sterr.StackFrames()
	for {
		frame, more := frames.Next()
		frameset = append(frameset, parseFrame(&frame))
		if !more {
			break
		}
	}

	return frameset.Build()
}

func MarshalCompactZerolog(err error) interface{} {
	if v := MarshalCompact(err); v != "" {
		return v
	}
	return nil
}

type frameSet []*frameInfo

func (fs frameSet) Build() string {
	frames := make([]string, 0, len(fs))
	isFirst := true
	prevPkg := ""

	for _, frame := range fs {
		if isFirst {
			f := fmt.Sprintf("%s.%s:%d", frame.Package, frame.Func, frame.Line)
			frames = append(frames, f)
			isFirst = false
		} else if frame.Package != prevPkg {
			f := fmt.Sprintf("%s.%s", frame.Package, frame.Func)
			frames = append(frames, f)
		} else {
			frames = append(frames, frame.Func)
		}

		prevPkg = frame.Package
	}

	return strings.Join(frames, "/")
}

type frameInfo struct {
	Package, Func string
	Line          int
}

func parseFrame(frame *runtime.Frame) *frameInfo {
	fun := frame.Function
	slashIdx := strings.LastIndexByte(fun, '/')
	if slashIdx >= 0 {
		fun = fun[slashIdx+1:]
	}

	pkg, fun, _ := strings.Cut(fun, ".")
	return &frameInfo{pkg, fun, frame.Line}
}
