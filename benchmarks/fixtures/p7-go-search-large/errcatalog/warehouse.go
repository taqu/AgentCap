package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"warehouse.missing": api.CodeNotFound,
	"warehouse.invalid": api.CodeBadRequest,
	"warehouse.locked": api.CodeConflict,
	"warehouse.expired": api.CodeBadRequest,
	"warehouse.limit": api.CodeInternal,
	"warehouse.timeout": api.CodeInternal,
	"warehouse.denied": api.CodeBadRequest,
	"warehouse.stale": api.CodeConflict,
	})
}
