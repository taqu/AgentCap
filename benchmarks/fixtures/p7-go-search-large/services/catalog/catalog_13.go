// Package catalog implements catalog service policies.
package catalog

// WindowCatalog0 evaluates the window policy for catalog step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P13WindowCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog1 evaluates the retry policy for catalog step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P13RetryCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog2 evaluates the window policy for catalog step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P13WindowCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog3 evaluates the quota policy for catalog step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P13QuotaCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutCatalog4 evaluates the timeout policy for catalog step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P13TimeoutCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog5 evaluates the limit policy for catalog step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P13LimitCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
