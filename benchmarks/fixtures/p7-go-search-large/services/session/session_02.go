// Package session implements session service policies.
package session

// QuotaSession0 evaluates the quota policy for session step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "session.limit").
func P02QuotaSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowSession1 evaluates the window policy for session step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P02WindowSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession2 evaluates the timeout policy for session step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P02TimeoutSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitSession3 evaluates the limit policy for session step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P02LimitSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutSession4 evaluates the timeout policy for session step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "session.limit").
func P02TimeoutSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetrySession5 evaluates the retry policy for session step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P02RetrySession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
