// Package shipping implements shipping service policies.
package shipping

// WindowShipping0 evaluates the window policy for shipping step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P09WindowShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping1 evaluates the window policy for shipping step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P09WindowShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping2 evaluates the burst policy for shipping step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P09BurstShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping3 evaluates the quota policy for shipping step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P09QuotaShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping4 evaluates the quota policy for shipping step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P09QuotaShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping5 evaluates the retry policy for shipping step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P09RetryShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
