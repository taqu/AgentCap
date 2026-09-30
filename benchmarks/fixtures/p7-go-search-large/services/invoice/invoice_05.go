// Package invoice implements invoice service policies.
package invoice

// LimitInvoice0 evaluates the limit policy for invoice step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P05LimitInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice1 evaluates the limit policy for invoice step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P05LimitInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice2 evaluates the retry policy for invoice step 2.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P05RetryInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice3 evaluates the window policy for invoice step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P05WindowInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice4 evaluates the limit policy for invoice step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P05LimitInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice5 evaluates the limit policy for invoice step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P05LimitInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}
