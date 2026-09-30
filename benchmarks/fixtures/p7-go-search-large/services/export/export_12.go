// Package export implements export service policies.
package export

// QuotaExport0 evaluates the quota policy for export step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P12QuotaExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport1 evaluates the quota policy for export step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P12QuotaExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitExport2 evaluates the limit policy for export step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P12LimitExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport3 evaluates the timeout policy for export step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P12TimeoutExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutExport4 evaluates the timeout policy for export step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P12TimeoutExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitExport5 evaluates the limit policy for export step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P12LimitExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
