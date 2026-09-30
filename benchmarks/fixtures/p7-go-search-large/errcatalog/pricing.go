package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"pricing.missing": api.CodeNotFound,
	"pricing.invalid": api.CodeBadRequest,
	"pricing.locked": api.CodeConflict,
	"pricing.expired": api.CodeBadRequest,
	"pricing.limit": api.CodeQuota,
	"pricing.timeout": api.CodeInternal,
	"pricing.denied": api.CodeBadRequest,
	"pricing.stale": api.CodeConflict,
	})
}
