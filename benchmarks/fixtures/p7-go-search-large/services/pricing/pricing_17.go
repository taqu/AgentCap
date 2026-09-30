// Package pricing implements pricing service policies.
package pricing

// BurstPricing0 evaluates the burst policy for pricing step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P17BurstPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing1 evaluates the quota policy for pricing step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P17QuotaPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing2 evaluates the timeout policy for pricing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P17TimeoutPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing3 evaluates the quota policy for pricing step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P17QuotaPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing4 evaluates the quota policy for pricing step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P17QuotaPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing5 evaluates the limit policy for pricing step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P17LimitPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
