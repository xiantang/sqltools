package sqltools

import (
	"fmt"
	"io"
	"strconv"
)

type Statement interface {
}

type ColumnStatement struct {
	ColumnName string
	DataType   Token
	Length     int
	Nullable   bool
	Comment    string
}

type AlterStatement struct {
	TableName string
	Column    Statement
	Option    Token
}

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
	IDENT  // fields, table_name,column_name
	NUMBER // 1234

	// Misc characters
	ASTERISK         // *
	COMMA            // ,
	LeftParentheses  // (
	RightParentheses // )
	COLON            // ;

	// Keywords
	SELECT
	FROM
	ALTER
	TABLE
	ADD
	DROP
	NULL
	COLUMN
	// DataType
	VARCHAR
	COMMENT
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

func (p *Parser) parseSelectStatement() (Statement, error) {
	stmt := &SelectStatement{}
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

func (p *Parser) parseDropStatement() (Statement, error) {
	tok, lit := p.scanWithoutWhiteSpace()
	if tok != COLUMN {
		return nil, fmt.Errorf("found %q, expected COLUMN", lit)
	}
	tok, lit = p.scanWithoutWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected IDENT", lit)
	}

	return ColumnStatement{
		ColumnName: lit,
	}, nil
}

func (p *Parser) parseAddStatement() (Statement, error) {
	tok, lit := p.scanWithoutWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected column name", lit)
	}
	cstmt := ColumnStatement{}
	cstmt.ColumnName = lit
	tok, lit = p.scanWithoutWhiteSpace()
	switch tok {
	case VARCHAR:
		cstmt.DataType = VARCHAR
		tok, lit = p.scanWithoutWhiteSpace()
		if tok != LeftParentheses {
			return nil, fmt.Errorf("found %q, expected LeftParentheses", lit)
		}
		ntok, nlit := p.scanWithoutWhiteSpace()
		if ntok != NUMBER {
			return nil, fmt.Errorf("found %q, expected NUMBER", nlit)
		}
		tok, lit = p.scanWithoutWhiteSpace()
		if tok != RightParentheses {
			return nil, fmt.Errorf("found %q, expected RightParentheses", lit)
		}
		number, err := strconv.Atoi(nlit)
		if err != nil {
			return nil, err
		}
		cstmt.Length = number
		tok, lit = p.scanWithoutWhiteSpace()
		if tok != NULL {
			return nil, fmt.Errorf("found %q, expected RightParentheses", lit)
		}
		cstmt.Nullable = true
		tok, lit = p.scanWithoutWhiteSpace()
		if tok == COMMENT {
			tok, lit = p.scanWithoutWhiteSpace()
			if tok != IDENT {
				return nil, fmt.Errorf("found %q, expected IDENT", lit)
			}
			cstmt.Comment = lit
		}
	}

	return cstmt, nil

}

func (p *Parser) parseAlterStatement() (Statement, error) {
	stmt := &AlterStatement{}
	tok, lit := p.scanWithoutWhiteSpace()
	if tok != TABLE {
		return nil, fmt.Errorf("found %q, expected TABLE", lit)
	}
	tok, lit = p.scanWithoutWhiteSpace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected table name", lit)
	}
	stmt.TableName = lit
	tok, lit = p.scanWithoutWhiteSpace()
	switch tok {
	case DROP:
		stmt.Option = DROP
		cstmt, err := p.parseDropStatement()
		if err != nil {
			return nil, err
		}
		stmt.Column = cstmt
		return stmt, nil
	case ADD:
		stmt.Option = ADD
		cstmt, err := p.parseAddStatement()
		if err != nil {
			return nil, err
		}
		stmt.Column = cstmt
		return stmt, nil
	default:
		return nil, fmt.Errorf("found %q, expected DROP", lit)
	}

}

func (p *Parser) Parse() (Statement, error) {
	tok, lit := p.scanWithoutWhiteSpace()
	switch tok {
	case SELECT:
		return p.parseSelectStatement()
	case ALTER:
		return p.parseAlterStatement()
	default:
		return nil, fmt.Errorf("found %q, expected SELECT or ALTER", lit)
	}

}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		s: NewScanner(reader),
	}
}
