// Package returns implements returns service policies.
package returns

// RetryReturns0 evaluates the retry policy for returns step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P03RetryReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns1 evaluates the retry policy for returns step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P03RetryReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns2 evaluates the window policy for returns step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P03WindowReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns3 evaluates the burst policy for returns step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P03BurstReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns4 evaluates the burst policy for returns step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P03BurstReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns5 evaluates the timeout policy for returns step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P03TimeoutReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
