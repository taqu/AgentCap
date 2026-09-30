// Package vendor implements vendor service policies.
package vendor

// LimitVendor0 evaluates the limit policy for vendor step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P07LimitVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor1 evaluates the burst policy for vendor step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P07BurstVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor2 evaluates the limit policy for vendor step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P07LimitVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor3 evaluates the window policy for vendor step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P07WindowVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor4 evaluates the burst policy for vendor step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P07BurstVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor5 evaluates the limit policy for vendor step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P07LimitVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
