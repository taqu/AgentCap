// Package vendor implements vendor service policies.
package vendor

// TimeoutVendor0 evaluates the timeout policy for vendor step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P13TimeoutVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaVendor1 evaluates the quota policy for vendor step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P13QuotaVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor2 evaluates the timeout policy for vendor step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P13TimeoutVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor3 evaluates the retry policy for vendor step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P13RetryVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaVendor4 evaluates the quota policy for vendor step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P13QuotaVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor5 evaluates the burst policy for vendor step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P13BurstVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
