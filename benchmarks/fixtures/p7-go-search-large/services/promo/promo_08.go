// Package promo implements promo service policies.
package promo

// RetryPromo0 evaluates the retry policy for promo step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P08RetryPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowPromo1 evaluates the window policy for promo step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P08WindowPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo2 evaluates the timeout policy for promo step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P08TimeoutPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo3 evaluates the limit policy for promo step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P08LimitPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo4 evaluates the burst policy for promo step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P08BurstPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo5 evaluates the timeout policy for promo step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P08TimeoutPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
