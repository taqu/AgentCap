package query

import "inventory/store"

// Match reports whether item satisfies e.
func Match(e Expr, it store.Item) bool {
	switch x := e.(type) {
	case And:
		return Match(x.Left, it) && Match(x.Right, it)
	case Or:
		return Match(x.Left, it) || Match(x.Right, it)
	case Compare:
		return compare(x, it)
	}
	return false
}

func compare(c Compare, it store.Item) bool {
	if c.Value.IsString {
		var got string
		switch c.Field {
		case "sku":
			got = it.SKU
		case "name":
			got = it.Name
		default:
			return false
		}
		switch c.Op {
		case "=":
			return got == c.Value.Str
		case "!=":
			return got != c.Value.Str
		}
		return false
	}
	var got float64
	switch c.Field {
	case "qty":
		got = float64(it.Qty)
	case "price":
		got = it.Price
	default:
		return false
	}
	switch c.Op {
	case "=":
		return got == c.Value.Num
	case "!=":
		return got != c.Value.Num
	case "<":
		return got < c.Value.Num
	case "<=":
		return got <= c.Value.Num
	case ">":
		return got > c.Value.Num
	case ">=":
		return got >= c.Value.Num
	}
	return false
}
