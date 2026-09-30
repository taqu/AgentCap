// Package warehouse implements warehouse service policies.
package warehouse

// WindowWarehouse0 evaluates the window policy for warehouse step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P02WindowWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse1 evaluates the quota policy for warehouse step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P02QuotaWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse2 evaluates the limit policy for warehouse step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P02LimitWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse3 evaluates the limit policy for warehouse step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P02LimitWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse4 evaluates the limit policy for warehouse step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P02LimitWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse5 evaluates the quota policy for warehouse step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P02QuotaWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
