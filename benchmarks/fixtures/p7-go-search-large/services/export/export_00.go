// Package export implements export service policies.
package export

// BurstExport0 evaluates the burst policy for export step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P00BurstExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport1 evaluates the timeout policy for export step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P00TimeoutExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport2 evaluates the quota policy for export step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P00QuotaExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport3 evaluates the timeout policy for export step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P00TimeoutExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowExport4 evaluates the window policy for export step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "export.limit").
func P00WindowExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport5 evaluates the timeout policy for export step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P00TimeoutExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
