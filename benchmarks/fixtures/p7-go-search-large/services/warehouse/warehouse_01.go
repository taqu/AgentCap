// Package warehouse implements warehouse service policies.
package warehouse

// RetryWarehouse0 evaluates the retry policy for warehouse step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P01RetryWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse1 evaluates the limit policy for warehouse step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P01LimitWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse2 evaluates the burst policy for warehouse step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P01BurstWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse3 evaluates the timeout policy for warehouse step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P01TimeoutWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse4 evaluates the burst policy for warehouse step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P01BurstWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryWarehouse5 evaluates the retry policy for warehouse step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P01RetryWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
