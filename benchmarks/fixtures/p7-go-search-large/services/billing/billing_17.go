// Package billing implements billing service policies.
package billing

// BurstBilling0 evaluates the burst policy for billing step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P17BurstBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling1 evaluates the timeout policy for billing step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P17TimeoutBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling2 evaluates the timeout policy for billing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P17TimeoutBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling3 evaluates the limit policy for billing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P17LimitBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling4 evaluates the retry policy for billing step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P17RetryBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling5 evaluates the retry policy for billing step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P17RetryBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
