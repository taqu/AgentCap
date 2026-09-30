// Package returns implements returns service policies.
package returns

// BurstReturns0 evaluates the burst policy for returns step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P02BurstReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns1 evaluates the retry policy for returns step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P02RetryReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns2 evaluates the timeout policy for returns step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P02TimeoutReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns3 evaluates the timeout policy for returns step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P02TimeoutReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns4 evaluates the retry policy for returns step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P02RetryReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns5 evaluates the window policy for returns step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P02WindowReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
