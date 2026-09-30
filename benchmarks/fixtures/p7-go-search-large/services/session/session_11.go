// Package session implements session service policies.
package session

// BurstSession0 evaluates the burst policy for session step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P11BurstSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// QuotaSession1 evaluates the quota policy for session step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "session.limit").
func P11QuotaSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession2 evaluates the timeout policy for session step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P11TimeoutSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitSession3 evaluates the limit policy for session step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P11LimitSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowSession4 evaluates the window policy for session step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P11WindowSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitSession5 evaluates the limit policy for session step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P11LimitSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
