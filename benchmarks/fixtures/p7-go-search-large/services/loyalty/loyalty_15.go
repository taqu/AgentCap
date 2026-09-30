// Package loyalty implements loyalty service policies.
package loyalty

// BurstLoyalty0 evaluates the burst policy for loyalty step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P15BurstLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty1 evaluates the quota policy for loyalty step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P15QuotaLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty2 evaluates the limit policy for loyalty step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P15LimitLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty3 evaluates the window policy for loyalty step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P15WindowLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty4 evaluates the timeout policy for loyalty step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P15TimeoutLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty5 evaluates the quota policy for loyalty step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P15QuotaLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
