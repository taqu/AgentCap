// Package vendor implements vendor service policies.
package vendor

// TimeoutVendor0 evaluates the timeout policy for vendor step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P17TimeoutVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaVendor1 evaluates the quota policy for vendor step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P17QuotaVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor2 evaluates the window policy for vendor step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P17WindowVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor3 evaluates the window policy for vendor step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P17WindowVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaVendor4 evaluates the quota policy for vendor step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P17QuotaVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowVendor5 evaluates the window policy for vendor step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P17WindowVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
