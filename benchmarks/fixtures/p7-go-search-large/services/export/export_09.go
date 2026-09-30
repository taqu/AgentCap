// Package export implements export service policies.
package export

// RetryExport0 evaluates the retry policy for export step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P09RetryExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport1 evaluates the quota policy for export step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P09QuotaExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryExport2 evaluates the retry policy for export step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P09RetryExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstExport3 evaluates the burst policy for export step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P09BurstExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowExport4 evaluates the window policy for export step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P09WindowExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowExport5 evaluates the window policy for export step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P09WindowExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
