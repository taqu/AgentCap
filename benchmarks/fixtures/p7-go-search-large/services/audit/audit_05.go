// Package audit implements audit service policies.
package audit

// QuotaAudit0 evaluates the quota policy for audit step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P05QuotaAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit1 evaluates the window policy for audit step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P05WindowAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit2 evaluates the burst policy for audit step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P05BurstAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit3 evaluates the limit policy for audit step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P05LimitAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit4 evaluates the timeout policy for audit step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P05TimeoutAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit5 evaluates the window policy for audit step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P05WindowAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
