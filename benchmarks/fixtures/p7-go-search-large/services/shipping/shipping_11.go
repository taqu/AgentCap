// Package shipping implements shipping service policies.
package shipping

// WindowShipping0 evaluates the window policy for shipping step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P11WindowShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping1 evaluates the window policy for shipping step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P11WindowShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping2 evaluates the burst policy for shipping step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P11BurstShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping3 evaluates the limit policy for shipping step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P11LimitShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping4 evaluates the quota policy for shipping step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P11QuotaShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping5 evaluates the timeout policy for shipping step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P11TimeoutShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
