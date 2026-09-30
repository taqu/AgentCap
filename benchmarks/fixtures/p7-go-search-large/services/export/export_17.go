// Package export implements export service policies.
package export

// WindowExport0 evaluates the window policy for export step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P17WindowExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowExport1 evaluates the window policy for export step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P17WindowExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstExport2 evaluates the burst policy for export step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P17BurstExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport3 evaluates the timeout policy for export step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P17TimeoutExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstExport4 evaluates the burst policy for export step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P17BurstExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport5 evaluates the timeout policy for export step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P17TimeoutExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
