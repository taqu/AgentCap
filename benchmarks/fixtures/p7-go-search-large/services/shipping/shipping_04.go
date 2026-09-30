// Package shipping implements shipping service policies.
package shipping

// LimitShipping0 evaluates the limit policy for shipping step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P04LimitShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping1 evaluates the limit policy for shipping step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P04LimitShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping2 evaluates the retry policy for shipping step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P04RetryShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping3 evaluates the timeout policy for shipping step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P04TimeoutShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping4 evaluates the burst policy for shipping step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P04BurstShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping5 evaluates the timeout policy for shipping step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P04TimeoutShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
