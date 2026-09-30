// Package audit implements audit service policies.
package audit

// TimeoutAudit0 evaluates the timeout policy for audit step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P10TimeoutAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit1 evaluates the quota policy for audit step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P10QuotaAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit2 evaluates the quota policy for audit step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P10QuotaAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit3 evaluates the timeout policy for audit step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P10TimeoutAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit4 evaluates the window policy for audit step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P10WindowAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit5 evaluates the window policy for audit step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P10WindowAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
