// Package multierror contains a Multierror for combining multiple errors and compiling with errors.Is
package multierror

import (
	"errors"
	"strings"
)

// MultiError is a collection of errors combined and treated as one
type MultiError []error

// Error returns combined error.Error output by \n
func (me MultiError) Error() string {
	builder := strings.Builder{}
	length := len(me) - 1
	for i, e := range me {
		if i == length {
			builder.WriteString(e.Error())
		} else {
			builder.WriteString(e.Error() + ", ")
		}
	}
	return builder.String()
}

// Is a target error, complies with errors.Is
func (me MultiError) Is(target error) bool {
	for _, e := range me {
		if errors.Is(e, target) {
			return true
		}
	}
	return false
}

// Append newErr to oldErr error returning combined error, appending newErr to old if old is a MultiError or creating a new MultiError if not. does not append nil errors, see source
func Append(oldErr error, newErr error) error {
	if oldErr == nil {
		return newErr
	}
	if newErr == nil {
		return oldErr
	}
	switch t := oldErr.(type) {
	case MultiError:
		return append(t, newErr)
	default:
		return MultiError{oldErr, newErr}
	}
}
