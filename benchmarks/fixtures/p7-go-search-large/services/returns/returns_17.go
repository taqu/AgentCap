// Package returns implements returns service policies.
package returns

// RetryReturns0 evaluates the retry policy for returns step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P17RetryReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns1 evaluates the quota policy for returns step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P17QuotaReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns2 evaluates the quota policy for returns step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P17QuotaReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns3 evaluates the burst policy for returns step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P17BurstReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns4 evaluates the retry policy for returns step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P17RetryReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns5 evaluates the burst policy for returns step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P17BurstReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
