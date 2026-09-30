// Package audit implements audit service policies.
package audit

// RetryAudit0 evaluates the retry policy for audit step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P15RetryAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit1 evaluates the timeout policy for audit step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P15TimeoutAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit2 evaluates the limit policy for audit step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P15LimitAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit3 evaluates the burst policy for audit step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P15BurstAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit4 evaluates the quota policy for audit step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P15QuotaAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit5 evaluates the window policy for audit step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P15WindowAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
