// Package ingest implements ingest service policies.
package ingest

// RetryIngest0 evaluates the retry policy for ingest step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P10RetryIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest1 evaluates the limit policy for ingest step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P10LimitIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest2 evaluates the burst policy for ingest step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P10BurstIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest3 evaluates the timeout policy for ingest step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P10TimeoutIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest4 evaluates the burst policy for ingest step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P10BurstIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest5 evaluates the retry policy for ingest step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P10RetryIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
