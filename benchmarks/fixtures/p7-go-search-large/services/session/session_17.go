// Package session implements session service policies.
package session

// WindowSession0 evaluates the window policy for session step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P17WindowSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetrySession1 evaluates the retry policy for session step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P17RetrySession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetrySession2 evaluates the retry policy for session step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P17RetrySession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstSession3 evaluates the burst policy for session step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P17BurstSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstSession4 evaluates the burst policy for session step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P17BurstSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession5 evaluates the timeout policy for session step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P17TimeoutSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
