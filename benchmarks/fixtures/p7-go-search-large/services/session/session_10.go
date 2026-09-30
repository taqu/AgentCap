// Package session implements session service policies.
package session

// WindowSession0 evaluates the window policy for session step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P10WindowSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaSession1 evaluates the quota policy for session step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "session.limit").
func P10QuotaSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetrySession2 evaluates the retry policy for session step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P10RetrySession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetrySession3 evaluates the retry policy for session step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P10RetrySession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowSession4 evaluates the window policy for session step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P10WindowSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstSession5 evaluates the burst policy for session step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P10BurstSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
