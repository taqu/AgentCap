// Package warehouse implements warehouse service policies.
package warehouse

// BurstWarehouse0 evaluates the burst policy for warehouse step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P15BurstWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse1 evaluates the timeout policy for warehouse step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P15TimeoutWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowWarehouse2 evaluates the window policy for warehouse step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P15WindowWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse3 evaluates the timeout policy for warehouse step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P15TimeoutWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse4 evaluates the burst policy for warehouse step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P15BurstWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryWarehouse5 evaluates the retry policy for warehouse step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P15RetryWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
