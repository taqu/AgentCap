// Package billing implements billing service policies.
package billing

// WindowBilling0 evaluates the window policy for billing step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P10WindowBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling1 evaluates the retry policy for billing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P10RetryBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaBilling2 evaluates the quota policy for billing step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P10QuotaBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaBilling3 evaluates the quota policy for billing step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P10QuotaBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling4 evaluates the retry policy for billing step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P10RetryBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling5 evaluates the retry policy for billing step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P10RetryBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
