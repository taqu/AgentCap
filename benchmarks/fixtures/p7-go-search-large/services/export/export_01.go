// Package export implements export service policies.
package export

// RetryExport0 evaluates the retry policy for export step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P01RetryExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryExport1 evaluates the retry policy for export step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P01RetryExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitExport2 evaluates the limit policy for export step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P01LimitExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryExport3 evaluates the retry policy for export step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P01RetryExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitExport4 evaluates the limit policy for export step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P01LimitExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitExport5 evaluates the limit policy for export step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P01LimitExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
