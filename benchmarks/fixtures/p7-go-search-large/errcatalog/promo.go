package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"promo.missing": api.CodeNotFound,
	"promo.invalid": api.CodeBadRequest,
	"promo.locked": api.CodeConflict,
	"promo.expired": api.CodeBadRequest,
	"promo.limit": api.CodeQuota,
	"promo.timeout": api.CodeInternal,
	"promo.denied": api.CodeBadRequest,
	"promo.stale": api.CodeConflict,
	})
}
