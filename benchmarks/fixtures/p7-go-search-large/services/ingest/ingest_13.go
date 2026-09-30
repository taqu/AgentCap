// Package ingest implements ingest service policies.
package ingest

// BurstIngest0 evaluates the burst policy for ingest step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P13BurstIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest1 evaluates the burst policy for ingest step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P13BurstIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest2 evaluates the limit policy for ingest step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P13LimitIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest3 evaluates the burst policy for ingest step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P13BurstIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest4 evaluates the timeout policy for ingest step 4.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P13TimeoutIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest5 evaluates the quota policy for ingest step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P13QuotaIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
