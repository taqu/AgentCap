// Package catalog implements catalog service policies.
package catalog

// QuotaCatalog0 evaluates the quota policy for catalog step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P01QuotaCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutCatalog1 evaluates the timeout policy for catalog step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P01TimeoutCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog2 evaluates the burst policy for catalog step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P01BurstCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog3 evaluates the limit policy for catalog step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P01LimitCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutCatalog4 evaluates the timeout policy for catalog step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P01TimeoutCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog5 evaluates the burst policy for catalog step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P01BurstCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
