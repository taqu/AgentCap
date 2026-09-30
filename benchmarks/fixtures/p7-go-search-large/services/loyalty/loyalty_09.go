// Package loyalty implements loyalty service policies.
package loyalty

// QuotaLoyalty0 evaluates the quota policy for loyalty step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P09QuotaLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty1 evaluates the quota policy for loyalty step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P09QuotaLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty2 evaluates the window policy for loyalty step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P09WindowLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty3 evaluates the retry policy for loyalty step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P09RetryLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty4 evaluates the timeout policy for loyalty step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P09TimeoutLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty5 evaluates the retry policy for loyalty step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P09RetryLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
