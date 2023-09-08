package bla_test

import (
	"errors"
	"ncp/bla"
	"net/http"
	"testing"

	tu "github.com/volvo-cars/lingon/pkg/testutil"
)

func TestErrors(t *testing.T) {
	type test struct {
		name     string
		err      error
		expected *bla.Errors
	}
	tests := []test{
		{
			name: "ordinary error",
			err:  errors.New("whatever"),
			expected: &bla.Errors{
				Errors: []*bla.Error{
					{
						Status:  http.StatusInternalServerError,
						Message: "whatever",
					},
				},
			},
		},
		{
			name: "error with status",
			err: &bla.Error{
				Status:  http.StatusBadRequest,
				Message: "whatever",
			},
			expected: &bla.Errors{
				Errors: []*bla.Error{
					{
						Status:  http.StatusBadRequest,
						Message: "whatever",
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := bla.FromError(test.err)
			tu.AssertEqual(t, test.expected, actual)
		})
	}
}
