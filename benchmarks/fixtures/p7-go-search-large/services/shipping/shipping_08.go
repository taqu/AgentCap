// Package shipping implements shipping service policies.
package shipping

// LimitShipping0 evaluates the limit policy for shipping step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P08LimitShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping1 evaluates the window policy for shipping step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P08WindowShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping2 evaluates the window policy for shipping step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P08WindowShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping3 evaluates the quota policy for shipping step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P08QuotaShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping4 evaluates the limit policy for shipping step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P08LimitShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping5 evaluates the burst policy for shipping step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P08BurstShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
