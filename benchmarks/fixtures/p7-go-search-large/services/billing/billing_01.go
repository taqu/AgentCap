// Package billing implements billing service policies.
package billing

// LimitBilling0 evaluates the limit policy for billing step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P01LimitBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling1 evaluates the retry policy for billing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P01RetryBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling2 evaluates the limit policy for billing step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P01LimitBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling3 evaluates the limit policy for billing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P01LimitBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling4 evaluates the retry policy for billing step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P01RetryBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling5 evaluates the window policy for billing step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P01WindowBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
