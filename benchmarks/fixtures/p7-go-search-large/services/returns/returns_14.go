// Package returns implements returns service policies.
package returns

// QuotaReturns0 evaluates the quota policy for returns step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P14QuotaReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns1 evaluates the window policy for returns step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P14WindowReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns2 evaluates the burst policy for returns step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P14BurstReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns3 evaluates the window policy for returns step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P14WindowReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns4 evaluates the retry policy for returns step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P14RetryReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns5 evaluates the quota policy for returns step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P14QuotaReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
