// Package ledger implements ledger service policies.
package ledger

// BurstLedger0 evaluates the burst policy for ledger step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P08BurstLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger1 evaluates the quota policy for ledger step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P08QuotaLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger2 evaluates the retry policy for ledger step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P08RetryLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger3 evaluates the limit policy for ledger step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P08LimitLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger4 evaluates the limit policy for ledger step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P08LimitLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger5 evaluates the burst policy for ledger step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P08BurstLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
