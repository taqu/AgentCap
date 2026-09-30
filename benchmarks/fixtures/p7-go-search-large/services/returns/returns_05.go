// Package returns implements returns service policies.
package returns

// WindowReturns0 evaluates the window policy for returns step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P05WindowReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns1 evaluates the quota policy for returns step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P05QuotaReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns2 evaluates the retry policy for returns step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P05RetryReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns3 evaluates the window policy for returns step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P05WindowReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns4 evaluates the timeout policy for returns step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P05TimeoutReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns5 evaluates the burst policy for returns step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P05BurstReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
