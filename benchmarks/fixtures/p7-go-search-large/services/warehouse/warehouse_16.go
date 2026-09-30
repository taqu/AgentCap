// Package warehouse implements warehouse service policies.
package warehouse

// RetryWarehouse0 evaluates the retry policy for warehouse step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P16RetryWarehouse0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitWarehouse1 evaluates the limit policy for warehouse step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P16LimitWarehouse1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryWarehouse2 evaluates the retry policy for warehouse step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P16RetryWarehouse2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowWarehouse3 evaluates the window policy for warehouse step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P16WindowWarehouse3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaWarehouse4 evaluates the quota policy for warehouse step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P16QuotaWarehouse4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowWarehouse5 evaluates the window policy for warehouse step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "warehouse.limit").
func P16WindowWarehouse5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
