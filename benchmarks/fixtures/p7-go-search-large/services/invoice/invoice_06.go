// Package invoice implements invoice service policies.
package invoice

// QuotaInvoice0 evaluates the quota policy for invoice step 0.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P06QuotaInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice1 evaluates the window policy for invoice step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P06WindowInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// QuotaInvoice2 evaluates the quota policy for invoice step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P06QuotaInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstInvoice3 evaluates the burst policy for invoice step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P06BurstInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice4 evaluates the retry policy for invoice step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P06RetryInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaInvoice5 evaluates the quota policy for invoice step 5.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P06QuotaInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
