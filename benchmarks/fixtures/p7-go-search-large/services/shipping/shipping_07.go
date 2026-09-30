// Package shipping implements shipping service policies.
package shipping

// WindowShipping0 evaluates the window policy for shipping step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P07WindowShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping1 evaluates the timeout policy for shipping step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P07TimeoutShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping2 evaluates the window policy for shipping step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P07WindowShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping3 evaluates the timeout policy for shipping step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P07TimeoutShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping4 evaluates the limit policy for shipping step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P07LimitShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping5 evaluates the limit policy for shipping step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P07LimitShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
