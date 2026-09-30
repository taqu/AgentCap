// Package warehouse implements warehouse service policies.
package warehouse

// QuotaWarehouse0 evaluates the quota policy for warehouse step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P00QuotaWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse1 evaluates the limit policy for warehouse step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P00LimitWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse2 evaluates the burst policy for warehouse step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P00BurstWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryWarehouse3 evaluates the retry policy for warehouse step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P00RetryWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse4 evaluates the limit policy for warehouse step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P00LimitWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse5 evaluates the burst policy for warehouse step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P00BurstWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
