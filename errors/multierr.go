package errors

import (
	"fmt"
	"io"
	"strings"
)

type multiError []error

func (e multiError) Unwrap() []error {
	return e
}

func (e multiError) toString(format func(error) string) string {
	errs := make([]string, len(e))
	for i := range e {
		errs[i] = format(e[i])
	}

	return strings.Join(errs, "; ")
}

func (e multiError) Error() string {
	return e.toString(errorStr)
}

func (e multiError) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			_, _ = io.WriteString(state, e.toString(errorVerboseStr))
			return
		}
		fallthrough
	case 's', 'q':
		_, _ = io.WriteString(state, e.Error())
	}
}

// Combine joins multiple errors into one.
func Combine(errs ...error) error {
	result := make(multiError, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			result = append(result, err)
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func errorStr(err error) string {
	return err.Error()
}

func errorVerboseStr(err error) string {
	return fmt.Sprintf("%+v", err)
}