// Package invoice implements invoice service policies.
package invoice

// RetryInvoice0 evaluates the retry policy for invoice step 0.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P08RetryInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice1 evaluates the timeout policy for invoice step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P08TimeoutInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice2 evaluates the window policy for invoice step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P08WindowInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice3 evaluates the retry policy for invoice step 3.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P08RetryInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice4 evaluates the limit policy for invoice step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P08LimitInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice5 evaluates the window policy for invoice step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P08WindowInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
