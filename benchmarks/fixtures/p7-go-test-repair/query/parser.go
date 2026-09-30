// Package query parses inventory filter expressions.
package query

import (
	"fmt"
	"strconv"

	"inventory/lexer"
)

type parser struct {
	toks []lexer.Token
	pos  int
}

// Parse parses src into an expression tree. "and" binds tighter than "or".
func Parse(src string) (Expr, error) {
	toks, err := lexer.Tokenize(src)
	if err != nil {
		return nil, err
	}
	p := &parser{toks: toks}
	e, err := p.or()
	if err != nil {
		return nil, err
	}
	if p.peek().Kind != lexer.EOF {
		return nil, fmt.Errorf("unexpected %q at %d", p.peek().Text, p.peek().Pos)
	}
	return e, nil
}

func (p *parser) peek() lexer.Token { return p.toks[p.pos] }
func (p *parser) next() lexer.Token  { t := p.toks[p.pos]; p.pos++; return t }

func (p *parser) or() (Expr, error) {
	left, err := p.and()
	if err != nil {
		return nil, err
	}
	for p.peek().Kind == lexer.Ident && p.peek().Text == "or" {
		p.next()
		right, err := p.and()
		if err != nil {
			return nil, err
		}
		left = Or{left, right}
	}
	return left, nil
}

func (p *parser) and() (Expr, error) {
	left, err := p.primary()
	if err != nil {
		return nil, err
	}
	for p.peek().Kind == lexer.Ident && p.peek().Text == "and" {
		p.next()
		right, err := p.primary()
		if err != nil {
			return nil, err
		}
		left = And{left, right}
	}
	return left, nil
}

func (p *parser) primary() (Expr, error) {
	if p.peek().Kind == lexer.LParen {
		p.next()
		e, err := p.or()
		if err != nil {
			return nil, err
		}
		if p.next().Kind != lexer.RParen {
			return nil, fmt.Errorf("expected )")
		}
		return e, nil
	}
	field := p.next()
	if field.Kind != lexer.Ident {
		return nil, fmt.Errorf("expected field name at %d", field.Pos)
	}
	op := p.next()
	if op.Kind != lexer.Op {
		return nil, fmt.Errorf("expected operator at %d", op.Pos)
	}
	val := p.next()
	switch val.Kind {
	case lexer.String:
		return Compare{field.Text, op.Text, Value{IsString: true, Str: val.Text}}, nil
	case lexer.Number:
		n, err := strconv.ParseFloat(val.Text, 64)
		if err != nil {
			return nil, err
		}
		return Compare{field.Text, op.Text, Value{Num: n}}, nil
	}
	return nil, fmt.Errorf("expected value at %d", val.Pos)
}
