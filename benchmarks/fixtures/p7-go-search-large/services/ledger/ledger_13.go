// Package ledger implements ledger service policies.
package ledger

// WindowLedger0 evaluates the window policy for ledger step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P13WindowLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger1 evaluates the quota policy for ledger step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P13QuotaLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger2 evaluates the limit policy for ledger step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P13LimitLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger3 evaluates the retry policy for ledger step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P13RetryLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger4 evaluates the quota policy for ledger step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P13QuotaLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger5 evaluates the retry policy for ledger step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P13RetryLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
