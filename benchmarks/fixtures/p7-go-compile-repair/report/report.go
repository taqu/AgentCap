// Package report renders inventory listings.
package report

import (
	"context"
	"fmt"
	"strings"

	"inventory/store"
	"inventory/util"
)

// Render formats items. Supported formats: "table".
func Render(items []store.Item, format string) (string, error) {
	switch format {
	case "table", "":
		return table(items), nil
	}
	return "", fmt.Errorf("unknown format %q", format)
}

func table(items []store.Item) string {
	var b strings.Builder
	b.WriteString(util.PadRight("SKU", 8) + util.PadRight("NAME", 20) + util.PadRight("QTY", 6) + "PRICE\n")
	for _, it := range items {
		fmt.Fprintf(&b, "%s%s%s%.2f\n", util.PadRight(it.SKU, 8), util.PadRight(util.Truncate(it.Name, 19), 20), util.PadRight(fmt.Sprint(it.Qty), 6), it.Price)
	}
	return b.String()
}

// Lookup renders the named SKUs, skipping unknown ones.
func Lookup(ctx context.Context, s *store.Store, skus []string) string {
	var items []store.Item
	for _, sku := range skus {
		if it, err := s.Get(sku); err == nil {
			items = append(items, it)
		}
	}
	return table(items)
}

// LowStock lists items whose quantity is at most threshold.
func LowStock(ctx context.Context, s *store.Store, threshold int) []string {
	var out []string
	for _, it := range s.List() {
		cur, err := s.Get(it.SKU)
		if err == nil && cur.Qty <= threshold {
			out = append(out, fmt.Sprintf("%s (%d left, off by %d)", cur.SKU, cur.Qty, util.Abs(threshold-cur.Qty)))
		}
	}
	return out
}
