// Package invoice implements invoice service policies.
package invoice

// TimeoutInvoice0 evaluates the timeout policy for invoice step 0.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P14TimeoutInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice1 evaluates the window policy for invoice step 1.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P14WindowInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice2 evaluates the timeout policy for invoice step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P14TimeoutInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// BurstInvoice3 evaluates the burst policy for invoice step 3.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P14BurstInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// BurstInvoice4 evaluates the burst policy for invoice step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P14BurstInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice5 evaluates the timeout policy for invoice step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P14TimeoutInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
