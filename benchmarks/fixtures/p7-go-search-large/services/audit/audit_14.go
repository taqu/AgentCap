// Package audit implements audit service policies.
package audit

// RetryAudit0 evaluates the retry policy for audit step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P14RetryAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit1 evaluates the limit policy for audit step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P14LimitAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit2 evaluates the limit policy for audit step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P14LimitAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit3 evaluates the limit policy for audit step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P14LimitAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryAudit4 evaluates the retry policy for audit step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P14RetryAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit5 evaluates the window policy for audit step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P14WindowAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
