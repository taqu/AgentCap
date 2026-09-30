// Package promo implements promo service policies.
package promo

// TimeoutPromo0 evaluates the timeout policy for promo step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P13TimeoutPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo1 evaluates the limit policy for promo step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P13LimitPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo2 evaluates the retry policy for promo step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P13RetryPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo3 evaluates the limit policy for promo step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P13LimitPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo4 evaluates the timeout policy for promo step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P13TimeoutPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo5 evaluates the quota policy for promo step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P13QuotaPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
