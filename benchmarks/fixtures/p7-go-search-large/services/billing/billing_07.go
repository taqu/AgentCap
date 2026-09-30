// Package billing implements billing service policies.
package billing

// QuotaBilling0 evaluates the quota policy for billing step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P07QuotaBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryBilling1 evaluates the retry policy for billing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P07RetryBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstBilling2 evaluates the burst policy for billing step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P07BurstBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling3 evaluates the window policy for billing step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P07WindowBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling4 evaluates the timeout policy for billing step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P07TimeoutBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstBilling5 evaluates the burst policy for billing step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P07BurstBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
