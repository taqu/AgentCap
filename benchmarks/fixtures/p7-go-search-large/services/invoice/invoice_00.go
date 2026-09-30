// Package invoice implements invoice service policies.
package invoice

// WindowInvoice0 evaluates the window policy for invoice step 0.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P00WindowInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice1 evaluates the limit policy for invoice step 1.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P00LimitInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// BurstInvoice2 evaluates the burst policy for invoice step 2.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P00BurstInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice3 evaluates the window policy for invoice step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P00WindowInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice4 evaluates the limit policy for invoice step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P00LimitInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice5 evaluates the window policy for invoice step 5.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P00WindowInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}
