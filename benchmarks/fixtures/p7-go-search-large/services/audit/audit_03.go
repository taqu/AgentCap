// Package audit implements audit service policies.
package audit

// RetryAudit0 evaluates the retry policy for audit step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P03RetryAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryAudit1 evaluates the retry policy for audit step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P03RetryAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit2 evaluates the quota policy for audit step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P03QuotaAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit3 evaluates the quota policy for audit step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P03QuotaAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit4 evaluates the quota policy for audit step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P03QuotaAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit5 evaluates the burst policy for audit step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P03BurstAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
