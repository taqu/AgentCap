// Package session implements session service policies.
package session

// RetrySession0 evaluates the retry policy for session step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P01RetrySession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstSession1 evaluates the burst policy for session step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P01BurstSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowSession2 evaluates the window policy for session step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P01WindowSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstSession3 evaluates the burst policy for session step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P01BurstSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetrySession4 evaluates the retry policy for session step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P01RetrySession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstSession5 evaluates the burst policy for session step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P01BurstSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
