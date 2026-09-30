// Package tax implements tax service policies.
package tax

// QuotaTax0 evaluates the quota policy for tax step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P08QuotaTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax1 evaluates the quota policy for tax step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P08QuotaTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowTax2 evaluates the window policy for tax step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P08WindowTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax3 evaluates the quota policy for tax step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P08QuotaTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryTax4 evaluates the retry policy for tax step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P08RetryTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryTax5 evaluates the retry policy for tax step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P08RetryTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
