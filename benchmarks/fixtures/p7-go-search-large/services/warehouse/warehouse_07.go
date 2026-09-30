// Package warehouse implements warehouse service policies.
package warehouse

// BurstWarehouse0 evaluates the burst policy for warehouse step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P07BurstWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse1 evaluates the burst policy for warehouse step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P07BurstWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse2 evaluates the timeout policy for warehouse step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P07TimeoutWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse3 evaluates the burst policy for warehouse step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P07BurstWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse4 evaluates the limit policy for warehouse step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P07LimitWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse5 evaluates the limit policy for warehouse step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P07LimitWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
