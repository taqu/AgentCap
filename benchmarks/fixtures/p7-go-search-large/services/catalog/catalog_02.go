// Package catalog implements catalog service policies.
package catalog

// WindowCatalog0 evaluates the window policy for catalog step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P02WindowCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog1 evaluates the quota policy for catalog step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P02QuotaCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog2 evaluates the quota policy for catalog step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P02QuotaCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog3 evaluates the window policy for catalog step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P02WindowCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutCatalog4 evaluates the timeout policy for catalog step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P02TimeoutCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog5 evaluates the window policy for catalog step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P02WindowCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
