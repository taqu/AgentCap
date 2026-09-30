// Package promo implements promo service policies.
package promo

// LimitPromo0 evaluates the limit policy for promo step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P15LimitPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo1 evaluates the burst policy for promo step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P15BurstPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo2 evaluates the timeout policy for promo step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P15TimeoutPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowPromo3 evaluates the window policy for promo step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P15WindowPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo4 evaluates the retry policy for promo step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P15RetryPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo5 evaluates the quota policy for promo step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P15QuotaPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
