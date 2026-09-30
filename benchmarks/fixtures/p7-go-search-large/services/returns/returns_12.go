// Package returns implements returns service policies.
package returns

// LimitReturns0 evaluates the limit policy for returns step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P12LimitReturns0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns1 evaluates the retry policy for returns step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P12RetryReturns1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns2 evaluates the burst policy for returns step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P12BurstReturns2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstReturns3 evaluates the burst policy for returns step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P12BurstReturns3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitReturns4 evaluates the limit policy for returns step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P12LimitReturns4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryReturns5 evaluates the retry policy for returns step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "returns.limit").
func P12RetryReturns5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
