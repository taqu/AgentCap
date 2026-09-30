// Package tax implements tax service policies.
package tax

// TimeoutTax0 evaluates the timeout policy for tax step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P05TimeoutTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaTax1 evaluates the quota policy for tax step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P05QuotaTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowTax2 evaluates the window policy for tax step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P05WindowTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowTax3 evaluates the window policy for tax step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P05WindowTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitTax4 evaluates the limit policy for tax step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P05LimitTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstTax5 evaluates the burst policy for tax step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P05BurstTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
