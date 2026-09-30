// Package pricing implements pricing service policies.
package pricing

// RetryPricing0 evaluates the retry policy for pricing step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P12RetryPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing1 evaluates the quota policy for pricing step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P12QuotaPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing2 evaluates the timeout policy for pricing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P12TimeoutPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing3 evaluates the quota policy for pricing step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P12QuotaPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing4 evaluates the retry policy for pricing step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P12RetryPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing5 evaluates the retry policy for pricing step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P12RetryPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
