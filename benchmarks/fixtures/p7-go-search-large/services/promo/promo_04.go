// Package promo implements promo service policies.
package promo

// TimeoutPromo0 evaluates the timeout policy for promo step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P04TimeoutPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo1 evaluates the quota policy for promo step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P04QuotaPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo2 evaluates the burst policy for promo step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P04BurstPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitPromo3 evaluates the limit policy for promo step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P04LimitPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo4 evaluates the burst policy for promo step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P04BurstPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstPromo5 evaluates the burst policy for promo step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P04BurstPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
