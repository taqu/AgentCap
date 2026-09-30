// Package tax implements tax service policies.
package tax

// RetryTax0 evaluates the retry policy for tax step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P12RetryTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowTax1 evaluates the window policy for tax step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P12WindowTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutTax2 evaluates the timeout policy for tax step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P12TimeoutTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutTax3 evaluates the timeout policy for tax step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P12TimeoutTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryTax4 evaluates the retry policy for tax step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P12RetryTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowTax5 evaluates the window policy for tax step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P12WindowTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
