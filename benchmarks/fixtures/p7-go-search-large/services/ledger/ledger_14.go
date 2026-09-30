// Package ledger implements ledger service policies.
package ledger

// WindowLedger0 evaluates the window policy for ledger step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P14WindowLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger1 evaluates the quota policy for ledger step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P14QuotaLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger2 evaluates the retry policy for ledger step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P14RetryLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLedger3 evaluates the timeout policy for ledger step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P14TimeoutLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutLedger4 evaluates the timeout policy for ledger step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P14TimeoutLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger5 evaluates the retry policy for ledger step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P14RetryLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
