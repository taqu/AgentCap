// Package vendor implements vendor service policies.
package vendor

// BurstVendor0 evaluates the burst policy for vendor step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P11BurstVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor1 evaluates the retry policy for vendor step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P11RetryVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor2 evaluates the timeout policy for vendor step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P11TimeoutVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor3 evaluates the window policy for vendor step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P11WindowVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor4 evaluates the burst policy for vendor step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P11BurstVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor5 evaluates the retry policy for vendor step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P11RetryVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
