// Package vendor implements vendor service policies.
package vendor

// TimeoutVendor0 evaluates the timeout policy for vendor step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P14TimeoutVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor1 evaluates the timeout policy for vendor step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P14TimeoutVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor2 evaluates the burst policy for vendor step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P14BurstVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor3 evaluates the burst policy for vendor step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P14BurstVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor4 evaluates the burst policy for vendor step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P14BurstVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor5 evaluates the burst policy for vendor step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P14BurstVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
