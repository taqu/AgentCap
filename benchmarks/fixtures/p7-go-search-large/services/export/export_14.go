// Package export implements export service policies.
package export

// TimeoutExport0 evaluates the timeout policy for export step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "export.limit").
func P14TimeoutExport0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitExport1 evaluates the limit policy for export step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "export.limit").
func P14LimitExport1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstExport2 evaluates the burst policy for export step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P14BurstExport2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryExport3 evaluates the retry policy for export step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "export.limit").
func P14RetryExport3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstExport4 evaluates the burst policy for export step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "export.limit").
func P14BurstExport4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaExport5 evaluates the quota policy for export step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "export.limit").
func P14QuotaExport5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
