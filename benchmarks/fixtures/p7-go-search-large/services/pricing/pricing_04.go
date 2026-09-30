// Package pricing implements pricing service policies.
package pricing

// QuotaPricing0 evaluates the quota policy for pricing step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P04QuotaPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing1 evaluates the limit policy for pricing step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P04LimitPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing2 evaluates the timeout policy for pricing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P04TimeoutPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing3 evaluates the limit policy for pricing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P04LimitPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing4 evaluates the limit policy for pricing step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P04LimitPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing5 evaluates the retry policy for pricing step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P04RetryPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
