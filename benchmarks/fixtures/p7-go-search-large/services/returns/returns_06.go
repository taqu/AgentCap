// Package returns implements returns service policies.
package returns

// QuotaReturns0 evaluates the quota policy for returns step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P06QuotaReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns1 evaluates the quota policy for returns step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P06QuotaReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns2 evaluates the window policy for returns step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P06WindowReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns3 evaluates the timeout policy for returns step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P06TimeoutReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns4 evaluates the quota policy for returns step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P06QuotaReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns5 evaluates the quota policy for returns step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P06QuotaReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
