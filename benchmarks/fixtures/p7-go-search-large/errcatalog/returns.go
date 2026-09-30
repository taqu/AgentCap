package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"returns.missing": api.CodeNotFound,
	"returns.invalid": api.CodeBadRequest,
	"returns.locked": api.CodeConflict,
	"returns.expired": api.CodeBadRequest,
	"returns.limit": api.CodeQuota,
	"returns.timeout": api.CodeInternal,
	"returns.denied": api.CodeBadRequest,
	"returns.stale": api.CodeConflict,
	})
}
