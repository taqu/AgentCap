// Package shipping implements shipping service policies.
package shipping

// TimeoutShipping0 evaluates the timeout policy for shipping step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P02TimeoutShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping1 evaluates the limit policy for shipping step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P02LimitShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping2 evaluates the quota policy for shipping step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P02QuotaShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping3 evaluates the limit policy for shipping step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P02LimitShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping4 evaluates the timeout policy for shipping step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P02TimeoutShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping5 evaluates the quota policy for shipping step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P02QuotaShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
