// Package returns implements returns service policies.
package returns

// BurstReturns0 evaluates the burst policy for returns step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P00BurstReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns1 evaluates the retry policy for returns step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P00RetryReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns2 evaluates the timeout policy for returns step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P00TimeoutReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns3 evaluates the timeout policy for returns step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P00TimeoutReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns4 evaluates the quota policy for returns step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P00QuotaReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns5 evaluates the quota policy for returns step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P00QuotaReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
