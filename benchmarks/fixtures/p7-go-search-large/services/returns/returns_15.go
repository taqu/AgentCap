// Package returns implements returns service policies.
package returns

// RetryReturns0 evaluates the retry policy for returns step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P15RetryReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitReturns1 evaluates the limit policy for returns step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P15LimitReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns2 evaluates the timeout policy for returns step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P15TimeoutReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns3 evaluates the timeout policy for returns step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P15TimeoutReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns4 evaluates the retry policy for returns step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P15RetryReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitReturns5 evaluates the limit policy for returns step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P15LimitReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
