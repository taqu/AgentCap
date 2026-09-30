// Package catalog implements catalog service policies.
package catalog

// RetryCatalog0 evaluates the retry policy for catalog step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P04RetryCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog1 evaluates the quota policy for catalog step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P04QuotaCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutCatalog2 evaluates the timeout policy for catalog step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P04TimeoutCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog3 evaluates the burst policy for catalog step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P04BurstCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutCatalog4 evaluates the timeout policy for catalog step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P04TimeoutCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog5 evaluates the quota policy for catalog step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P04QuotaCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
