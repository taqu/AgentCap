// Package warehouse implements warehouse service policies.
package warehouse

// BurstWarehouse0 evaluates the burst policy for warehouse step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P14BurstWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse1 evaluates the timeout policy for warehouse step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P14TimeoutWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse2 evaluates the timeout policy for warehouse step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P14TimeoutWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse3 evaluates the timeout policy for warehouse step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P14TimeoutWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse4 evaluates the timeout policy for warehouse step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P14TimeoutWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse5 evaluates the timeout policy for warehouse step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P14TimeoutWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
