// Package billing implements billing service policies.
package billing

// TimeoutBilling0 evaluates the timeout policy for billing step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P04TimeoutBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling1 evaluates the window policy for billing step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P04WindowBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling2 evaluates the window policy for billing step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P04WindowBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling3 evaluates the window policy for billing step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P04WindowBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling4 evaluates the retry policy for billing step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P04RetryBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling5 evaluates the retry policy for billing step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P04RetryBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
