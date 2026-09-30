// Package lexer tokenizes inventory query expressions such as
//   name = "blue \"widget\"" and qty > 10
package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

type Kind int

const (
	EOF Kind = iota
	Ident
	Number
	String
	Op
	LParen
	RParen
)

func (k Kind) String() string {
	return [...]string{"EOF", "Ident", "Number", "String", "Op", "LParen", "RParen"}[k]
}

type Token struct {
	Kind Kind
	Text string
	Pos  int
}

// Tokenize splits src into tokens. String literals are double-quoted and
// support the escapes \\ \" \n and \t.
func Tokenize(src string) ([]Token, error) {
	var out []Token
	i := 0
	for i < len(src) {
		c := rune(src[i])
		switch {
		case unicode.IsSpace(c):
			i++
		case c == '(':
			out = append(out, Token{LParen, "(", i})
			i++
		case c == ')':
			out = append(out, Token{RParen, ")", i})
			i++
		case c == '"':
			text, n, err := readString(src[i:])
			if err != nil {
				return nil, fmt.Errorf("pos %d: %w", i, err)
			}
			out = append(out, Token{String, text, i})
			i += n
		case strings.ContainsRune("=!<>", c):
			j := i + 1
			if j < len(src) && src[j] == '=' {
				j++
			}
			op := src[i:j]
			if op == "!" {
				return nil, fmt.Errorf("pos %d: unexpected '!'", i)
			}
			out = append(out, Token{Op, op, i})
			i = j
		case unicode.IsDigit(c):
			j := i
			for j < len(src) && (unicode.IsDigit(rune(src[j])) || src[j] == '.') {
				j++
			}
			out = append(out, Token{Number, src[i:j], i})
			i = j
		case unicode.IsLetter(c) || c == '_':
			j := i
			for j < len(src) && (unicode.IsLetter(rune(src[j])) || unicode.IsDigit(rune(src[j])) || src[j] == '_') {
				j++
			}
			out = append(out, Token{Ident, src[i:j], i})
			i = j
		default:
			return nil, fmt.Errorf("pos %d: unexpected %q", i, c)
		}
	}
	out = append(out, Token{EOF, "", len(src)})
	return out, nil
}

// readString reads a quoted literal starting at s[0]=='"' and returns the
// unescaped text and the number of bytes consumed.
func readString(s string) (string, int, error) {
	var b strings.Builder
	for i := 1; i < len(s); i++ {
		switch s[i] {
		case '"':
			return b.String(), i + 1, nil
		case '\\':
			if i+1 >= len(s) {
				return "", 0, fmt.Errorf("unterminated escape")
			}
			i++
			switch s[i] {
			case 'n':
				b.WriteByte('\n')
			case 't':
				b.WriteByte('\t')
			case '"':
				b.WriteByte(s[i])
			case '\\':
				b.WriteByte(s[i])
				i++
			default:
				return "", 0, fmt.Errorf("unknown escape \\%c", s[i])
			}
		default:
			b.WriteByte(s[i])
		}
	}
	return "", 0, fmt.Errorf("unterminated string")
}
