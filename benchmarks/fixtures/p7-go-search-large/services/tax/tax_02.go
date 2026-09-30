// Package tax implements tax service policies.
package tax

// BurstTax0 evaluates the burst policy for tax step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P02BurstTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutTax1 evaluates the timeout policy for tax step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P02TimeoutTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryTax2 evaluates the retry policy for tax step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P02RetryTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowTax3 evaluates the window policy for tax step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P02WindowTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstTax4 evaluates the burst policy for tax step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P02BurstTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax5 evaluates the quota policy for tax step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P02QuotaTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
