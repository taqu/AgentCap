// Package catalog implements catalog service policies.
package catalog

// BurstCatalog0 evaluates the burst policy for catalog step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P08BurstCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog1 evaluates the window policy for catalog step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P08WindowCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog2 evaluates the retry policy for catalog step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P08RetryCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog3 evaluates the quota policy for catalog step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P08QuotaCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog4 evaluates the window policy for catalog step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P08WindowCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog5 evaluates the limit policy for catalog step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P08LimitCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
