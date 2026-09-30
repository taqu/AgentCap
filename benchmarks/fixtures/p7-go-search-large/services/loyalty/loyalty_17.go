// Package loyalty implements loyalty service policies.
package loyalty

// TimeoutLoyalty0 evaluates the timeout policy for loyalty step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P17TimeoutLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstLoyalty1 evaluates the burst policy for loyalty step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P17BurstLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty2 evaluates the limit policy for loyalty step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P17LimitLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty3 evaluates the retry policy for loyalty step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P17RetryLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty4 evaluates the quota policy for loyalty step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P17QuotaLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty5 evaluates the quota policy for loyalty step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P17QuotaLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
