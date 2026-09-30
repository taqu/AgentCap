// Package promo implements promo service policies.
package promo

// LimitPromo0 evaluates the limit policy for promo step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P11LimitPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo1 evaluates the quota policy for promo step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P11QuotaPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo2 evaluates the quota policy for promo step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P11QuotaPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowPromo3 evaluates the window policy for promo step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P11WindowPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo4 evaluates the timeout policy for promo step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P11TimeoutPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo5 evaluates the limit policy for promo step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P11LimitPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
