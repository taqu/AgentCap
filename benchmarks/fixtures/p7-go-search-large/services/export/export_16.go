// Package export implements export service policies.
package export

// WindowExport0 evaluates the window policy for export step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P16WindowExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryExport1 evaluates the retry policy for export step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P16RetryExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport2 evaluates the timeout policy for export step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P16TimeoutExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryExport3 evaluates the retry policy for export step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P16RetryExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowExport4 evaluates the window policy for export step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P16WindowExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstExport5 evaluates the burst policy for export step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P16BurstExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
