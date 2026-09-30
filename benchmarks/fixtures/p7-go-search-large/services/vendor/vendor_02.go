// Package vendor implements vendor service policies.
package vendor

// LimitVendor0 evaluates the limit policy for vendor step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P02LimitVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor1 evaluates the window policy for vendor step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P02WindowVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor2 evaluates the retry policy for vendor step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P02RetryVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor3 evaluates the limit policy for vendor step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P02LimitVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaVendor4 evaluates the quota policy for vendor step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P02QuotaVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor5 evaluates the limit policy for vendor step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P02LimitVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
