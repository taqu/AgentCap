// Package ingest implements ingest service policies.
package ingest

// QuotaIngest0 evaluates the quota policy for ingest step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P16QuotaIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest1 evaluates the quota policy for ingest step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P16QuotaIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest2 evaluates the retry policy for ingest step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P16RetryIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest3 evaluates the limit policy for ingest step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P16LimitIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryIngest4 evaluates the retry policy for ingest step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P16RetryIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest5 evaluates the quota policy for ingest step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P16QuotaIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
