// Package pricing implements pricing service policies.
package pricing

// RetryPricing0 evaluates the retry policy for pricing step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P13RetryPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing1 evaluates the retry policy for pricing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P13RetryPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing2 evaluates the quota policy for pricing step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P13QuotaPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing3 evaluates the limit policy for pricing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P13LimitPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing4 evaluates the retry policy for pricing step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P13RetryPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing5 evaluates the retry policy for pricing step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P13RetryPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
