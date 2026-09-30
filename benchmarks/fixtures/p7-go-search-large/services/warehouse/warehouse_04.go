// Package warehouse implements warehouse service policies.
package warehouse

// WindowWarehouse0 evaluates the window policy for warehouse step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P04WindowWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse1 evaluates the burst policy for warehouse step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P04BurstWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowWarehouse2 evaluates the window policy for warehouse step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P04WindowWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowWarehouse3 evaluates the window policy for warehouse step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P04WindowWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse4 evaluates the limit policy for warehouse step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P04LimitWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse5 evaluates the burst policy for warehouse step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P04BurstWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
