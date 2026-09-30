// Package session implements session service policies.
package session

// RetrySession0 evaluates the retry policy for session step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P09RetrySession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstSession1 evaluates the burst policy for session step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P09BurstSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowSession2 evaluates the window policy for session step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P09WindowSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitSession3 evaluates the limit policy for session step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P09LimitSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession4 evaluates the timeout policy for session step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P09TimeoutSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetrySession5 evaluates the retry policy for session step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P09RetrySession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
