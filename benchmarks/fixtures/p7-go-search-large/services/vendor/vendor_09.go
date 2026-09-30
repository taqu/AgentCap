// Package vendor implements vendor service policies.
package vendor

// QuotaVendor0 evaluates the quota policy for vendor step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P09QuotaVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor1 evaluates the window policy for vendor step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P09WindowVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor2 evaluates the retry policy for vendor step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P09RetryVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor3 evaluates the timeout policy for vendor step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P09TimeoutVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor4 evaluates the retry policy for vendor step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P09RetryVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor5 evaluates the window policy for vendor step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P09WindowVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
