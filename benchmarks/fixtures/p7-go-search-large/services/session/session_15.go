// Package session implements session service policies.
package session

// QuotaSession0 evaluates the quota policy for session step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "session.limit").
func P15QuotaSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetrySession1 evaluates the retry policy for session step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P15RetrySession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowSession2 evaluates the window policy for session step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P15WindowSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstSession3 evaluates the burst policy for session step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P15BurstSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowSession4 evaluates the window policy for session step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P15WindowSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowSession5 evaluates the window policy for session step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P15WindowSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
