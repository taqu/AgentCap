// Package audit implements audit service policies.
package audit

// RetryAudit0 evaluates the retry policy for audit step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P06RetryAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit1 evaluates the quota policy for audit step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P06QuotaAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit2 evaluates the limit policy for audit step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P06LimitAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit3 evaluates the window policy for audit step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P06WindowAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryAudit4 evaluates the retry policy for audit step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P06RetryAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit5 evaluates the timeout policy for audit step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P06TimeoutAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
