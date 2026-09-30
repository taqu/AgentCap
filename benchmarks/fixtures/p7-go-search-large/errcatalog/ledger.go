package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"ledger.missing": api.CodeNotFound,
	"ledger.invalid": api.CodeBadRequest,
	"ledger.locked": api.CodeConflict,
	"ledger.expired": api.CodeBadRequest,
	"ledger.limit": api.CodeQuota,
	"ledger.timeout": api.CodeInternal,
	"ledger.denied": api.CodeBadRequest,
	"ledger.stale": api.CodeConflict,
	})
}
