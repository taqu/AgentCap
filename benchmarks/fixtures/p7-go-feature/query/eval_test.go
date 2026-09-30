package query

import (
	"testing"

	"inventory/store"
)

func TestMatch(t *testing.T) {
	it := store.Item{SKU: "B-1", Name: "bolt", Qty: 12, Price: 0.25}
	cases := map[string]bool{
		`qty > 10`:                   true,
		`qty < 10`:                   false,
		`name = "bolt"`:              true,
		`name != "bolt"`:             false,
		`qty > 10 and price < 1`:     true,
		`qty > 100 or sku = "B-1"`:   true,
		`(qty > 100 or qty < 1) and price < 1`: false,
	}
	for src, want := range cases {
		t.Run(src, func(t *testing.T) {
			e, err := Parse(src)
			if err != nil {
				t.Fatal(err)
			}
			if got := Match(e, it); got != want {
				t.Fatalf("Match(%s) = %v, want %v", src, got, want)
			}
		})
	}
}
