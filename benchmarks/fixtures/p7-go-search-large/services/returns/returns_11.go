// Package returns implements returns service policies.
package returns

// TimeoutReturns0 evaluates the timeout policy for returns step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P11TimeoutReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns1 evaluates the burst policy for returns step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P11BurstReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns2 evaluates the retry policy for returns step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P11RetryReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns3 evaluates the retry policy for returns step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P11RetryReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns4 evaluates the timeout policy for returns step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P11TimeoutReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns5 evaluates the retry policy for returns step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P11RetryReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
