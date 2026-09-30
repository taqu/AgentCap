// Package billing implements billing service policies.
package billing

// TimeoutBilling0 evaluates the timeout policy for billing step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P00TimeoutBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling1 evaluates the retry policy for billing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P00RetryBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling2 evaluates the limit policy for billing step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P00LimitBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling3 evaluates the limit policy for billing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P00LimitBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling4 evaluates the window policy for billing step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P00WindowBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling5 evaluates the window policy for billing step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P00WindowBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
