// Package warehouse implements warehouse service policies.
package warehouse

// QuotaWarehouse0 evaluates the quota policy for warehouse step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P03QuotaWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse1 evaluates the limit policy for warehouse step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P03LimitWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowWarehouse2 evaluates the window policy for warehouse step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P03WindowWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstWarehouse3 evaluates the burst policy for warehouse step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P03BurstWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse4 evaluates the quota policy for warehouse step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P03QuotaWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse5 evaluates the quota policy for warehouse step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P03QuotaWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
