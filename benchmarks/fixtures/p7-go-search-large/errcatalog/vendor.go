package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"vendor.missing": api.CodeNotFound,
	"vendor.invalid": api.CodeBadRequest,
	"vendor.locked": api.CodeConflict,
	"vendor.expired": api.CodeBadRequest,
	"vendor.limit": api.CodeQuota,
	"vendor.timeout": api.CodeInternal,
	"vendor.denied": api.CodeBadRequest,
	"vendor.stale": api.CodeConflict,
	})
}
