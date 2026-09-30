// Package catalog implements catalog service policies.
package catalog

// LimitCatalog0 evaluates the limit policy for catalog step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P12LimitCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog1 evaluates the quota policy for catalog step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P12QuotaCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog2 evaluates the retry policy for catalog step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P12RetryCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog3 evaluates the burst policy for catalog step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P12BurstCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog4 evaluates the window policy for catalog step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P12WindowCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutCatalog5 evaluates the timeout policy for catalog step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P12TimeoutCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
