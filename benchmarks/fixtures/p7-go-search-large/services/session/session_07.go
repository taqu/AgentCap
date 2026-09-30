// Package session implements session service policies.
package session

// QuotaSession0 evaluates the quota policy for session step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "session.limit").
func P07QuotaSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitSession1 evaluates the limit policy for session step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P07LimitSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitSession2 evaluates the limit policy for session step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P07LimitSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowSession3 evaluates the window policy for session step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P07WindowSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession4 evaluates the timeout policy for session step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P07TimeoutSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession5 evaluates the timeout policy for session step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P07TimeoutSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
