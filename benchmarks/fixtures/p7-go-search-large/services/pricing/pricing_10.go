// Package pricing implements pricing service policies.
package pricing

// WindowPricing0 evaluates the window policy for pricing step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P10WindowPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing1 evaluates the retry policy for pricing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P10RetryPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing2 evaluates the timeout policy for pricing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P10TimeoutPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing3 evaluates the limit policy for pricing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P10LimitPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstPricing4 evaluates the burst policy for pricing step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P10BurstPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing5 evaluates the timeout policy for pricing step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P10TimeoutPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
