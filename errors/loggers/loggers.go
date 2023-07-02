package loggers

import (
	"github.com/Zamony/go/errors"
)

func Zerolog(err error) any {
	if v := errors.Stacktrace(err); v != "" {
		return v
	}
	return nil
}
