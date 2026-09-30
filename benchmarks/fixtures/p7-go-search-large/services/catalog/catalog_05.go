// Package catalog implements catalog service policies.
package catalog

// QuotaCatalog0 evaluates the quota policy for catalog step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P05QuotaCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog1 evaluates the quota policy for catalog step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P05QuotaCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog2 evaluates the quota policy for catalog step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P05QuotaCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog3 evaluates the window policy for catalog step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P05WindowCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog4 evaluates the limit policy for catalog step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P05LimitCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog5 evaluates the burst policy for catalog step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P05BurstCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
