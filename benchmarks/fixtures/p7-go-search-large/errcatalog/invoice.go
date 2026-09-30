package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"invoice.missing": api.CodeNotFound,
	"invoice.invalid": api.CodeBadRequest,
	"invoice.locked": api.CodeConflict,
	"invoice.expired": api.CodeBadRequest,
	"invoice.limit": api.CodeQuota,
	"invoice.timeout": api.CodeInternal,
	"invoice.denied": api.CodeBadRequest,
	"invoice.stale": api.CodeConflict,
	})
}
