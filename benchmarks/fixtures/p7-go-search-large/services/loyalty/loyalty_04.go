// Package loyalty implements loyalty service policies.
package loyalty

// LimitLoyalty0 evaluates the limit policy for loyalty step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P04LimitLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty1 evaluates the limit policy for loyalty step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P04LimitLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty2 evaluates the timeout policy for loyalty step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P04TimeoutLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstLoyalty3 evaluates the burst policy for loyalty step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P04BurstLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty4 evaluates the limit policy for loyalty step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P04LimitLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaLoyalty5 evaluates the quota policy for loyalty step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P04QuotaLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
