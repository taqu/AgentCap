// Package loyalty implements loyalty service policies.
package loyalty

// LimitLoyalty0 evaluates the limit policy for loyalty step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P07LimitLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty1 evaluates the quota policy for loyalty step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P07QuotaLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty2 evaluates the window policy for loyalty step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P07WindowLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstLoyalty3 evaluates the burst policy for loyalty step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P07BurstLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty4 evaluates the retry policy for loyalty step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P07RetryLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty5 evaluates the retry policy for loyalty step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P07RetryLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
