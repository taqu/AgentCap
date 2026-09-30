// Package audit implements audit service policies.
package audit

// BurstAudit0 evaluates the burst policy for audit step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P08BurstAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit1 evaluates the limit policy for audit step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P08LimitAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit2 evaluates the quota policy for audit step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P08QuotaAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit3 evaluates the timeout policy for audit step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P08TimeoutAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryAudit4 evaluates the retry policy for audit step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P08RetryAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit5 evaluates the burst policy for audit step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P08BurstAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
