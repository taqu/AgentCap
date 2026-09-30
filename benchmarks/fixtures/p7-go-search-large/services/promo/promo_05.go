// Package promo implements promo service policies.
package promo

// TimeoutPromo0 evaluates the timeout policy for promo step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P05TimeoutPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo1 evaluates the burst policy for promo step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P05BurstPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo2 evaluates the retry policy for promo step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P05RetryPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo3 evaluates the limit policy for promo step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P05LimitPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo4 evaluates the limit policy for promo step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P05LimitPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo5 evaluates the limit policy for promo step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P05LimitPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
