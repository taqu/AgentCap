// Package shipping implements shipping service policies.
package shipping

// LimitShipping0 evaluates the limit policy for shipping step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P15LimitShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping1 evaluates the window policy for shipping step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P15WindowShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping2 evaluates the burst policy for shipping step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P15BurstShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryShipping3 evaluates the retry policy for shipping step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P15RetryShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping4 evaluates the burst policy for shipping step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P15BurstShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping5 evaluates the quota policy for shipping step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P15QuotaShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
