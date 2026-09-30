// Package billing implements billing service policies.
package billing

// BurstBilling0 evaluates the burst policy for billing step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P09BurstBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstBilling1 evaluates the burst policy for billing step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P09BurstBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling2 evaluates the timeout policy for billing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P09TimeoutBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstBilling3 evaluates the burst policy for billing step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P09BurstBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling4 evaluates the limit policy for billing step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P09LimitBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling5 evaluates the timeout policy for billing step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P09TimeoutBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
