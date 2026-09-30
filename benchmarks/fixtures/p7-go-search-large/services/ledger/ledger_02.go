// Package ledger implements ledger service policies.
package ledger

// TimeoutLedger0 evaluates the timeout policy for ledger step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P02TimeoutLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger1 evaluates the burst policy for ledger step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P02BurstLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger2 evaluates the quota policy for ledger step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P02QuotaLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger3 evaluates the retry policy for ledger step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P02RetryLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger4 evaluates the limit policy for ledger step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P02LimitLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger5 evaluates the burst policy for ledger step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P02BurstLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
