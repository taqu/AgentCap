// Package catalog implements catalog service policies.
package catalog

// QuotaCatalog0 evaluates the quota policy for catalog step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P16QuotaCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog1 evaluates the limit policy for catalog step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P16LimitCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog2 evaluates the retry policy for catalog step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P16RetryCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog3 evaluates the limit policy for catalog step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P16LimitCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog4 evaluates the retry policy for catalog step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P16RetryCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog5 evaluates the quota policy for catalog step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P16QuotaCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
