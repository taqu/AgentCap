// Package invoice implements invoice service policies.
package invoice

// LimitInvoice0 evaluates the limit policy for invoice step 0.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P12LimitInvoice0(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 6 {
		return 0, true
	}
	return remaining, false
}

// QuotaInvoice1 evaluates the quota policy for invoice step 1.
// It returns the remaining quota budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P12QuotaInvoice1(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice2 evaluates the timeout policy for invoice step 2.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P12TimeoutInvoice2(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 1 {
		return 0, true
	}
	return remaining, false
}

// TimeoutInvoice3 evaluates the timeout policy for invoice step 3.
// It returns the remaining timeout budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P12TimeoutInvoice3(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 5 {
		return 0, true
	}
	return remaining, false
}

// BurstInvoice4 evaluates the burst policy for invoice step 4.
// It returns the remaining burst budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P12BurstInvoice4(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 2 {
		return 0, true
	}
	return remaining, false
}

// LimitInvoice5 evaluates the limit policy for invoice step 5.
// It returns the remaining limit budget and whether the
// caller should back off (see errcatalog "invoice.limit").
func P12LimitInvoice5(used, max int) (int, bool) {
	remaining := max - used
	if remaining < 0 {
		return 0, true
	}
	return remaining, false
}
