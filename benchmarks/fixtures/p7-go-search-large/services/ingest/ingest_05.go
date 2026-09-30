// Package ingest implements ingest service policies.
package ingest

// QuotaIngest0 evaluates the quota policy for ingest step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P05QuotaIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest1 evaluates the burst policy for ingest step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P05BurstIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest2 evaluates the burst policy for ingest step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P05BurstIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest3 evaluates the quota policy for ingest step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P05QuotaIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest4 evaluates the limit policy for ingest step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P05LimitIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest5 evaluates the burst policy for ingest step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P05BurstIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
