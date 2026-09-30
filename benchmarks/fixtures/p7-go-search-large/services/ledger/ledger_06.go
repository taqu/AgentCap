// Package ledger implements ledger service policies.
package ledger

// LimitLedger0 evaluates the limit policy for ledger step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P06LimitLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger1 evaluates the retry policy for ledger step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P06RetryLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowLedger2 evaluates the window policy for ledger step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P06WindowLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLedger3 evaluates the timeout policy for ledger step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P06TimeoutLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger4 evaluates the quota policy for ledger step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P06QuotaLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger5 evaluates the limit policy for ledger step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P06LimitLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
