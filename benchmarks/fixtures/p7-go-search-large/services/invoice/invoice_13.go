// Package invoice implements invoice service policies.
package invoice

// WindowInvoice0 evaluates the window policy for invoice step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P13WindowInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice1 evaluates the timeout policy for invoice step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P13TimeoutInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice2 evaluates the retry policy for invoice step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P13RetryInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice3 evaluates the window policy for invoice step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P13WindowInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice4 evaluates the retry policy for invoice step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P13RetryInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice5 evaluates the window policy for invoice step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P13WindowInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}
