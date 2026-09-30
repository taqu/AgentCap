package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"session.missing": api.CodeNotFound,
	"session.invalid": api.CodeBadRequest,
	"session.locked": api.CodeConflict,
	"session.expired": api.CodeBadRequest,
	"session.limit": api.CodeQuota,
	"session.timeout": api.CodeInternal,
	"session.denied": api.CodeBadRequest,
	"session.stale": api.CodeConflict,
	})
}
