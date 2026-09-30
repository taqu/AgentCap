// Package session implements session service policies.
package session

// WindowSession0 evaluates the window policy for session step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P04WindowSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetrySession1 evaluates the retry policy for session step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P04RetrySession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstSession2 evaluates the burst policy for session step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P04BurstSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession3 evaluates the timeout policy for session step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P04TimeoutSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession4 evaluates the timeout policy for session step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P04TimeoutSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetrySession5 evaluates the retry policy for session step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P04RetrySession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
