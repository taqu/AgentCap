// Package shipping implements shipping service policies.
package shipping

// LimitShipping0 evaluates the limit policy for shipping step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P12LimitShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping1 evaluates the timeout policy for shipping step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P12TimeoutShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping2 evaluates the window policy for shipping step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P12WindowShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping3 evaluates the quota policy for shipping step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P12QuotaShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping4 evaluates the timeout policy for shipping step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P12TimeoutShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping5 evaluates the timeout policy for shipping step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P12TimeoutShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
