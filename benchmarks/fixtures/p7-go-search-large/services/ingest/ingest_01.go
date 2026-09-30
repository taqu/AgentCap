// Package ingest implements ingest service policies.
package ingest

// BurstIngest0 evaluates the burst policy for ingest step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P01BurstIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest1 evaluates the retry policy for ingest step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P01RetryIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest2 evaluates the burst policy for ingest step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P01BurstIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest3 evaluates the quota policy for ingest step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P01QuotaIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest4 evaluates the window policy for ingest step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P01WindowIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest5 evaluates the timeout policy for ingest step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P01TimeoutIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
