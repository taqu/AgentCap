// Package audit implements audit service policies.
package audit

// TimeoutAudit0 evaluates the timeout policy for audit step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P13TimeoutAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit1 evaluates the quota policy for audit step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P13QuotaAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit2 evaluates the burst policy for audit step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P13BurstAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit3 evaluates the limit policy for audit step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P13LimitAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit4 evaluates the burst policy for audit step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P13BurstAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit5 evaluates the timeout policy for audit step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P13TimeoutAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
