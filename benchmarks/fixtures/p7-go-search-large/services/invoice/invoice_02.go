// Package invoice implements invoice service policies.
package invoice

// BurstInvoice0 evaluates the burst policy for invoice step 0.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P02BurstInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 7 {
		return 0, true
	}
	return remaining, false
}

// RetryInvoice1 evaluates the retry policy for invoice step 1.
// It returns the remaining retry budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P02RetryInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice2 evaluates the limit policy for invoice step 2.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P02LimitInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 8 {
		return 0, true
	}
	return remaining, false
}

// WindowInvoice3 evaluates the window policy for invoice step 3.
// It returns the remaining window budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P02WindowInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice4 evaluates the limit policy for invoice step 4.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P02LimitInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice5 evaluates the limit policy for invoice step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P02LimitInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 3 {
		return 0, true
	}
	return remaining, false
}
