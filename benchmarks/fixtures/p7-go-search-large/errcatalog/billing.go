package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"billing.missing": api.CodeNotFound,
	"billing.invalid": api.CodeBadRequest,
	"billing.locked": api.CodeConflict,
	"billing.expired": api.CodeBadRequest,
	"billing.limit": api.CodeQuota,
	"billing.timeout": api.CodeInternal,
	"billing.denied": api.CodeBadRequest,
	"billing.stale": api.CodeConflict,
	})
}
