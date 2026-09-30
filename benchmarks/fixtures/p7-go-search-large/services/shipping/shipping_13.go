// Package shipping implements shipping service policies.
package shipping

// LimitShipping0 evaluates the limit policy for shipping step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P13LimitShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping1 evaluates the retry policy for shipping step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P13RetryShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping2 evaluates the retry policy for shipping step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P13RetryShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping3 evaluates the window policy for shipping step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P13WindowShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping4 evaluates the quota policy for shipping step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P13QuotaShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping5 evaluates the limit policy for shipping step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P13LimitShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
