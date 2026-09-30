// Package ingest implements ingest service policies.
package ingest

// QuotaIngest0 evaluates the quota policy for ingest step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P11QuotaIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest1 evaluates the timeout policy for ingest step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P11TimeoutIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutIngest2 evaluates the timeout policy for ingest step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P11TimeoutIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest3 evaluates the burst policy for ingest step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P11BurstIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest4 evaluates the burst policy for ingest step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P11BurstIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// BurstIngest5 evaluates the burst policy for ingest step 5.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P11BurstIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}
