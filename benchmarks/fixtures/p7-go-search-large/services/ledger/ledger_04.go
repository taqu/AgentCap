// Package ledger implements ledger service policies.
package ledger

// QuotaLedger0 evaluates the quota policy for ledger step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P04QuotaLedger0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger1 evaluates the retry policy for ledger step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P04RetryLedger1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryLedger2 evaluates the retry policy for ledger step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P04RetryLedger2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger3 evaluates the burst policy for ledger step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P04BurstLedger3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger4 evaluates the burst policy for ledger step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P04BurstLedger4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstLedger5 evaluates the burst policy for ledger step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ledger.limit").
func P04BurstLedger5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
