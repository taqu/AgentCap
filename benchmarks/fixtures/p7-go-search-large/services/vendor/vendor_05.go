// Package vendor implements vendor service policies.
package vendor

// WindowVendor0 evaluates the window policy for vendor step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P05WindowVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor1 evaluates the retry policy for vendor step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P05RetryVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor2 evaluates the retry policy for vendor step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P05RetryVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor3 evaluates the retry policy for vendor step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P05RetryVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor4 evaluates the retry policy for vendor step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P05RetryVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor5 evaluates the limit policy for vendor step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P05LimitVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
