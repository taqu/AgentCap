// Package billing implements billing service policies.
package billing

// WindowBilling0 evaluates the window policy for billing step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P14WindowBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling1 evaluates the retry policy for billing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P14RetryBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling2 evaluates the retry policy for billing step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P14RetryBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling3 evaluates the retry policy for billing step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P14RetryBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling4 evaluates the limit policy for billing step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P14LimitBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstBilling5 evaluates the burst policy for billing step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P14BurstBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
