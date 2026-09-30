// Package export implements export service policies.
package export

// WindowExport0 evaluates the window policy for export step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P03WindowExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryExport1 evaluates the retry policy for export step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P03RetryExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport2 evaluates the quota policy for export step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P03QuotaExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstExport3 evaluates the burst policy for export step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P03BurstExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport4 evaluates the quota policy for export step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P03QuotaExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport5 evaluates the timeout policy for export step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P03TimeoutExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
