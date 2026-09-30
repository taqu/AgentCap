// Package tax implements tax service policies.
package tax

// LimitTax0 evaluates the limit policy for tax step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P10LimitTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitTax1 evaluates the limit policy for tax step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P10LimitTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryTax2 evaluates the retry policy for tax step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P10RetryTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax3 evaluates the quota policy for tax step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P10QuotaTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax4 evaluates the quota policy for tax step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P10QuotaTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstTax5 evaluates the burst policy for tax step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P10BurstTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
