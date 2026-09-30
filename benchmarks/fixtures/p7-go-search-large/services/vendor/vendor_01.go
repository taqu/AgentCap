// Package vendor implements vendor service policies.
package vendor

// BurstVendor0 evaluates the burst policy for vendor step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P01BurstVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor1 evaluates the burst policy for vendor step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P01BurstVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor2 evaluates the window policy for vendor step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P01WindowVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor3 evaluates the window policy for vendor step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P01WindowVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryVendor4 evaluates the retry policy for vendor step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P01RetryVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor5 evaluates the window policy for vendor step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P01WindowVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
