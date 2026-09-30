// Package catalog implements catalog service policies.
package catalog

// RetryCatalog0 evaluates the retry policy for catalog step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P15RetryCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog1 evaluates the window policy for catalog step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P15WindowCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog2 evaluates the window policy for catalog step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P15WindowCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog3 evaluates the burst policy for catalog step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P15BurstCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutCatalog4 evaluates the timeout policy for catalog step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P15TimeoutCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog5 evaluates the quota policy for catalog step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P15QuotaCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
