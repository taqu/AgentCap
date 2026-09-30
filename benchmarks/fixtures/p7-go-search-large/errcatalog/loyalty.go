package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"loyalty.missing": api.CodeNotFound,
	"loyalty.invalid": api.CodeBadRequest,
	"loyalty.locked": api.CodeConflict,
	"loyalty.expired": api.CodeBadRequest,
	"loyalty.limit": api.CodeQuota,
	"loyalty.timeout": api.CodeInternal,
	"loyalty.denied": api.CodeBadRequest,
	"loyalty.stale": api.CodeConflict,
	})
}
