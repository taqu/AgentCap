// Package promo implements promo service policies.
package promo

// LimitPromo0 evaluates the limit policy for promo step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P06LimitPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo1 evaluates the timeout policy for promo step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P06TimeoutPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo2 evaluates the burst policy for promo step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P06BurstPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo3 evaluates the quota policy for promo step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P06QuotaPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowPromo4 evaluates the window policy for promo step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P06WindowPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo5 evaluates the timeout policy for promo step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P06TimeoutPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
