// Package ledger implements ledger service policies.
package ledger

// LimitLedger0 evaluates the limit policy for ledger step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P00LimitLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger1 evaluates the retry policy for ledger step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P00RetryLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger2 evaluates the limit policy for ledger step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P00LimitLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger3 evaluates the burst policy for ledger step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P00BurstLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger4 evaluates the burst policy for ledger step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P00BurstLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger5 evaluates the retry policy for ledger step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P00RetryLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
