// Package promo implements promo service policies.
package promo

// TimeoutPromo0 evaluates the timeout policy for promo step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P17TimeoutPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowPromo1 evaluates the window policy for promo step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P17WindowPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo2 evaluates the limit policy for promo step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P17LimitPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo3 evaluates the quota policy for promo step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P17QuotaPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo4 evaluates the retry policy for promo step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P17RetryPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowPromo5 evaluates the window policy for promo step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P17WindowPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
