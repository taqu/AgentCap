// Package tax implements tax service policies.
package tax

// WindowTax0 evaluates the window policy for tax step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P13WindowTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax1 evaluates the quota policy for tax step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P13QuotaTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitTax2 evaluates the limit policy for tax step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P13LimitTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax3 evaluates the quota policy for tax step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P13QuotaTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstTax4 evaluates the burst policy for tax step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P13BurstTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax5 evaluates the quota policy for tax step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P13QuotaTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
