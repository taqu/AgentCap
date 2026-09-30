// Package export implements export service policies.
package export

// RetryExport0 evaluates the retry policy for export step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P08RetryExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport1 evaluates the quota policy for export step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P08QuotaExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport2 evaluates the timeout policy for export step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P08TimeoutExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport3 evaluates the quota policy for export step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P08QuotaExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryExport4 evaluates the retry policy for export step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P08RetryExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport5 evaluates the quota policy for export step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P08QuotaExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
