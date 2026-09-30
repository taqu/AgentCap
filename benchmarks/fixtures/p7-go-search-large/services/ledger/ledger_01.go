// Package ledger implements ledger service policies.
package ledger

// LimitLedger0 evaluates the limit policy for ledger step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P01LimitLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger1 evaluates the limit policy for ledger step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P01LimitLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowLedger2 evaluates the window policy for ledger step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P01WindowLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger3 evaluates the burst policy for ledger step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P01BurstLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger4 evaluates the limit policy for ledger step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P01LimitLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLedger5 evaluates the timeout policy for ledger step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P01TimeoutLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
