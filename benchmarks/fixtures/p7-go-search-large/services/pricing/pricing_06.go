// Package pricing implements pricing service policies.
package pricing

// TimeoutPricing0 evaluates the timeout policy for pricing step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P06TimeoutPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing1 evaluates the timeout policy for pricing step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P06TimeoutPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing2 evaluates the timeout policy for pricing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P06TimeoutPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing3 evaluates the timeout policy for pricing step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P06TimeoutPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing4 evaluates the quota policy for pricing step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P06QuotaPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing5 evaluates the quota policy for pricing step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P06QuotaPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
