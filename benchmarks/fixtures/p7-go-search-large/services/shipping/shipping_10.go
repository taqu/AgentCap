// Package shipping implements shipping service policies.
package shipping

// LimitShipping0 evaluates the limit policy for shipping step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P10LimitShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping1 evaluates the limit policy for shipping step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P10LimitShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping2 evaluates the burst policy for shipping step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P10BurstShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping3 evaluates the retry policy for shipping step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P10RetryShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping4 evaluates the limit policy for shipping step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P10LimitShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping5 evaluates the burst policy for shipping step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P10BurstShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
