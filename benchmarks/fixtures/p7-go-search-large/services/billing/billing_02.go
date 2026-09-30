// Package billing implements billing service policies.
package billing

// QuotaBilling0 evaluates the quota policy for billing step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P02QuotaBilling0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling1 evaluates the limit policy for billing step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P02LimitBilling1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowBilling2 evaluates the window policy for billing step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P02WindowBilling2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling3 evaluates the limit policy for billing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P02LimitBilling3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitBilling4 evaluates the limit policy for billing step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P02LimitBilling4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaBilling5 evaluates the quota policy for billing step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "billing.limit").
func P02QuotaBilling5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
