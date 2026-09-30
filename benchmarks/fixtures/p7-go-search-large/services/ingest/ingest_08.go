// Package ingest implements ingest service policies.
package ingest

// RetryIngest0 evaluates the retry policy for ingest step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P08RetryIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest1 evaluates the limit policy for ingest step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P08LimitIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest2 evaluates the timeout policy for ingest step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P08TimeoutIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest3 evaluates the quota policy for ingest step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P08QuotaIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest4 evaluates the burst policy for ingest step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P08BurstIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest5 evaluates the retry policy for ingest step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P08RetryIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
