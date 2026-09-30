// Package invoice implements invoice service policies.
package invoice

// RetryInvoice0 evaluates the retry policy for invoice step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P17RetryInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice1 evaluates the timeout policy for invoice step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P17TimeoutInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// QuotaInvoice2 evaluates the quota policy for invoice step 2.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P17QuotaInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice3 evaluates the window policy for invoice step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P17WindowInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice4 evaluates the retry policy for invoice step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P17RetryInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice5 evaluates the timeout policy for invoice step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P17TimeoutInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
