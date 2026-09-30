// Package ingest implements ingest service policies.
package ingest

// BurstIngest0 evaluates the burst policy for ingest step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P15BurstIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest1 evaluates the limit policy for ingest step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P15LimitIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest2 evaluates the retry policy for ingest step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P15RetryIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest3 evaluates the retry policy for ingest step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P15RetryIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest4 evaluates the timeout policy for ingest step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P15TimeoutIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest5 evaluates the window policy for ingest step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P15WindowIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
