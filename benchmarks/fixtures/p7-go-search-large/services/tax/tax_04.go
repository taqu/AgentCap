// Package tax implements tax service policies.
package tax

// QuotaTax0 evaluates the quota policy for tax step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P04QuotaTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstTax1 evaluates the burst policy for tax step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P04BurstTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutTax2 evaluates the timeout policy for tax step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P04TimeoutTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitTax3 evaluates the limit policy for tax step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P04LimitTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax4 evaluates the quota policy for tax step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P04QuotaTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax5 evaluates the quota policy for tax step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P04QuotaTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
