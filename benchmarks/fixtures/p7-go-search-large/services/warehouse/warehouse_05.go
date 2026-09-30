// Package warehouse implements warehouse service policies.
package warehouse

// BurstWarehouse0 evaluates the burst policy for warehouse step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P05BurstWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse1 evaluates the burst policy for warehouse step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P05BurstWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse2 evaluates the limit policy for warehouse step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P05LimitWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryWarehouse3 evaluates the retry policy for warehouse step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P05RetryWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse4 evaluates the limit policy for warehouse step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P05LimitWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse5 evaluates the quota policy for warehouse step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P05QuotaWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
