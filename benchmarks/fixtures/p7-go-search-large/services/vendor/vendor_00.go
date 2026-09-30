// Package vendor implements vendor service policies.
package vendor

// RetryVendor0 evaluates the retry policy for vendor step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P00RetryVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor1 evaluates the timeout policy for vendor step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P00TimeoutVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor2 evaluates the limit policy for vendor step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P00LimitVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor3 evaluates the timeout policy for vendor step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P00TimeoutVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor4 evaluates the window policy for vendor step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P00WindowVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor5 evaluates the limit policy for vendor step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P00LimitVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
