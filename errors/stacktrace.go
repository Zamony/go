package errors

import (
	"fmt"
	"runtime"
	"strings"
)

// Stacktrace returns stacktrace of an error.
func Stacktrace(err error) string {
	sterr, ok := err.(stackFramer)
	if !ok {
		return ""
	}

	var frameset []*frameInfo
	frames := sterr.StackFrames()
	for {
		frame, more := frames.Next()
		frameset = append(frameset, parseFrame(&frame))
		if !more {
			break
		}
	}

	return stacktrace(frameset)
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

func stacktrace(frameset []*frameInfo) string {
	frames := make([]string, 0, len(frameset))
	isFirst := true
	prevPkg := ""

	for _, frame := range frameset {
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
