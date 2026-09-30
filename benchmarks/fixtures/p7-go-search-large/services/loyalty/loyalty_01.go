// Package loyalty implements loyalty service policies.
package loyalty

// QuotaLoyalty0 evaluates the quota policy for loyalty step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P01QuotaLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty1 evaluates the limit policy for loyalty step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P01LimitLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty2 evaluates the limit policy for loyalty step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P01LimitLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty3 evaluates the window policy for loyalty step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P01WindowLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstLoyalty4 evaluates the burst policy for loyalty step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P01BurstLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty5 evaluates the retry policy for loyalty step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P01RetryLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
