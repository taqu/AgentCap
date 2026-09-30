// Package warehouse implements warehouse service policies.
package warehouse

// WindowWarehouse0 evaluates the window policy for warehouse step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P11WindowWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse1 evaluates the timeout policy for warehouse step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P11TimeoutWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse2 evaluates the burst policy for warehouse step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P11BurstWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse3 evaluates the limit policy for warehouse step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P11LimitWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse4 evaluates the quota policy for warehouse step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P11QuotaWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse5 evaluates the burst policy for warehouse step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P11BurstWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
