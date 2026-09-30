// Package loyalty implements loyalty service policies.
package loyalty

// QuotaLoyalty0 evaluates the quota policy for loyalty step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P10QuotaLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty1 evaluates the timeout policy for loyalty step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P10TimeoutLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty2 evaluates the timeout policy for loyalty step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P10TimeoutLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty3 evaluates the quota policy for loyalty step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P10QuotaLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty4 evaluates the retry policy for loyalty step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P10RetryLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty5 evaluates the timeout policy for loyalty step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P10TimeoutLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
