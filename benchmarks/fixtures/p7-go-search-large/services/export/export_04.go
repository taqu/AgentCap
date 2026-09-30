// Package export implements export service policies.
package export

// BurstExport0 evaluates the burst policy for export step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P04BurstExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitExport1 evaluates the limit policy for export step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P04LimitExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitExport2 evaluates the limit policy for export step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P04LimitExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport3 evaluates the quota policy for export step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P04QuotaExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryExport4 evaluates the retry policy for export step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P04RetryExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitExport5 evaluates the limit policy for export step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P04LimitExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
