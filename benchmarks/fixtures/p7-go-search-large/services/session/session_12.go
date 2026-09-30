// Package session implements session service policies.
package session

// BurstSession0 evaluates the burst policy for session step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P12BurstSession0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowSession1 evaluates the window policy for session step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P12WindowSession1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstSession2 evaluates the burst policy for session step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "session.limit").
func P12BurstSession2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowSession3 evaluates the window policy for session step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P12WindowSession3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowSession4 evaluates the window policy for session step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P12WindowSession4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowSession5 evaluates the window policy for session step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "session.limit").
func P12WindowSession5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
