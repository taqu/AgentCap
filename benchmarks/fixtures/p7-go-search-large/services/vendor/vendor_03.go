// Package vendor implements vendor service policies.
package vendor

// BurstVendor0 evaluates the burst policy for vendor step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P03BurstVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor1 evaluates the burst policy for vendor step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P03BurstVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor2 evaluates the limit policy for vendor step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P03LimitVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor3 evaluates the window policy for vendor step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P03WindowVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor4 evaluates the window policy for vendor step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P03WindowVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor5 evaluates the burst policy for vendor step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P03BurstVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
