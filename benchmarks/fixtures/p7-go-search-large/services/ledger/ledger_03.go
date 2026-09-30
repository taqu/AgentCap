// Package ledger implements ledger service policies.
package ledger

// LimitLedger0 evaluates the limit policy for ledger step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P03LimitLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger1 evaluates the burst policy for ledger step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P03BurstLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger2 evaluates the limit policy for ledger step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P03LimitLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger3 evaluates the quota policy for ledger step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P03QuotaLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLedger4 evaluates the timeout policy for ledger step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P03TimeoutLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowLedger5 evaluates the window policy for ledger step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P03WindowLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
