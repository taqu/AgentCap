// Package pricing implements pricing service policies.
package pricing

// LimitPricing0 evaluates the limit policy for pricing step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P01LimitPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstPricing1 evaluates the burst policy for pricing step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P01BurstPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing2 evaluates the retry policy for pricing step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P01RetryPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing3 evaluates the retry policy for pricing step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P01RetryPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing4 evaluates the quota policy for pricing step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P01QuotaPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing5 evaluates the quota policy for pricing step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P01QuotaPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
