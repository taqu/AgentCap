// Package pricing implements pricing service policies.
package pricing

// WindowPricing0 evaluates the window policy for pricing step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P03WindowPricing0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing1 evaluates the retry policy for pricing step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P03RetryPricing1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaPricing2 evaluates the quota policy for pricing step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P03QuotaPricing2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitPricing3 evaluates the limit policy for pricing step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P03LimitPricing3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryPricing4 evaluates the retry policy for pricing step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P03RetryPricing4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowPricing5 evaluates the window policy for pricing step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "pricing.limit").
func P03WindowPricing5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
