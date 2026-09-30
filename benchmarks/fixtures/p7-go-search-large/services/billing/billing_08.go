// Package billing implements billing service policies.
package billing

// WindowBilling0 evaluates the window policy for billing step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P08WindowBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling1 evaluates the window policy for billing step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P08WindowBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling2 evaluates the limit policy for billing step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P08LimitBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutBilling3 evaluates the timeout policy for billing step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P08TimeoutBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstBilling4 evaluates the burst policy for billing step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P08BurstBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling5 evaluates the limit policy for billing step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P08LimitBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
