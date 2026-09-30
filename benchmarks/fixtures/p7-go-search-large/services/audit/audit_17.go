// Package audit implements audit service policies.
package audit

// LimitAudit0 evaluates the limit policy for audit step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P17LimitAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit1 evaluates the timeout policy for audit step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P17TimeoutAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit2 evaluates the burst policy for audit step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P17BurstAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit3 evaluates the window policy for audit step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P17WindowAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit4 evaluates the timeout policy for audit step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P17TimeoutAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutAudit5 evaluates the timeout policy for audit step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P17TimeoutAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
