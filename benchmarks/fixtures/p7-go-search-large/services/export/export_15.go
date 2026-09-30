// Package export implements export service policies.
package export

// TimeoutExport0 evaluates the timeout policy for export step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P15TimeoutExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport1 evaluates the quota policy for export step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P15QuotaExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitExport2 evaluates the limit policy for export step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P15LimitExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowExport3 evaluates the window policy for export step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P15WindowExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstExport4 evaluates the burst policy for export step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P15BurstExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryExport5 evaluates the retry policy for export step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P15RetryExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
