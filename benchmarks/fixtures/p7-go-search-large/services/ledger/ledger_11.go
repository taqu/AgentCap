// Package ledger implements ledger service policies.
package ledger

// RetryLedger0 evaluates the retry policy for ledger step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P11RetryLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger1 evaluates the limit policy for ledger step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P11LimitLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger2 evaluates the retry policy for ledger step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P11RetryLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger3 evaluates the limit policy for ledger step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P11LimitLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger4 evaluates the limit policy for ledger step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P11LimitLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger5 evaluates the burst policy for ledger step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P11BurstLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
