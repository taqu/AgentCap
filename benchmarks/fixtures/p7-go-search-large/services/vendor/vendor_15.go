// Package vendor implements vendor service policies.
package vendor

// QuotaVendor0 evaluates the quota policy for vendor step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P15QuotaVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor1 evaluates the timeout policy for vendor step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P15TimeoutVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor2 evaluates the burst policy for vendor step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P15BurstVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor3 evaluates the retry policy for vendor step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P15RetryVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor4 evaluates the timeout policy for vendor step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P15TimeoutVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaVendor5 evaluates the quota policy for vendor step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P15QuotaVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
