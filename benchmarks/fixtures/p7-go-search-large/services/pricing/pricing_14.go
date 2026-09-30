// Package pricing implements pricing service policies.
package pricing

// BurstPricing0 evaluates the burst policy for pricing step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P14BurstPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing1 evaluates the retry policy for pricing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P14RetryPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing2 evaluates the timeout policy for pricing step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P14TimeoutPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing3 evaluates the timeout policy for pricing step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P14TimeoutPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing4 evaluates the retry policy for pricing step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P14RetryPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutPricing5 evaluates the timeout policy for pricing step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P14TimeoutPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
