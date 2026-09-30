// Package ingest implements ingest service policies.
package ingest

// BurstIngest0 evaluates the burst policy for ingest step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P17BurstIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest1 evaluates the limit policy for ingest step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P17LimitIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest2 evaluates the window policy for ingest step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P17WindowIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest3 evaluates the retry policy for ingest step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P17RetryIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest4 evaluates the timeout policy for ingest step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P17TimeoutIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest5 evaluates the burst policy for ingest step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P17BurstIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
