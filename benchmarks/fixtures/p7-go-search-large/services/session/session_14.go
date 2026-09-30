// Package session implements session service policies.
package session

// LimitSession0 evaluates the limit policy for session step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P14LimitSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession1 evaluates the timeout policy for session step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P14TimeoutSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitSession2 evaluates the limit policy for session step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P14LimitSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetrySession3 evaluates the retry policy for session step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P14RetrySession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowSession4 evaluates the window policy for session step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P14WindowSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstSession5 evaluates the burst policy for session step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P14BurstSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
