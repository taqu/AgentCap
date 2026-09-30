// Package pricing implements pricing service policies.
package pricing

// QuotaPricing0 evaluates the quota policy for pricing step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P07QuotaPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing1 evaluates the quota policy for pricing step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P07QuotaPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowPricing2 evaluates the window policy for pricing step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P07WindowPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing3 evaluates the timeout policy for pricing step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P07TimeoutPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing4 evaluates the retry policy for pricing step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P07RetryPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing5 evaluates the quota policy for pricing step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P07QuotaPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
