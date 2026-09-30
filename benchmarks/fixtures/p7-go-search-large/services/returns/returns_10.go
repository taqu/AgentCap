// Package returns implements returns service policies.
package returns

// RetryReturns0 evaluates the retry policy for returns step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P10RetryReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns1 evaluates the window policy for returns step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P10WindowReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns2 evaluates the quota policy for returns step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P10QuotaReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns3 evaluates the quota policy for returns step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P10QuotaReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns4 evaluates the burst policy for returns step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P10BurstReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns5 evaluates the retry policy for returns step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P10RetryReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
