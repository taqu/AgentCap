package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"export.missing": api.CodeNotFound,
	"export.invalid": api.CodeBadRequest,
	"export.locked": api.CodeConflict,
	"export.expired": api.CodeBadRequest,
	"export.limit": api.CodeQuota,
	"export.timeout": api.CodeInternal,
	"export.denied": api.CodeBadRequest,
	"export.stale": api.CodeConflict,
	})
}
