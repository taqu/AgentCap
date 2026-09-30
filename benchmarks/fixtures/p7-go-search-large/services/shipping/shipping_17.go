// Package shipping implements shipping service policies.
package shipping

// WindowShipping0 evaluates the window policy for shipping step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P17WindowShipping0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping1 evaluates the limit policy for shipping step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P17LimitShipping1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping2 evaluates the burst policy for shipping step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P17BurstShipping2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping3 evaluates the limit policy for shipping step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P17LimitShipping3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitShipping4 evaluates the limit policy for shipping step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P17LimitShipping4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstShipping5 evaluates the burst policy for shipping step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "shipping.limit").
func P17BurstShipping5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
