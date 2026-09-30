// Package audit implements audit service policies.
package audit

// QuotaAudit0 evaluates the quota policy for audit step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P11QuotaAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryAudit1 evaluates the retry policy for audit step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P11RetryAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit2 evaluates the window policy for audit step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P11WindowAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryAudit3 evaluates the retry policy for audit step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P11RetryAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit4 evaluates the quota policy for audit step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P11QuotaAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit5 evaluates the window policy for audit step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P11WindowAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
