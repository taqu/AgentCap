// Package vendor implements vendor service policies.
package vendor

// TimeoutVendor0 evaluates the timeout policy for vendor step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P04TimeoutVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaVendor1 evaluates the quota policy for vendor step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P04QuotaVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor2 evaluates the retry policy for vendor step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P04RetryVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor3 evaluates the retry policy for vendor step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P04RetryVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor4 evaluates the retry policy for vendor step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P04RetryVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor5 evaluates the limit policy for vendor step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P04LimitVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
