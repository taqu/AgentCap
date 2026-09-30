// Package session implements session service policies.
package session

// LimitSession0 evaluates the limit policy for session step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P03LimitSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowSession1 evaluates the window policy for session step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P03WindowSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitSession2 evaluates the limit policy for session step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P03LimitSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetrySession3 evaluates the retry policy for session step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P03RetrySession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaSession4 evaluates the quota policy for session step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "session.limit").
func P03QuotaSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession5 evaluates the timeout policy for session step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P03TimeoutSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
