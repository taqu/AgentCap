// Package shipping implements shipping service policies.
package shipping

// QuotaShipping0 evaluates the quota policy for shipping step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P00QuotaShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping1 evaluates the quota policy for shipping step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P00QuotaShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping2 evaluates the window policy for shipping step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P00WindowShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping3 evaluates the timeout policy for shipping step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P00TimeoutShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping4 evaluates the quota policy for shipping step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P00QuotaShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping5 evaluates the timeout policy for shipping step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P00TimeoutShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
