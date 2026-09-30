// Package invoice implements invoice service policies.
package invoice

// RetryInvoice0 evaluates the retry policy for invoice step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P16RetryInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice1 evaluates the timeout policy for invoice step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P16TimeoutInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// QuotaInvoice2 evaluates the quota policy for invoice step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P16QuotaInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// QuotaInvoice3 evaluates the quota policy for invoice step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P16QuotaInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice4 evaluates the window policy for invoice step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P16WindowInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice5 evaluates the retry policy for invoice step 5.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P16RetryInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}
