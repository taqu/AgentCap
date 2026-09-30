// Package pricing implements pricing service policies.
package pricing

// WindowPricing0 evaluates the window policy for pricing step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P08WindowPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstPricing1 evaluates the burst policy for pricing step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P08BurstPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstPricing2 evaluates the burst policy for pricing step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P08BurstPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing3 evaluates the limit policy for pricing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P08LimitPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing4 evaluates the limit policy for pricing step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P08LimitPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing5 evaluates the quota policy for pricing step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P08QuotaPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
