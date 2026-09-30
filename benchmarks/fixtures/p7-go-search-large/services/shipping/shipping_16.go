// Package shipping implements shipping service policies.
package shipping

// BurstShipping0 evaluates the burst policy for shipping step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P16BurstShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping1 evaluates the burst policy for shipping step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P16BurstShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping2 evaluates the limit policy for shipping step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P16LimitShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping3 evaluates the burst policy for shipping step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P16BurstShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping4 evaluates the burst policy for shipping step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P16BurstShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaShipping5 evaluates the quota policy for shipping step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P16QuotaShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
