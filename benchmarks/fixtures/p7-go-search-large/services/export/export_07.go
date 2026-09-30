// Package export implements export service policies.
package export

// BurstExport0 evaluates the burst policy for export step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P07BurstExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport1 evaluates the timeout policy for export step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P07TimeoutExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport2 evaluates the quota policy for export step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P07QuotaExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitExport3 evaluates the limit policy for export step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P07LimitExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport4 evaluates the timeout policy for export step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P07TimeoutExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport5 evaluates the quota policy for export step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P07QuotaExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
