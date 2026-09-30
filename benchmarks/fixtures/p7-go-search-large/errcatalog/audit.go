package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"audit.missing": api.CodeNotFound,
	"audit.invalid": api.CodeBadRequest,
	"audit.locked": api.CodeConflict,
	"audit.expired": api.CodeBadRequest,
	"audit.limit": api.CodeQuota,
	"audit.timeout": api.CodeInternal,
	"audit.denied": api.CodeBadRequest,
	"audit.stale": api.CodeConflict,
	})
}
