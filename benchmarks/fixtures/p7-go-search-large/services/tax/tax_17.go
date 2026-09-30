// Package tax implements tax service policies.
package tax

// LimitTax0 evaluates the limit policy for tax step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P17LimitTax0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitTax1 evaluates the limit policy for tax step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P17LimitTax1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitTax2 evaluates the limit policy for tax step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P17LimitTax2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutTax3 evaluates the timeout policy for tax step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P17TimeoutTax3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutTax4 evaluates the timeout policy for tax step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P17TimeoutTax4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstTax5 evaluates the burst policy for tax step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "tax.limit").
func P17BurstTax5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
