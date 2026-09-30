// Package billing implements billing service policies.
package billing

// QuotaBilling0 evaluates the quota policy for billing step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P11QuotaBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling1 evaluates the retry policy for billing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P11RetryBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling2 evaluates the timeout policy for billing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P11TimeoutBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling3 evaluates the retry policy for billing step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P11RetryBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling4 evaluates the timeout policy for billing step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P11TimeoutBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling5 evaluates the timeout policy for billing step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P11TimeoutBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
