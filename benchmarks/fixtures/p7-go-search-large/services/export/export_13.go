// Package export implements export service policies.
package export

// LimitExport0 evaluates the limit policy for export step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P13LimitExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport1 evaluates the timeout policy for export step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P13TimeoutExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstExport2 evaluates the burst policy for export step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P13BurstExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport3 evaluates the timeout policy for export step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P13TimeoutExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryExport4 evaluates the retry policy for export step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P13RetryExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstExport5 evaluates the burst policy for export step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P13BurstExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
