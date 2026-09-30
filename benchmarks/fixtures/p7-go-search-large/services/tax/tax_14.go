// Package tax implements tax service policies.
package tax

// BurstTax0 evaluates the burst policy for tax step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P14BurstTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitTax1 evaluates the limit policy for tax step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P14LimitTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstTax2 evaluates the burst policy for tax step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P14BurstTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryTax3 evaluates the retry policy for tax step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P14RetryTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowTax4 evaluates the window policy for tax step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P14WindowTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryTax5 evaluates the retry policy for tax step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P14RetryTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
