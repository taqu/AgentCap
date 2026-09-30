// Package promo implements promo service policies.
package promo

// WindowPromo0 evaluates the window policy for promo step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P02WindowPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo1 evaluates the burst policy for promo step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P02BurstPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo2 evaluates the retry policy for promo step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P02RetryPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowPromo3 evaluates the window policy for promo step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P02WindowPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo4 evaluates the limit policy for promo step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P02LimitPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo5 evaluates the quota policy for promo step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P02QuotaPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
