// Package warehouse implements warehouse service policies.
package warehouse

// BurstWarehouse0 evaluates the burst policy for warehouse step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P09BurstWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryWarehouse1 evaluates the retry policy for warehouse step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P09RetryWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowWarehouse2 evaluates the window policy for warehouse step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P09WindowWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse3 evaluates the limit policy for warehouse step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P09LimitWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse4 evaluates the burst policy for warehouse step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P09BurstWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowWarehouse5 evaluates the window policy for warehouse step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P09WindowWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
