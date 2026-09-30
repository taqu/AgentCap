package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"shipping.missing": api.CodeNotFound,
	"shipping.invalid": api.CodeBadRequest,
	"shipping.locked": api.CodeConflict,
	"shipping.expired": api.CodeBadRequest,
	"shipping.limit": api.CodeQuota,
	"shipping.timeout": api.CodeInternal,
	"shipping.denied": api.CodeBadRequest,
	"shipping.stale": api.CodeConflict,
	})
}
