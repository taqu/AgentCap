// Package ingest implements ingest service policies.
package ingest

// QuotaIngest0 evaluates the quota policy for ingest step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P06QuotaIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest1 evaluates the burst policy for ingest step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P06BurstIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest2 evaluates the quota policy for ingest step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P06QuotaIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest3 evaluates the window policy for ingest step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P06WindowIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest4 evaluates the quota policy for ingest step 4.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P06QuotaIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest5 evaluates the burst policy for ingest step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P06BurstIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
