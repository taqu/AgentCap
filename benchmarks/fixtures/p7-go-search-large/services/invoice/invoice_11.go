// Package invoice implements invoice service policies.
package invoice

// BurstInvoice0 evaluates the burst policy for invoice step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P11BurstInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice1 evaluates the timeout policy for invoice step 1.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P11TimeoutInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice2 evaluates the timeout policy for invoice step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P11TimeoutInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice3 evaluates the limit policy for invoice step 3.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P11LimitInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 9 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice4 evaluates the limit policy for invoice step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P11LimitInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 4 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice5 evaluates the limit policy for invoice step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P11LimitInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}
