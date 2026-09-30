// Package ingest implements ingest service policies.
package ingest

// WindowIngest0 evaluates the window policy for ingest step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P04WindowIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest1 evaluates the quota policy for ingest step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P04QuotaIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest2 evaluates the quota policy for ingest step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P04QuotaIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest3 evaluates the timeout policy for ingest step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P04TimeoutIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest4 evaluates the limit policy for ingest step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P04LimitIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest5 evaluates the burst policy for ingest step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P04BurstIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
