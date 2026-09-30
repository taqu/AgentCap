// Package returns implements returns service policies.
package returns

// RetryReturns0 evaluates the retry policy for returns step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P01RetryReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitReturns1 evaluates the limit policy for returns step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P01LimitReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns2 evaluates the retry policy for returns step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P01RetryReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitReturns3 evaluates the limit policy for returns step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P01LimitReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns4 evaluates the timeout policy for returns step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P01TimeoutReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns5 evaluates the retry policy for returns step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P01RetryReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
