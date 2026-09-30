// Package loyalty implements loyalty service policies.
package loyalty

// RetryLoyalty0 evaluates the retry policy for loyalty step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P11RetryLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty1 evaluates the quota policy for loyalty step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P11QuotaLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty2 evaluates the limit policy for loyalty step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P11LimitLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty3 evaluates the window policy for loyalty step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P11WindowLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty4 evaluates the timeout policy for loyalty step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P11TimeoutLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty5 evaluates the timeout policy for loyalty step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P11TimeoutLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
