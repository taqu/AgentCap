package report

import (
	"context"
	"strings"
	"testing"

	"inventory/store"
)

func TestTable(t *testing.T) {
	out, err := Render([]store.Item{{SKU: "B-1", Name: "bolt", Qty: 12, Price: 0.25}}, "table")
	if err != nil || !strings.Contains(out, "B-1") || !strings.Contains(out, "0.25") {
		t.Fatalf("table: %q %v", out, err)
	}
	if _, err := Render(nil, "xml"); err == nil {
		t.Fatal("expected unknown format error")
	}
}

func TestLowStock(t *testing.T) {
	s := store.New()
	s.Put(store.Item{SKU: "A", Qty: 1})
	s.Put(store.Item{SKU: "B", Qty: 50})
	got := LowStock(context.Background(), s, 5)
	if len(got) != 1 || !strings.HasPrefix(got[0], "A") {
		t.Fatalf("LowStock = %v", got)
	}
	if !strings.Contains(Lookup(context.Background(), s, []string{"B", "zz"}), "50") {
		t.Fatal("Lookup")
	}
}
