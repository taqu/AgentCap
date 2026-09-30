// Package billing implements billing service policies.
package billing

// RetryBilling0 evaluates the retry policy for billing step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P05RetryBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling1 evaluates the retry policy for billing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P05RetryBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling2 evaluates the timeout policy for billing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P05TimeoutBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaBilling3 evaluates the quota policy for billing step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P05QuotaBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling4 evaluates the limit policy for billing step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P05LimitBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling5 evaluates the timeout policy for billing step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P05TimeoutBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
