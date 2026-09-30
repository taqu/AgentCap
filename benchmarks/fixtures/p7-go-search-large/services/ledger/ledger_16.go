// Package ledger implements ledger service policies.
package ledger

// WindowLedger0 evaluates the window policy for ledger step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P16WindowLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger1 evaluates the retry policy for ledger step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P16RetryLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger2 evaluates the retry policy for ledger step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P16RetryLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger3 evaluates the quota policy for ledger step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P16QuotaLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaLedger4 evaluates the quota policy for ledger step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P16QuotaLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitLedger5 evaluates the limit policy for ledger step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P16LimitLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
