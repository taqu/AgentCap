// Package export implements export service policies.
package export

// RetryExport0 evaluates the retry policy for export step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P05RetryExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport1 evaluates the timeout policy for export step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P05TimeoutExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport2 evaluates the timeout policy for export step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P05TimeoutExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryExport3 evaluates the retry policy for export step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P05RetryExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryExport4 evaluates the retry policy for export step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P05RetryExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitExport5 evaluates the limit policy for export step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P05LimitExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
