// Package warehouse implements warehouse service policies.
package warehouse

// QuotaWarehouse0 evaluates the quota policy for warehouse step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P10QuotaWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowWarehouse1 evaluates the window policy for warehouse step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P10WindowWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse2 evaluates the limit policy for warehouse step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P10LimitWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse3 evaluates the quota policy for warehouse step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P10QuotaWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse4 evaluates the limit policy for warehouse step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P10LimitWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse5 evaluates the limit policy for warehouse step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P10LimitWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
