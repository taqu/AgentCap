// Package promo implements promo service policies.
package promo

// WindowPromo0 evaluates the window policy for promo step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P01WindowPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo1 evaluates the limit policy for promo step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P01LimitPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo2 evaluates the retry policy for promo step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P01RetryPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo3 evaluates the burst policy for promo step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P01BurstPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo4 evaluates the retry policy for promo step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P01RetryPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo5 evaluates the limit policy for promo step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P01LimitPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
