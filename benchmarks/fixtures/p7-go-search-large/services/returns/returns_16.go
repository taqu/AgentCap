// Package returns implements returns service policies.
package returns

// LimitReturns0 evaluates the limit policy for returns step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P16LimitReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns1 evaluates the retry policy for returns step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P16RetryReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns2 evaluates the window policy for returns step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P16WindowReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitReturns3 evaluates the limit policy for returns step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P16LimitReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutReturns4 evaluates the timeout policy for returns step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P16TimeoutReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowReturns5 evaluates the window policy for returns step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P16WindowReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
