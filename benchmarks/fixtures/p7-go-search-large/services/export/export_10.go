// Package export implements export service policies.
package export

// WindowExport0 evaluates the window policy for export step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P10WindowExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryExport1 evaluates the retry policy for export step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P10RetryExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowExport2 evaluates the window policy for export step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P10WindowExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport3 evaluates the quota policy for export step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P10QuotaExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstExport4 evaluates the burst policy for export step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P10BurstExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitExport5 evaluates the limit policy for export step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P10LimitExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
