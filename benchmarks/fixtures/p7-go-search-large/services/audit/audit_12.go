// Package audit implements audit service policies.
package audit

// RetryAudit0 evaluates the retry policy for audit step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P12RetryAudit0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit1 evaluates the burst policy for audit step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P12BurstAudit1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaAudit2 evaluates the quota policy for audit step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P12QuotaAudit2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryAudit3 evaluates the retry policy for audit step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P12RetryAudit3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstAudit4 evaluates the burst policy for audit step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P12BurstAudit4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitAudit5 evaluates the limit policy for audit step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "audit.limit").
func P12LimitAudit5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
