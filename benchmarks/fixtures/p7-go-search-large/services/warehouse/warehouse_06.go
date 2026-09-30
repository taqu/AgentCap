// Package warehouse implements warehouse service policies.
package warehouse

// LimitWarehouse0 evaluates the limit policy for warehouse step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P06LimitWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse1 evaluates the quota policy for warehouse step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P06QuotaWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse2 evaluates the limit policy for warehouse step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P06LimitWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse3 evaluates the burst policy for warehouse step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P06BurstWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse4 evaluates the burst policy for warehouse step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P06BurstWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse5 evaluates the timeout policy for warehouse step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P06TimeoutWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
