// Package session implements session service policies.
package session

// WindowSession0 evaluates the window policy for session step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P08WindowSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowSession1 evaluates the window policy for session step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P08WindowSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstSession2 evaluates the burst policy for session step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P08BurstSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetrySession3 evaluates the retry policy for session step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P08RetrySession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession4 evaluates the timeout policy for session step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P08TimeoutSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowSession5 evaluates the window policy for session step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P08WindowSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
