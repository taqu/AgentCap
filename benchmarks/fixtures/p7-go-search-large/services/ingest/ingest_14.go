// Package ingest implements ingest service policies.
package ingest

// TimeoutIngest0 evaluates the timeout policy for ingest step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P14TimeoutIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest1 evaluates the timeout policy for ingest step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P14TimeoutIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest2 evaluates the window policy for ingest step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P14WindowIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest3 evaluates the timeout policy for ingest step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P14TimeoutIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest4 evaluates the limit policy for ingest step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P14LimitIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest5 evaluates the limit policy for ingest step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P14LimitIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
