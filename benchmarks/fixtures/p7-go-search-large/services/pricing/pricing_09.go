// Package pricing implements pricing service policies.
package pricing

// TimeoutPricing0 evaluates the timeout policy for pricing step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P09TimeoutPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing1 evaluates the timeout policy for pricing step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P09TimeoutPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing2 evaluates the limit policy for pricing step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P09LimitPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing3 evaluates the quota policy for pricing step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P09QuotaPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowPricing4 evaluates the window policy for pricing step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P09WindowPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing5 evaluates the limit policy for pricing step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P09LimitPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
