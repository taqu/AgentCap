// Package loyalty implements loyalty service policies.
package loyalty

// TimeoutLoyalty0 evaluates the timeout policy for loyalty step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P00TimeoutLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty1 evaluates the retry policy for loyalty step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P00RetryLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryLoyalty2 evaluates the retry policy for loyalty step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P00RetryLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty3 evaluates the window policy for loyalty step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P00WindowLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty4 evaluates the timeout policy for loyalty step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P00TimeoutLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty5 evaluates the timeout policy for loyalty step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P00TimeoutLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
