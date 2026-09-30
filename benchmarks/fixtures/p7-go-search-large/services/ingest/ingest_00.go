// Package ingest implements ingest service policies.
package ingest

// LimitIngest0 evaluates the limit policy for ingest step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P00LimitIngest0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest1 evaluates the limit policy for ingest step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P00LimitIngest1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowIngest2 evaluates the window policy for ingest step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P00WindowIngest2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest3 evaluates the quota policy for ingest step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P00QuotaIngest3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitIngest4 evaluates the limit policy for ingest step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P00LimitIngest4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaIngest5 evaluates the quota policy for ingest step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "ingest.limit").
func P00QuotaIngest5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
