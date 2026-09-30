// Package tax implements tax service policies.
package tax

// TimeoutTax0 evaluates the timeout policy for tax step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P11TimeoutTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstTax1 evaluates the burst policy for tax step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P11BurstTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutTax2 evaluates the timeout policy for tax step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P11TimeoutTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowTax3 evaluates the window policy for tax step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P11WindowTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax4 evaluates the quota policy for tax step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P11QuotaTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryTax5 evaluates the retry policy for tax step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P11RetryTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
