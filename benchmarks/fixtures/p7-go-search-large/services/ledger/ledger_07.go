// Package ledger implements ledger service policies.
package ledger

// LimitLedger0 evaluates the limit policy for ledger step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P07LimitLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger1 evaluates the burst policy for ledger step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P07BurstLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLedger2 evaluates the timeout policy for ledger step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P07TimeoutLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger3 evaluates the quota policy for ledger step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P07QuotaLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger4 evaluates the burst policy for ledger step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P07BurstLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLedger5 evaluates the timeout policy for ledger step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P07TimeoutLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
