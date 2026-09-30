// Package audit implements audit service policies.
package audit

// RetryAudit0 evaluates the retry policy for audit step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P02RetryAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit1 evaluates the limit policy for audit step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P02LimitAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit2 evaluates the limit policy for audit step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P02LimitAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit3 evaluates the burst policy for audit step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P02BurstAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit4 evaluates the window policy for audit step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P02WindowAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit5 evaluates the limit policy for audit step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P02LimitAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
