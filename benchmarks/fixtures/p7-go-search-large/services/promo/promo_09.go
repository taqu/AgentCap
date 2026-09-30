// Package promo implements promo service policies.
package promo

// BurstPromo0 evaluates the burst policy for promo step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P09BurstPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo1 evaluates the quota policy for promo step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P09QuotaPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo2 evaluates the timeout policy for promo step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P09TimeoutPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo3 evaluates the retry policy for promo step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P09RetryPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo4 evaluates the retry policy for promo step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P09RetryPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo5 evaluates the burst policy for promo step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P09BurstPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
