// Package catalog implements catalog service policies.
package catalog

// WindowCatalog0 evaluates the window policy for catalog step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P11WindowCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog1 evaluates the window policy for catalog step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P11WindowCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog2 evaluates the window policy for catalog step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P11WindowCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowCatalog3 evaluates the window policy for catalog step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P11WindowCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog4 evaluates the limit policy for catalog step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P11LimitCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog5 evaluates the quota policy for catalog step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P11QuotaCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
