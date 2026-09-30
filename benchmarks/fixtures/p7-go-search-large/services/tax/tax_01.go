// Package tax implements tax service policies.
package tax

// BurstTax0 evaluates the burst policy for tax step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P01BurstTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstTax1 evaluates the burst policy for tax step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P01BurstTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryTax2 evaluates the retry policy for tax step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P01RetryTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryTax3 evaluates the retry policy for tax step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P01RetryTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstTax4 evaluates the burst policy for tax step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P01BurstTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutTax5 evaluates the timeout policy for tax step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P01TimeoutTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
