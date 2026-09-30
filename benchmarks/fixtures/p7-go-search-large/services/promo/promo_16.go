// Package promo implements promo service policies.
package promo

// LimitPromo0 evaluates the limit policy for promo step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P16LimitPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo1 evaluates the timeout policy for promo step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P16TimeoutPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo2 evaluates the timeout policy for promo step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P16TimeoutPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo3 evaluates the limit policy for promo step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P16LimitPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo4 evaluates the timeout policy for promo step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P16TimeoutPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo5 evaluates the limit policy for promo step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P16LimitPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
