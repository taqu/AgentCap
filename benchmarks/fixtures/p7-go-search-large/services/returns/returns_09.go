// Package returns implements returns service policies.
package returns

// BurstReturns0 evaluates the burst policy for returns step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P09BurstReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns1 evaluates the window policy for returns step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P09WindowReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns2 evaluates the retry policy for returns step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P09RetryReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns3 evaluates the quota policy for returns step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P09QuotaReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns4 evaluates the window policy for returns step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P09WindowReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitReturns5 evaluates the limit policy for returns step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P09LimitReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
