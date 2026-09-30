// Package audit implements audit service policies.
package audit

// QuotaAudit0 evaluates the quota policy for audit step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P16QuotaAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit1 evaluates the timeout policy for audit step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P16TimeoutAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit2 evaluates the limit policy for audit step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P16LimitAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit3 evaluates the quota policy for audit step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P16QuotaAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit4 evaluates the limit policy for audit step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P16LimitAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit5 evaluates the limit policy for audit step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P16LimitAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
