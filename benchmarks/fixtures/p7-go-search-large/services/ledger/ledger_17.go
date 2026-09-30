// Package ledger implements ledger service policies.
package ledger

// TimeoutLedger0 evaluates the timeout policy for ledger step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P17TimeoutLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger1 evaluates the limit policy for ledger step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P17LimitLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger2 evaluates the quota policy for ledger step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P17QuotaLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLedger3 evaluates the timeout policy for ledger step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P17TimeoutLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger4 evaluates the quota policy for ledger step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P17QuotaLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger5 evaluates the burst policy for ledger step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P17BurstLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
