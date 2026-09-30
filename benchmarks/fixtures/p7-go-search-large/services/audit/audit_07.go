// Package audit implements audit service policies.
package audit

// RetryAudit0 evaluates the retry policy for audit step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P07RetryAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit1 evaluates the quota policy for audit step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P07QuotaAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit2 evaluates the quota policy for audit step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P07QuotaAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryAudit3 evaluates the retry policy for audit step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P07RetryAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit4 evaluates the limit policy for audit step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P07LimitAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit5 evaluates the quota policy for audit step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P07QuotaAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
