// Package loyalty implements loyalty service policies.
package loyalty

// RetryLoyalty0 evaluates the retry policy for loyalty step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P06RetryLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty1 evaluates the timeout policy for loyalty step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P06TimeoutLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty2 evaluates the retry policy for loyalty step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P06RetryLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty3 evaluates the retry policy for loyalty step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P06RetryLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty4 evaluates the quota policy for loyalty step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P06QuotaLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty5 evaluates the retry policy for loyalty step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P06RetryLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
