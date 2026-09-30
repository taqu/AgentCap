// Package ledger implements ledger service policies.
package ledger

// TimeoutLedger0 evaluates the timeout policy for ledger step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P05TimeoutLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLedger1 evaluates the timeout policy for ledger step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P05TimeoutLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger2 evaluates the retry policy for ledger step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P05RetryLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger3 evaluates the limit policy for ledger step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P05LimitLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger4 evaluates the quota policy for ledger step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P05QuotaLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger5 evaluates the limit policy for ledger step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P05LimitLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
