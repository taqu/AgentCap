// Package ledger implements ledger service policies.
package ledger

// TimeoutLedger0 evaluates the timeout policy for ledger step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P15TimeoutLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowLedger1 evaluates the window policy for ledger step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P15WindowLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger2 evaluates the limit policy for ledger step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P15LimitLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger3 evaluates the burst policy for ledger step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P15BurstLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger4 evaluates the limit policy for ledger step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P15LimitLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowLedger5 evaluates the window policy for ledger step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P15WindowLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
