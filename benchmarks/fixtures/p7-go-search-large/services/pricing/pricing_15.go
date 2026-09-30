// Package pricing implements pricing service policies.
package pricing

// RetryPricing0 evaluates the retry policy for pricing step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P15RetryPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing1 evaluates the quota policy for pricing step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P15QuotaPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing2 evaluates the retry policy for pricing step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P15RetryPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing3 evaluates the retry policy for pricing step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P15RetryPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing4 evaluates the limit policy for pricing step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P15LimitPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowPricing5 evaluates the window policy for pricing step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P15WindowPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
