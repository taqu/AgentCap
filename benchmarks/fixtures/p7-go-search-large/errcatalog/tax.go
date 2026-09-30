package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"tax.missing": api.CodeNotFound,
	"tax.invalid": api.CodeBadRequest,
	"tax.locked": api.CodeConflict,
	"tax.expired": api.CodeBadRequest,
	"tax.limit": api.CodeQuota,
	"tax.timeout": api.CodeInternal,
	"tax.denied": api.CodeBadRequest,
	"tax.stale": api.CodeConflict,
	})
}
