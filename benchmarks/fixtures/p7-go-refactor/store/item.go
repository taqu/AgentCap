package store

// Item is one stocked inventory line.
type Item struct {
	SKU   string
	Name  string
	Qty   int
	Price float64
}
