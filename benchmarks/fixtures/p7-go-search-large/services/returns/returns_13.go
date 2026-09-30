// Package returns implements returns service policies.
package returns

// LimitReturns0 evaluates the limit policy for returns step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P13LimitReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns1 evaluates the window policy for returns step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P13WindowReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaReturns2 evaluates the quota policy for returns step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P13QuotaReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns3 evaluates the window policy for returns step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P13WindowReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitReturns4 evaluates the limit policy for returns step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P13LimitReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns5 evaluates the retry policy for returns step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P13RetryReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
