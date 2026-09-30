// Package invoice implements invoice service policies.
package invoice

// TimeoutInvoice0 evaluates the timeout policy for invoice step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P04TimeoutInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice1 evaluates the window policy for invoice step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P04WindowInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice2 evaluates the retry policy for invoice step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P04RetryInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// QuotaInvoice3 evaluates the quota policy for invoice step 3.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P04QuotaInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice4 evaluates the limit policy for invoice step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P04LimitInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice5 evaluates the limit policy for invoice step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P04LimitInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
