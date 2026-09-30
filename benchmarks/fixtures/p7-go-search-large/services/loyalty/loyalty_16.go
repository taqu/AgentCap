// Package loyalty implements loyalty service policies.
package loyalty

// WindowLoyalty0 evaluates the window policy for loyalty step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P16WindowLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty1 evaluates the window policy for loyalty step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P16WindowLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty2 evaluates the quota policy for loyalty step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P16QuotaLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty3 evaluates the limit policy for loyalty step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P16LimitLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty4 evaluates the retry policy for loyalty step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P16RetryLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty5 evaluates the window policy for loyalty step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P16WindowLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
