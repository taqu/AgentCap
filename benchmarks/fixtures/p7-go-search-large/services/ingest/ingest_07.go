// Package ingest implements ingest service policies.
package ingest

// WindowIngest0 evaluates the window policy for ingest step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P07WindowIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest1 evaluates the quota policy for ingest step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P07QuotaIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest2 evaluates the window policy for ingest step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P07WindowIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest3 evaluates the retry policy for ingest step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P07RetryIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest4 evaluates the limit policy for ingest step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P07LimitIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest5 evaluates the limit policy for ingest step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P07LimitIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
