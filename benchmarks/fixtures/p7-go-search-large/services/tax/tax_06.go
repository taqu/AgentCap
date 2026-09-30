// Package tax implements tax service policies.
package tax

// QuotaTax0 evaluates the quota policy for tax step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P06QuotaTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstTax1 evaluates the burst policy for tax step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P06BurstTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowTax2 evaluates the window policy for tax step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P06WindowTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstTax3 evaluates the burst policy for tax step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P06BurstTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstTax4 evaluates the burst policy for tax step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P06BurstTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowTax5 evaluates the window policy for tax step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P06WindowTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
