// Package catalog implements catalog service policies.
package catalog

// RetryCatalog0 evaluates the retry policy for catalog step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P14RetryCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog1 evaluates the limit policy for catalog step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P14LimitCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog2 evaluates the retry policy for catalog step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P14RetryCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog3 evaluates the window policy for catalog step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P14WindowCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog4 evaluates the window policy for catalog step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P14WindowCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog5 evaluates the quota policy for catalog step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P14QuotaCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
