// Package pricing implements pricing service policies.
package pricing

// QuotaPricing0 evaluates the quota policy for pricing step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P05QuotaPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing1 evaluates the quota policy for pricing step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P05QuotaPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing2 evaluates the retry policy for pricing step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P05RetryPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstPricing3 evaluates the burst policy for pricing step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P05BurstPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstPricing4 evaluates the burst policy for pricing step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P05BurstPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstPricing5 evaluates the burst policy for pricing step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P05BurstPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
