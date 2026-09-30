package api

import "net/http"

// Code is a stable machine-readable API error code.
type Code string

const (
	CodeNotFound   Code = "not_found"
	CodeBadQuery   Code = "bad_query"
	CodeConflict   Code = "conflict"
	CodeQuota      Code = "quota_exceeded"
	CodeInternal   Code = "internal"
	CodeBadRequest Code = "bad_request"
)

// HTTPStatus maps an API error code to its HTTP status.
func HTTPStatus(c Code) int {
	switch c {
	case CodeNotFound:
		return http.StatusNotFound
	case CodeBadQuery, CodeBadRequest:
		return http.StatusBadRequest
	case CodeConflict:
		return http.StatusConflict
	case CodeQuota:
		return http.StatusTooManyRequests
	}
	return http.StatusInternalServerError
}
