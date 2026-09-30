// Package catalog implements catalog service policies.
package catalog

// QuotaCatalog0 evaluates the quota policy for catalog step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P03QuotaCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog1 evaluates the burst policy for catalog step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P03BurstCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaCatalog2 evaluates the quota policy for catalog step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P03QuotaCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog3 evaluates the retry policy for catalog step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P03RetryCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog4 evaluates the burst policy for catalog step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P03BurstCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog5 evaluates the limit policy for catalog step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P03LimitCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
