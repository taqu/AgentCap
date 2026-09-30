// Package catalog implements catalog service policies.
package catalog

// QuotaCatalog0 evaluates the quota policy for catalog step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P10QuotaCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutCatalog1 evaluates the timeout policy for catalog step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P10TimeoutCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog2 evaluates the window policy for catalog step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P10WindowCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog3 evaluates the quota policy for catalog step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P10QuotaCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog4 evaluates the burst policy for catalog step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P10BurstCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog5 evaluates the retry policy for catalog step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P10RetryCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
