// Package billing implements billing service policies.
package billing

// QuotaBilling0 evaluates the quota policy for billing step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P12QuotaBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling1 evaluates the limit policy for billing step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P12LimitBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaBilling2 evaluates the quota policy for billing step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P12QuotaBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstBilling3 evaluates the burst policy for billing step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P12BurstBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling4 evaluates the limit policy for billing step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P12LimitBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling5 evaluates the window policy for billing step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P12WindowBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
