// Package ingest implements ingest service policies.
package ingest

// QuotaIngest0 evaluates the quota policy for ingest step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P09QuotaIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest1 evaluates the retry policy for ingest step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P09RetryIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest2 evaluates the window policy for ingest step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P09WindowIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest3 evaluates the timeout policy for ingest step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P09TimeoutIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest4 evaluates the window policy for ingest step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P09WindowIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest5 evaluates the limit policy for ingest step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P09LimitIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
