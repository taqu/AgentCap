// Package audit implements audit service policies.
package audit

// BurstAudit0 evaluates the burst policy for audit step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P04BurstAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit1 evaluates the limit policy for audit step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P04LimitAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit2 evaluates the limit policy for audit step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P04LimitAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit3 evaluates the quota policy for audit step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P04QuotaAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowAudit4 evaluates the window policy for audit step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P04WindowAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit5 evaluates the burst policy for audit step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P04BurstAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
