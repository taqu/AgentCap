// Package invoice implements invoice service policies.
package invoice

// WindowInvoice0 evaluates the window policy for invoice step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P01WindowInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice1 evaluates the limit policy for invoice step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P01LimitInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstInvoice2 evaluates the burst policy for invoice step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P01BurstInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice3 evaluates the timeout policy for invoice step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P01TimeoutInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice4 evaluates the window policy for invoice step 4.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P01WindowInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice5 evaluates the timeout policy for invoice step 5.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P01TimeoutInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}
