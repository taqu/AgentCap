package errcatalog

import "inventory/api"

func init() {
	register(map[string]api.Code{
	"ingest.missing": api.CodeNotFound,
	"ingest.invalid": api.CodeBadRequest,
	"ingest.locked": api.CodeConflict,
	"ingest.expired": api.CodeBadRequest,
	"ingest.limit": api.CodeQuota,
	"ingest.timeout": api.CodeInternal,
	"ingest.denied": api.CodeBadRequest,
	"ingest.stale": api.CodeConflict,
	})
}
