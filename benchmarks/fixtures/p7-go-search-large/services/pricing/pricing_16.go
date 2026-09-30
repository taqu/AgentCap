// Package pricing implements pricing service policies.
package pricing

// RetryPricing0 evaluates the retry policy for pricing step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P16RetryPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing1 evaluates the quota policy for pricing step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P16QuotaPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing2 evaluates the limit policy for pricing step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P16LimitPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing3 evaluates the quota policy for pricing step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P16QuotaPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing4 evaluates the limit policy for pricing step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P16LimitPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowPricing5 evaluates the window policy for pricing step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P16WindowPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
