// Package tax implements tax service policies.
package tax

// RetryTax0 evaluates the retry policy for tax step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P03RetryTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitTax1 evaluates the limit policy for tax step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P03LimitTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutTax2 evaluates the timeout policy for tax step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P03TimeoutTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowTax3 evaluates the window policy for tax step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P03WindowTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstTax4 evaluates the burst policy for tax step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P03BurstTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitTax5 evaluates the limit policy for tax step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P03LimitTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
