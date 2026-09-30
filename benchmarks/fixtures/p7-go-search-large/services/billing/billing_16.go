// Package billing implements billing service policies.
package billing

// WindowBilling0 evaluates the window policy for billing step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P16WindowBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling1 evaluates the window policy for billing step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P16WindowBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling2 evaluates the timeout policy for billing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P16TimeoutBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling3 evaluates the limit policy for billing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P16LimitBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaBilling4 evaluates the quota policy for billing step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P16QuotaBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling5 evaluates the retry policy for billing step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P16RetryBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
