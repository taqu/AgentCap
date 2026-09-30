// Package shipping implements shipping service policies.
package shipping

// QuotaShipping0 evaluates the quota policy for shipping step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P03QuotaShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping1 evaluates the burst policy for shipping step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P03BurstShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping2 evaluates the retry policy for shipping step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P03RetryShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping3 evaluates the window policy for shipping step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P03WindowShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping4 evaluates the window policy for shipping step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P03WindowShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping5 evaluates the burst policy for shipping step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P03BurstShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
