// Package session implements session service policies.
package session

// TimeoutSession0 evaluates the timeout policy for session step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P00TimeoutSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitSession1 evaluates the limit policy for session step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P00LimitSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowSession2 evaluates the window policy for session step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P00WindowSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowSession3 evaluates the window policy for session step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P00WindowSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetrySession4 evaluates the retry policy for session step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P00RetrySession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowSession5 evaluates the window policy for session step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P00WindowSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
