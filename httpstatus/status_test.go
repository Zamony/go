package httpstatus_test

import (
	"net/http"
	"testing"

	"github.com/Zamony/httpstatus"
)

func TestStatus(t *testing.T) {
	testCases := []struct {
		Code int
		Want httpstatus.Type
	}{
		{
			Code: http.StatusContinue,
			Want: httpstatus.Information,
		},
		{
			Code: http.StatusOK,
			Want: httpstatus.Success,
		},
		{
			Code: http.StatusMovedPermanently,
			Want: httpstatus.Redirection,
		},
		{
			Code: http.StatusNotFound,
			Want: httpstatus.ClientError,
		},
		{
			Code: http.StatusInternalServerError,
			Want: httpstatus.ServerError,
		},
	}

	for _, tc := range testCases {
		status := httpstatus.From(tc.Code)
		if status != tc.Want {
			t.Fatalf(
				"For [code=%d] wanted [status_type=%d], got [status_type=%d]",
				tc.Code, tc.Want, status,
			)
		}
	}
}
