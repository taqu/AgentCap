// Package warehouse implements warehouse service policies.
package warehouse

// QuotaWarehouse0 evaluates the quota policy for warehouse step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P17QuotaWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse1 evaluates the limit policy for warehouse step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P17LimitWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse2 evaluates the quota policy for warehouse step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P17QuotaWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse3 evaluates the limit policy for warehouse step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P17LimitWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutWarehouse4 evaluates the timeout policy for warehouse step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P17TimeoutWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse5 evaluates the limit policy for warehouse step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P17LimitWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
