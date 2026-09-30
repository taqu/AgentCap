// Package loyalty implements loyalty service policies.
package loyalty

// QuotaLoyalty0 evaluates the quota policy for loyalty step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P13QuotaLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty1 evaluates the retry policy for loyalty step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P13RetryLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty2 evaluates the timeout policy for loyalty step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P13TimeoutLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty3 evaluates the quota policy for loyalty step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P13QuotaLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty4 evaluates the limit policy for loyalty step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P13LimitLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty5 evaluates the timeout policy for loyalty step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P13TimeoutLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
