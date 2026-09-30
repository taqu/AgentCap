// Package promo implements promo service policies.
package promo

// WindowPromo0 evaluates the window policy for promo step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P14WindowPromo0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo1 evaluates the timeout policy for promo step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P14TimeoutPromo1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPromo2 evaluates the timeout policy for promo step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P14TimeoutPromo2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowPromo3 evaluates the window policy for promo step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P14WindowPromo3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo4 evaluates the quota policy for promo step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P14QuotaPromo4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaPromo5 evaluates the quota policy for promo step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "promo.limit").
func P14QuotaPromo5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
