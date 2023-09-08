package bla

import (
	"fmt"
	"net/http"
	"strings"
)

var _ error = (*Errors)(nil)

type Errors struct {
	Errors []*Error `json:"errors,omitempty"`
}

func (e *Errors) Error() string {
	ss := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		ss[i] = err.Error()
	}
	return strings.Join(ss, ",")
}

func (e *Errors) HasError() bool {
	return len(e.Errors) > 0
}

func (e *Errors) AddError(err *Error) {
	e.Errors = append(e.Errors, err)
}

var _ error = (*Error)(nil)

type Error struct {
	// Status is the HTTP status code applicable to this problem.
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (status %d)", e.Message, e.Status)
}

func FromError(err error) *Errors {
	switch val := err.(type) {
	case *Errors:
		return val
	case *Error:
		// Error should be wrapped in Errors.
		return &Errors{
			Errors: []*Error{val},
		}
	default:
		// Any other error should be wrapped in Errors.
		return FromError(&Error{
			Status:  http.StatusInternalServerError,
			Message: val.Error(),
		})
	}
}
