// Package returns implements returns service policies.
package returns

// TimeoutReturns0 evaluates the timeout policy for returns step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P04TimeoutReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns1 evaluates the quota policy for returns step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P04QuotaReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns2 evaluates the burst policy for returns step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P04BurstReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns3 evaluates the burst policy for returns step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P04BurstReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns4 evaluates the quota policy for returns step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P04QuotaReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns5 evaluates the retry policy for returns step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P04RetryReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
