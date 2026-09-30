// Package ingest implements ingest service policies.
package ingest

// BurstIngest0 evaluates the burst policy for ingest step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P12BurstIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest1 evaluates the limit policy for ingest step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P12LimitIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest2 evaluates the burst policy for ingest step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P12BurstIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest3 evaluates the limit policy for ingest step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P12LimitIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest4 evaluates the quota policy for ingest step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P12QuotaIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest5 evaluates the timeout policy for ingest step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P12TimeoutIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
