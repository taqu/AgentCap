// Package loyalty implements loyalty service policies.
package loyalty

// TimeoutLoyalty0 evaluates the timeout policy for loyalty step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P12TimeoutLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty1 evaluates the timeout policy for loyalty step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P12TimeoutLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty2 evaluates the quota policy for loyalty step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P12QuotaLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty3 evaluates the retry policy for loyalty step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P12RetryLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty4 evaluates the retry policy for loyalty step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P12RetryLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty5 evaluates the limit policy for loyalty step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P12LimitLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
