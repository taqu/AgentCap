// Package catalog implements catalog service policies.
package catalog

// LimitCatalog0 evaluates the limit policy for catalog step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P17LimitCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog1 evaluates the burst policy for catalog step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P17BurstCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog2 evaluates the quota policy for catalog step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P17QuotaCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog3 evaluates the quota policy for catalog step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P17QuotaCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog4 evaluates the quota policy for catalog step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P17QuotaCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog5 evaluates the retry policy for catalog step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P17RetryCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
