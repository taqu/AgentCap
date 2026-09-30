package query

import "fmt"

// Expr is a boolean expression over item fields.
type Expr interface{ String() string }

type Compare struct {
	Field string
	Op    string
	Value Value
}

type Value struct {
	IsString bool
	Str      string
	Num      float64
}

type And struct{ Left, Right Expr }
type Or struct{ Left, Right Expr }

func (c Compare) String() string {
	if c.Value.IsString {
		return fmt.Sprintf("(%s %s %q)", c.Field, c.Op, c.Value.Str)
	}
	return fmt.Sprintf("(%s %s %g)", c.Field, c.Op, c.Value.Num)
}
func (a And) String() string { return "(" + a.Left.String() + " and " + a.Right.String() + ")" }
func (o Or) String() string  { return "(" + o.Left.String() + " or " + o.Right.String() + ")" }
