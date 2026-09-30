// Package warehouse implements warehouse service policies.
package warehouse

// TimeoutWarehouse0 evaluates the timeout policy for warehouse step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P08TimeoutWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse1 evaluates the burst policy for warehouse step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P08BurstWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse2 evaluates the quota policy for warehouse step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P08QuotaWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse3 evaluates the quota policy for warehouse step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P08QuotaWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse4 evaluates the timeout policy for warehouse step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P08TimeoutWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse5 evaluates the quota policy for warehouse step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P08QuotaWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
