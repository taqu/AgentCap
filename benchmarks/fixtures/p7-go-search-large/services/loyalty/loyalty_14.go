// Package loyalty implements loyalty service policies.
package loyalty

// TimeoutLoyalty0 evaluates the timeout policy for loyalty step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P14TimeoutLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty1 evaluates the limit policy for loyalty step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P14LimitLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty2 evaluates the quota policy for loyalty step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P14QuotaLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty3 evaluates the window policy for loyalty step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P14WindowLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty4 evaluates the quota policy for loyalty step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P14QuotaLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty5 evaluates the window policy for loyalty step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P14WindowLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
