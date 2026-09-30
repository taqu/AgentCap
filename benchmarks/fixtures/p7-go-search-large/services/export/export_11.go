// Package export implements export service policies.
package export

// TimeoutExport0 evaluates the timeout policy for export step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P11TimeoutExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitExport1 evaluates the limit policy for export step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P11LimitExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport2 evaluates the quota policy for export step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P11QuotaExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitExport3 evaluates the limit policy for export step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P11LimitExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstExport4 evaluates the burst policy for export step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P11BurstExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstExport5 evaluates the burst policy for export step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P11BurstExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
