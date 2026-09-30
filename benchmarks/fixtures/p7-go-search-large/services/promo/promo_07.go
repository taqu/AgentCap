// Package promo implements promo service policies.
package promo

// WindowPromo0 evaluates the window policy for promo step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P07WindowPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo1 evaluates the retry policy for promo step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P07RetryPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo2 evaluates the quota policy for promo step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P07QuotaPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo3 evaluates the timeout policy for promo step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P07TimeoutPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo4 evaluates the burst policy for promo step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P07BurstPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryPromo5 evaluates the retry policy for promo step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P07RetryPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
