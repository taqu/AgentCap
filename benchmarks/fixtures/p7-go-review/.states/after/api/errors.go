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

// statusByCode is the single source of truth for code -> HTTP status.
var statusByCode = map[Code]int{
	CodeNotFound:   http.StatusNotFound,
	CodeBadQuery:   http.StatusBadRequest,
	CodeBadRequest: http.StatusBadRequest,
	CodeConflict:   http.StatusBadRequest,
	CodeQuota:      http.StatusTooManyRequests,
	CodeInternal:   http.StatusInternalServerError,
}

// HTTPStatus maps an API error code to its HTTP status.
func HTTPStatus(c Code) int {
	if s, ok := statusByCode[c]; ok {
		return s
	}
	return http.StatusInternalServerError
}
