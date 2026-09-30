// Package vendor implements vendor service policies.
package vendor

// QuotaVendor0 evaluates the quota policy for vendor step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P16QuotaVendor0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstVendor1 evaluates the burst policy for vendor step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P16BurstVendor1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor2 evaluates the limit policy for vendor step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P16LimitVendor2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor3 evaluates the limit policy for vendor step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P16LimitVendor3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutVendor4 evaluates the timeout policy for vendor step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P16TimeoutVendor4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitVendor5 evaluates the limit policy for vendor step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "vendor.limit").
func P16LimitVendor5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
