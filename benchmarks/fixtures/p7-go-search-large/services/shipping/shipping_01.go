// Package shipping implements shipping service policies.
package shipping

// BurstShipping0 evaluates the burst policy for shipping step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P01BurstShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping1 evaluates the limit policy for shipping step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P01LimitShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping2 evaluates the window policy for shipping step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P01WindowShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping3 evaluates the retry policy for shipping step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P01RetryShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping4 evaluates the retry policy for shipping step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P01RetryShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping5 evaluates the window policy for shipping step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P01WindowShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
