// Package vendor implements vendor service policies.
package vendor

// WindowVendor0 evaluates the window policy for vendor step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P10WindowVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaVendor1 evaluates the quota policy for vendor step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P10QuotaVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor2 evaluates the timeout policy for vendor step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P10TimeoutVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaVendor3 evaluates the quota policy for vendor step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P10QuotaVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor4 evaluates the retry policy for vendor step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P10RetryVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor5 evaluates the burst policy for vendor step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P10BurstVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
