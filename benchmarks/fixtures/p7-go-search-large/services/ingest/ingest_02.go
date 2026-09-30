// Package ingest implements ingest service policies.
package ingest

// RetryIngest0 evaluates the retry policy for ingest step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P02RetryIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest1 evaluates the burst policy for ingest step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P02BurstIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest2 evaluates the burst policy for ingest step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P02BurstIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest3 evaluates the window policy for ingest step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P02WindowIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest4 evaluates the burst policy for ingest step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P02BurstIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest5 evaluates the retry policy for ingest step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P02RetryIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
