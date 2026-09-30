// Package vendor implements vendor service policies.
package vendor

// WindowVendor0 evaluates the window policy for vendor step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P12WindowVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor1 evaluates the burst policy for vendor step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P12BurstVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor2 evaluates the timeout policy for vendor step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P12TimeoutVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor3 evaluates the burst policy for vendor step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P12BurstVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor4 evaluates the burst policy for vendor step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P12BurstVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor5 evaluates the timeout policy for vendor step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P12TimeoutVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
