// Package tax implements tax service policies.
package tax

// RetryTax0 evaluates the retry policy for tax step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P09RetryTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryTax1 evaluates the retry policy for tax step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P09RetryTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstTax2 evaluates the burst policy for tax step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P09BurstTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax3 evaluates the quota policy for tax step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P09QuotaTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstTax4 evaluates the burst policy for tax step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P09BurstTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitTax5 evaluates the limit policy for tax step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P09LimitTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
