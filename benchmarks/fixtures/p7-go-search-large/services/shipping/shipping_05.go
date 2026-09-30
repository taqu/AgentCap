// Package shipping implements shipping service policies.
package shipping

// WindowShipping0 evaluates the window policy for shipping step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P05WindowShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping1 evaluates the limit policy for shipping step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P05LimitShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping2 evaluates the limit policy for shipping step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P05LimitShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping3 evaluates the limit policy for shipping step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P05LimitShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowShipping4 evaluates the window policy for shipping step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P05WindowShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutShipping5 evaluates the timeout policy for shipping step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P05TimeoutShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
