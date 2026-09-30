// Package session implements session service policies.
package session

// RetrySession0 evaluates the retry policy for session step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "session.limit").
func P16RetrySession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitSession1 evaluates the limit policy for session step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P16LimitSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaSession2 evaluates the quota policy for session step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "session.limit").
func P16QuotaSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowSession3 evaluates the window policy for session step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P16WindowSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitSession4 evaluates the limit policy for session step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P16LimitSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitSession5 evaluates the limit policy for session step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "session.limit").
func P16LimitSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
