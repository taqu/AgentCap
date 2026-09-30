// Package shipping implements shipping service policies.
package shipping

// TimeoutShipping0 evaluates the timeout policy for shipping step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P14TimeoutShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping1 evaluates the quota policy for shipping step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P14QuotaShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping2 evaluates the window policy for shipping step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P14WindowShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping3 evaluates the retry policy for shipping step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P14RetryShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping4 evaluates the timeout policy for shipping step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P14TimeoutShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping5 evaluates the burst policy for shipping step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P14BurstShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
