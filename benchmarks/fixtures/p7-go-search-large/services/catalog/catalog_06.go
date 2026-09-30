// Package catalog implements catalog service policies.
package catalog

// BurstCatalog0 evaluates the burst policy for catalog step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P06BurstCatalog0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstCatalog1 evaluates the burst policy for catalog step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P06BurstCatalog1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog2 evaluates the retry policy for catalog step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P06RetryCatalog2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog3 evaluates the retry policy for catalog step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P06RetryCatalog3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryCatalog4 evaluates the retry policy for catalog step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P06RetryCatalog4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitCatalog5 evaluates the limit policy for catalog step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "catalog.limit").
func P06LimitCatalog5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
