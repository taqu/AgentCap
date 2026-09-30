// Package errcatalog maps internal domain error keys ("domain.kind")
// to public API error codes.
package errcatalog

import "inventory/api"

var table = map[string]api.Code{}

func register(m map[string]api.Code) {
	for k, v := range m {
		table[k] = v
	}
}

// Lookup returns the API code for key, or CodeInternal when unknown.
func Lookup(key string) api.Code {
	if c, ok := table[key]; ok {
		return c
	}
	return api.CodeInternal
}

// Status returns the HTTP status for a domain error key.
func Status(key string) int { return api.HTTPStatus(Lookup(key)) }
