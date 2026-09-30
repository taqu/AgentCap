// Package invoice implements invoice service policies.
package invoice

// BurstInvoice0 evaluates the burst policy for invoice step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P07BurstInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// BurstInvoice1 evaluates the burst policy for invoice step 1.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P07BurstInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice2 evaluates the window policy for invoice step 2.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P07WindowInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice3 evaluates the window policy for invoice step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P07WindowInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice4 evaluates the retry policy for invoice step 4.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P07RetryInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice5 evaluates the timeout policy for invoice step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P07TimeoutInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}
