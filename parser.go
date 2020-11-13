package sqltools

import (
	"fmt"
	"io"
)

type SelectStatement struct {
	Fields    []string
	TableName string
}

type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS

	// Literals
	IDENT // fields, table_name

	// Misc characters
	ASTERISK // *
	COMMA    // ,

	// Keywords
	SELECT
	FROM
	ALERT
	TABLE
	DROP
)

func isWhitespace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n'
}

func isLetter(ch rune) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

func (p *Parser) scan() (Token, string) {
	if p.buf.n == 1 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}
	p.buf.tok, p.buf.lit = p.s.Scan()
	return p.buf.tok, p.buf.lit
}

func (p *Parser) unscan() {
	p.buf.n = 1
}

func (p *Parser) scanWithoutWhiteSpace() (Token, string) {
	token, lit := p.scan()
	if token == WS {
		return p.scan()
	}

	return token, lit
}

func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{}
	if tok, lit := p.scanWithoutWhiteSpace(); tok != SELECT {
		return nil, fmt.Errorf("found %q, expected SELECT", lit)
	}
	for {
		tok, lit := p.scanWithoutWhiteSpace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("found %q, expected field", lit)
		}

		stmt.Fields = append(stmt.Fields, lit)
		if tok, _ := p.scanWithoutWhiteSpace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	tok, lit := p.scanWithoutWhiteSpace()
	if tok != FROM {
		return nil, fmt.Errorf("found %q, expected FROM", lit)
	}
	tok, lit = p.scanWithoutWhiteSpace()

	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit

	return stmt, nil

}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		s: NewScanner(reader),
	}
}
