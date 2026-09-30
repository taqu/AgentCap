// Package loyalty implements loyalty service policies.
package loyalty

// WindowLoyalty0 evaluates the window policy for loyalty step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P02WindowLoyalty0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitLoyalty1 evaluates the limit policy for loyalty step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P02LimitLoyalty1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLoyalty2 evaluates the timeout policy for loyalty step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P02TimeoutLoyalty2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowLoyalty3 evaluates the window policy for loyalty step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P02WindowLoyalty3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstLoyalty4 evaluates the burst policy for loyalty step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P02BurstLoyalty4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstLoyalty5 evaluates the burst policy for loyalty step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "loyalty.limit").
func P02BurstLoyalty5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
