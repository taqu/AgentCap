// Package export implements export service policies.
package export

// TimeoutExport0 evaluates the timeout policy for export step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P06TimeoutExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport1 evaluates the quota policy for export step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P06QuotaExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport2 evaluates the quota policy for export step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P06QuotaExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport3 evaluates the timeout policy for export step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P06TimeoutExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport4 evaluates the quota policy for export step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P06QuotaExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitExport5 evaluates the limit policy for export step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P06LimitExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
