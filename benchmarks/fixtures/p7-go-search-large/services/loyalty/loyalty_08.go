// Package loyalty implements loyalty service policies.
package loyalty

// BurstLoyalty0 evaluates the burst policy for loyalty step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P08BurstLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty1 evaluates the retry policy for loyalty step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P08RetryLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty2 evaluates the timeout policy for loyalty step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P08TimeoutLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty3 evaluates the retry policy for loyalty step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P08RetryLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstLoyalty4 evaluates the burst policy for loyalty step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P08BurstLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty5 evaluates the quota policy for loyalty step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P08QuotaLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
